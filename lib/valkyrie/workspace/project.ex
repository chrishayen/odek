defmodule Valkyrie.Workspace.Project do
  use Ecto.Schema
  import Ecto.Changeset

  @primary_key {:id, :binary_id, autogenerate: true}
  @foreign_key_type :binary_id

  schema "projects" do
    field :name, :string
    field :definition_of_done, :string, default: ""
    field :deleted_at, :utc_datetime_usec

    belongs_to :organization, Valkyrie.Organizations.Organization

    timestamps(type: :utc_datetime_usec)
  end

  def changeset(project, attrs) do
    project
    |> cast(attrs, [:organization_id, :name, :definition_of_done, :deleted_at])
    |> update_change(:name, &String.trim/1)
    |> validate_required([:organization_id, :name])
    |> validate_length(:name, max: 160)
    |> unique_constraint(:name, name: :projects_org_lower_name_uniq)
  end

  def update_changeset(project, attrs) do
    project
    |> cast(attrs, [:name, :definition_of_done])
    |> update_change(:name, &String.trim/1)
    |> validate_length(:name, max: 160)
    |> unique_constraint(:name, name: :projects_org_lower_name_uniq)
  end
end
