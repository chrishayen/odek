defmodule Valkyrie.Workflow do
  @moduledoc """
  Story workflow: polling, claiming, state transitions, comments, and history.
  """

  import Ecto.Query, warn: false

  alias Valkyrie.Accounts.Scope
  alias Valkyrie.Repo
  alias Valkyrie.Workflow.{StoryClaim, StoryComment, StoryHistoryEvent}
  alias Valkyrie.Workspace
  alias Valkyrie.Workspace.{Project, Story}

  def poll_stories(%Scope{} = scope, project_id, cursor, limit) do
    Workspace.poll_ready_stories(scope, project_id, cursor, limit)
  end

  def claim_story(%Scope{} = scope, story_id) do
    Repo.transact(fn ->
      story =
        Repo.one(
          from s in Story,
            join: p in Project,
            on: p.id == s.project_id,
            where:
              s.id == ^story_id and p.organization_id == ^scope.organization_id and
                is_nil(s.deleted_at) and
                is_nil(p.deleted_at),
            lock: "FOR UPDATE"
        )

      if is_nil(story) do
        Repo.rollback(:not_found)
      end

      existing_claim =
        Repo.one(
          from c in StoryClaim,
            where: c.story_id == ^story.id and is_nil(c.released_at),
            lock: "FOR UPDATE"
        )

      if existing_claim do
        Repo.rollback(:already_claimed)
      end

      {:ok, claim} =
        %StoryClaim{}
        |> StoryClaim.changeset(%{story_id: story.id, claimed_by_user_id: scope.user.id})
        |> Repo.insert()

      {:ok, _story} = Workspace.update_story_state(scope, story.id, "in_progress")

      _ =
        add_history(scope, story.id, "assignment_change", %{
          "action" => "claim",
          "claim_id" => claim.id
        })

      {:ok, claim}
    end)
    |> normalize_tx_result()
  end

  def update_story_state(%Scope{} = scope, story_id, state) do
    with {:ok, story} <- Workspace.update_story_state(scope, story_id, state),
         {:ok, _event} <- add_history(scope, story.id, "state_change", %{"state" => story.state}) do
      {:ok, story}
    end
  end

  def record_field_update(%Scope{} = _scope, _story_id, changed_fields)
      when map_size(changed_fields) == 0,
      do: :ok

  def record_field_update(%Scope{} = scope, story_id, changed_fields) do
    case add_history(scope, story_id, "field_update", %{"fields" => changed_fields}) do
      {:ok, _event} -> :ok
      _ -> :ok
    end
  end

  def add_comment(%Scope{} = scope, story_id, body) do
    with {:ok, _story} <- fetch_story_for_write(scope, story_id),
         {:ok, comment} <-
           %StoryComment{}
           |> StoryComment.changeset(%{
             story_id: story_id,
             author_user_id: scope.user.id,
             body: body
           })
           |> Repo.insert(),
         {:ok, _event} <-
           add_history(scope, story_id, "comment_added", %{"comment_id" => comment.id}) do
      {:ok, comment}
    end
  end

  def list_comments(%Scope{} = scope, story_id, limit) do
    with {:ok, _story} <- fetch_story_for_history(scope, story_id) do
      limit = normalize_limit(limit, 50)

      Repo.all(
        from c in StoryComment,
          where: c.story_id == ^story_id,
          order_by: [desc: c.inserted_at, desc: c.id],
          limit: ^limit
      )
    else
      _ -> []
    end
  end

  def list_history(%Scope{} = scope, story_id, limit) do
    with {:ok, _story} <- fetch_story_for_history(scope, story_id) do
      limit = normalize_limit(limit, 100)

      Repo.all(
        from h in StoryHistoryEvent,
          where: h.story_id == ^story_id,
          order_by: [desc: h.inserted_at, desc: h.id],
          limit: ^limit
      )
    else
      _ -> []
    end
  end

  defp add_history(%Scope{} = scope, story_id, event_type, payload) do
    %StoryHistoryEvent{}
    |> StoryHistoryEvent.changeset(%{
      story_id: story_id,
      actor_user_id: scope.user.id,
      event_type: event_type,
      payload: payload || %{}
    })
    |> Repo.insert()
  end

  defp fetch_story_for_write(%Scope{} = scope, story_id) do
    case Workspace.get_story(scope, story_id) do
      nil -> {:error, :not_found}
      story -> {:ok, story}
    end
  end

  defp fetch_story_for_history(%Scope{} = scope, story_id) do
    story =
      Repo.one(
        from s in Story,
          join: p in Project,
          on: p.id == s.project_id,
          where:
            s.id == ^story_id and p.organization_id == ^scope.organization_id and
              is_nil(p.deleted_at)
      )

    if story, do: {:ok, story}, else: {:error, :not_found}
  end

  defp normalize_tx_result({:ok, value}), do: {:ok, value}
  defp normalize_tx_result({:error, reason}), do: {:error, reason}

  defp normalize_limit(limit, _fallback) when is_integer(limit) and limit > 0, do: limit
  defp normalize_limit(_limit, fallback), do: fallback
end
