defmodule Valkyrie.Context do
  @moduledoc """
  Project context aggregation from stories, activity, and chat.
  """

  import Ecto.Query, warn: false

  alias Valkyrie.Accounts.Scope
  alias Valkyrie.Chat.{Message, Thread}
  alias Valkyrie.Repo
  alias Valkyrie.Workflow.{StoryComment, StoryHistoryEvent}
  alias Valkyrie.Workspace.{Project, Story}

  def build_project_context(%Scope{} = scope, project_id) do
    with %Project{} = project <-
           Repo.one(
             from p in Project,
               where:
                 p.id == ^project_id and p.organization_id == ^scope.organization_id and
                   is_nil(p.deleted_at)
           ) do
      stories =
        Repo.all(
          from s in Story,
            where: s.project_id == ^project_id and is_nil(s.deleted_at),
            order_by: [asc: s.inserted_at, asc: s.id]
        )

      story_ids = Enum.map(stories, & &1.id)

      comments =
        Repo.all(
          from c in StoryComment,
            where: c.story_id in ^story_ids,
            order_by: [asc: c.inserted_at, asc: c.id]
        )

      history =
        Repo.all(
          from h in StoryHistoryEvent,
            where: h.story_id in ^story_ids,
            order_by: [asc: h.inserted_at, asc: h.id]
        )

      threads =
        Repo.all(
          from t in Thread,
            where:
              t.organization_id == ^scope.organization_id and t.project_id == ^project_id and
                is_nil(t.deleted_at),
            order_by: [asc: t.inserted_at, asc: t.id]
        )

      thread_ids = Enum.map(threads, & &1.id)

      messages =
        Repo.all(
          from m in Message,
            where: m.thread_id in ^thread_ids,
            order_by: [asc: m.inserted_at, asc: m.id]
        )

      {:ok,
       %{
         project: project,
         stories: stories,
         comments: comments,
         history: history,
         threads: threads,
         messages: messages
       }}
    else
      nil -> {:error, :not_found}
    end
  end
end
