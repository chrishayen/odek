defmodule Valkyrie.Workspace.Feature do
  use Ecto.Schema
  import Ecto.Changeset

  @primary_key {:id, :binary_id, autogenerate: true}
  @foreign_key_type :binary_id

  schema "features" do
    field :name, :string
    field :description, :string, default: ""
    field :deleted_at, :utc_datetime_usec

    belongs_to :project, Valkyrie.Workspace.Project

    timestamps(type: :utc_datetime_usec)
  end

  def changeset(feature, attrs) do
    feature
    |> cast(attrs, [:project_id, :name, :description, :deleted_at])
    |> validate_required([:project_id, :name])
    |> validate_length(:name, max: 160)
  end
end
