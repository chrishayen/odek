defmodule Valkyrie.Organizations.Membership do
  use Ecto.Schema
  import Ecto.Changeset

  @roles ~w(owner admin member viewer)

  @primary_key {:id, :binary_id, autogenerate: true}
  @foreign_key_type :binary_id

  schema "organization_memberships" do
    field :role, :string

    belongs_to :organization, Valkyrie.Organizations.Organization
    belongs_to :user, Valkyrie.Accounts.User, type: :id

    timestamps(type: :utc_datetime_usec)
  end

  def changeset(membership, attrs) do
    membership
    |> cast(attrs, [:organization_id, :user_id, :role])
    |> validate_required([:organization_id, :user_id, :role])
    |> validate_inclusion(:role, @roles)
    |> unique_constraint([:organization_id, :user_id])
  end
end
