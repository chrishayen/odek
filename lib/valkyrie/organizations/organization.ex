defmodule Valkyrie.Organizations.Organization do
  use Ecto.Schema
  import Ecto.Changeset

  @primary_key {:id, :binary_id, autogenerate: true}
  @foreign_key_type :binary_id

  schema "organizations" do
    field :name, :string
    field :deleted_at, :utc_datetime_usec

    has_many :memberships, Valkyrie.Organizations.Membership

    timestamps(type: :utc_datetime_usec)
  end

  def changeset(org, attrs) do
    org
    |> cast(attrs, [:name, :deleted_at])
    |> validate_required([:name])
    |> validate_length(:name, max: 160)
    |> unique_constraint(:name)
  end
end
