defmodule Valkyrie.APIKeys.APIKey do
  use Ecto.Schema
  import Ecto.Changeset

  @primary_key {:id, :binary_id, autogenerate: true}
  @foreign_key_type :binary_id

  schema "api_keys" do
    field :key_prefix, :string
    field :key_hash, :binary
    field :revoked_at, :utc_datetime_usec

    belongs_to :organization, Valkyrie.Organizations.Organization
    belongs_to :user, Valkyrie.Accounts.User, type: :id

    timestamps(type: :utc_datetime_usec)
  end

  def changeset(api_key, attrs) do
    api_key
    |> cast(attrs, [:organization_id, :user_id, :key_prefix, :key_hash, :revoked_at])
    |> validate_required([:organization_id, :user_id, :key_prefix, :key_hash])
    |> unique_constraint(:key_hash)
  end
end
