defmodule Valkyrie.Workflow.StoryHistoryEvent do
  use Ecto.Schema
  import Ecto.Changeset

  @primary_key {:id, :binary_id, autogenerate: true}
  @foreign_key_type :binary_id

  schema "story_history_events" do
    field :event_type, :string
    field :payload, :map, default: %{}

    belongs_to :story, Valkyrie.Workspace.Story
    belongs_to :actor_user, Valkyrie.Accounts.User, type: :id

    timestamps(type: :utc_datetime_usec, updated_at: false)
  end

  def changeset(event, attrs) do
    event
    |> cast(attrs, [:story_id, :actor_user_id, :event_type, :payload])
    |> validate_required([:story_id, :actor_user_id, :event_type])
  end
end
