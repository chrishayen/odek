defmodule Valkyrie.Workspace.Story do
  use Ecto.Schema
  import Ecto.Changeset

  @states ~w(backlog ready in_progress review done)

  @primary_key {:id, :binary_id, autogenerate: true}
  @foreign_key_type :binary_id

  schema "stories" do
    field :name, :string
    field :description, :string, default: ""
    field :state, :string, default: "backlog"
    field :deleted_at, :utc_datetime_usec

    belongs_to :project, Valkyrie.Workspace.Project
    belongs_to :feature, Valkyrie.Workspace.Feature

    timestamps(type: :utc_datetime_usec)
  end

  def states, do: @states

  def valid_state?(state), do: state in @states

  def changeset(story, attrs) do
    story
    |> cast(attrs, [:project_id, :feature_id, :name, :description, :state, :deleted_at])
    |> validate_required([:project_id, :name])
    |> validate_length(:name, max: 160)
    |> validate_inclusion(:state, @states)
  end

  def update_details_changeset(story, attrs) do
    story
    |> cast(attrs, [:feature_id, :name, :description])
    |> validate_length(:name, max: 160)
  end

  def state_changeset(story, attrs) do
    story
    |> cast(attrs, [:state])
    |> validate_required([:state])
    |> validate_inclusion(:state, @states)
  end
end
