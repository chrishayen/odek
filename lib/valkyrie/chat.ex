defmodule Valkyrie.Chat do
  @moduledoc """
  Project-scoped chat thread/message persistence and pubsub fanout.
  """

  import Ecto.Query, warn: false

  alias Valkyrie.Accounts.Scope
  alias Valkyrie.Chat.{Message, Thread}
  alias Valkyrie.Repo
  alias Valkyrie.Workspace.{Project, Story}

  def topic(project_id), do: "project_events:" <> project_id

  def subscribe(project_id) do
    Phoenix.PubSub.subscribe(Valkyrie.PubSub, topic(project_id))
  end

  def create_thread(%Scope{} = scope, attrs) do
    params = normalize_attrs(attrs)
    project_id = Map.get(params, "project_id")
    story_id = Map.get(params, "story_id")

    with {:ok, %Project{}} <- fetch_project(scope.organization_id, project_id),
         :ok <- validate_story_link(project_id, story_id) do
      %Thread{}
      |> Thread.changeset(
        Map.merge(params, %{
          "organization_id" => scope.organization_id,
          "created_by_user_id" => scope.user.id,
          "story_id" => normalize_optional(story_id)
        })
      )
      |> Repo.insert()
    end
  end

  def list_threads(%Scope{} = scope, project_id, limit) do
    limit = normalize_limit(limit, 50)

    Repo.all(
      from t in Thread,
        where:
          t.organization_id == ^scope.organization_id and t.project_id == ^project_id and
            is_nil(t.deleted_at),
        order_by: [desc: t.inserted_at, desc: t.id],
        limit: ^limit
    )
  end

  def get_thread(%Scope{} = scope, thread_id) do
    Repo.one(
      from t in Thread,
        where:
          t.id == ^thread_id and t.organization_id == ^scope.organization_id and
            is_nil(t.deleted_at)
    )
  end

  def create_message(%Scope{} = scope, thread_id, attrs) do
    params = normalize_attrs(attrs)

    with %Thread{} = thread <- get_thread(scope, thread_id),
         {:ok, message} <-
           %Message{}
           |> Message.changeset(%{
             thread_id: thread.id,
             sender_kind: Map.get(params, "sender_kind", "runtime"),
             sender_user_id: scope.user.id,
             body: Map.get(params, "body", ""),
             metadata: Map.get(params, "metadata", %{})
           })
           |> Repo.insert() do
      broadcast_chat_message(thread.project_id, thread.id, message.id)
      {:ok, message, thread}
    else
      nil -> {:error, :not_found}
      error -> error
    end
  end

  def list_messages(%Scope{} = scope, thread_id, limit) do
    limit = normalize_limit(limit, 100)

    case get_thread(scope, thread_id) do
      nil ->
        []

      _thread ->
        Repo.all(
          from m in Message,
            where: m.thread_id == ^thread_id,
            order_by: [desc: m.inserted_at, desc: m.id],
            limit: ^limit
        )
    end
  end

  defp broadcast_chat_message(project_id, thread_id, message_id) do
    Phoenix.PubSub.broadcast(
      Valkyrie.PubSub,
      topic(project_id),
      {:chat_message,
       %{
         event: "chat_message",
         project_id: project_id,
         thread_id: thread_id,
         message_id: message_id
       }}
    )
  end

  defp fetch_project(org_id, project_id) do
    case Repo.one(
           from p in Project,
             where: p.id == ^project_id and p.organization_id == ^org_id and is_nil(p.deleted_at)
         ) do
      nil -> {:error, :not_found}
      project -> {:ok, project}
    end
  end

  defp validate_story_link(_project_id, nil), do: :ok
  defp validate_story_link(_project_id, ""), do: :ok

  defp validate_story_link(project_id, story_id) do
    case Repo.one(
           from s in Story,
             where: s.id == ^story_id and s.project_id == ^project_id and is_nil(s.deleted_at)
         ) do
      nil -> {:error, :not_found}
      _story -> :ok
    end
  end

  defp normalize_attrs(attrs), do: Map.new(attrs, fn {k, v} -> {to_string(k), v} end)
  defp normalize_optional(nil), do: nil
  defp normalize_optional(""), do: nil
  defp normalize_optional(v), do: v

  defp normalize_limit(limit, _fallback) when is_integer(limit) and limit > 0, do: limit
  defp normalize_limit(_limit, fallback), do: fallback
end
