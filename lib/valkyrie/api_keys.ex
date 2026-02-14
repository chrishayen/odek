defmodule Valkyrie.APIKeys do
  @moduledoc """
  API key lifecycle and bearer authentication.
  """

  import Ecto.Query, warn: false

  alias Valkyrie.APIKeys.APIKey
  alias Valkyrie.Accounts
  alias Valkyrie.Accounts.Scope
  alias Valkyrie.Organizations
  alias Valkyrie.Repo

  def create_api_key(%Scope{user: %{id: user_id}} = _scope, organization_id) do
    if Organizations.member?(organization_id, user_id) do
      raw_key = "sk_" <> Base.url_encode64(:crypto.strong_rand_bytes(24), padding: false)
      prefix = String.slice(raw_key, 0, 12)

      attrs = %{
        organization_id: organization_id,
        user_id: user_id,
        key_prefix: prefix,
        key_hash: hash_key(raw_key)
      }

      case %APIKey{} |> APIKey.changeset(attrs) |> Repo.insert() do
        {:ok, api_key} -> {:ok, api_key, raw_key}
        error -> error
      end
    else
      {:error, :forbidden}
    end
  end

  def list_api_keys(%Scope{user: %{id: user_id}}, organization_id \\ nil) do
    base_query =
      from k in APIKey,
        where: k.user_id == ^user_id,
        where: is_nil(k.revoked_at),
        order_by: [desc: k.inserted_at]

    query =
      if is_binary(organization_id) and organization_id != "" do
        from k in base_query, where: k.organization_id == ^organization_id
      else
        base_query
      end

    Repo.all(query)
  end

  def revoke_api_key(%Scope{user: %{id: user_id}}, key_id) do
    case Repo.get(APIKey, key_id) do
      %APIKey{user_id: ^user_id} = key ->
        key
        |> APIKey.changeset(%{revoked_at: DateTime.utc_now()})
        |> Repo.update()

      nil ->
        {:error, :not_found}

      _ ->
        {:error, :forbidden}
    end
  end

  def authenticate_api_key(raw_key) when is_binary(raw_key) do
    key_hash = hash_key(raw_key)

    api_key =
      Repo.one(
        from k in APIKey,
          where: k.key_hash == ^key_hash and is_nil(k.revoked_at),
          preload: [:user]
      )

    with %APIKey{} = key <- api_key,
         membership when not is_nil(membership) <-
           Organizations.get_membership(key.organization_id, key.user_id),
         %Accounts.User{} = user <- Accounts.get_user(key.user_id) do
      {:ok, Scope.for_api_key(user, key.organization_id, membership.role, key.id)}
    else
      _ -> {:error, :unauthorized}
    end
  end

  def authenticate_api_key(_), do: {:error, :unauthorized}

  def key_metadata(api_key) do
    %{
      id: api_key.id,
      organization_id: api_key.organization_id,
      user_id: api_key.user_id,
      key_prefix: api_key.key_prefix,
      created_at: api_key.inserted_at,
      revoked_at: api_key.revoked_at
    }
  end

  defp hash_key(raw_key) do
    :crypto.hash(:sha256, raw_key)
  end
end
