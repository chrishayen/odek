defmodule Valkyrie.Prompts.Skill do
  use Ecto.Schema
  import Ecto.Changeset

  @primary_key {:id, :binary_id, autogenerate: true}
  @foreign_key_type :binary_id

  schema "skills" do
    field :name, :string
    field :body, :string, default: ""
    field :deleted_at, :utc_datetime_usec

    belongs_to :organization, Valkyrie.Organizations.Organization
    belongs_to :owner_user, Valkyrie.Accounts.User, type: :id

    timestamps(type: :utc_datetime_usec)
  end

  def changeset(skill, attrs) do
    skill
    |> cast(attrs, [:organization_id, :owner_user_id, :name, :body, :deleted_at])
    |> validate_required([:organization_id, :owner_user_id, :name])
    |> validate_length(:name, max: 160)
  end
end
