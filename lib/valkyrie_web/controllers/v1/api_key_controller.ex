defmodule ValkyrieWeb.V1.APIKeyController do
  use ValkyrieWeb.V1.BaseController

  alias Valkyrie.APIKeys

  def create(conn, params) do
    scope = current_scope!(conn)
    organization_id = Map.get(params, "organization_id", scope.organization_id)

    case APIKeys.create_api_key(scope, organization_id) do
      {:ok, key, raw_key} ->
        conn
        |> put_status(:created)
        |> json(%{
          id: key.id,
          organization_id: key.organization_id,
          user_id: key.user_id,
          key_prefix: key.key_prefix,
          api_key: raw_key
        })

      {:error, :forbidden} ->
        forbidden(conn, "organization mismatch")

      {:error, changeset} ->
        validation_error(conn, "api_key_error", changeset)
    end
  end

  def index(conn, params) do
    scope = current_scope!(conn)
    organization_id = Map.get(params, "organization_id")

    keys =
      APIKeys.list_api_keys(scope, organization_id)
      |> Enum.map(fn key ->
        meta = APIKeys.key_metadata(key)

        %{
          id: meta.id,
          organization_id: meta.organization_id,
          user_id: meta.user_id,
          key_prefix: meta.key_prefix,
          created_at: meta.created_at,
          revoked_at: meta.revoked_at
        }
      end)

    json(conn, %{api_keys: keys})
  end

  def revoke(conn, %{"key_id" => key_id}) do
    scope = current_scope!(conn)

    case APIKeys.revoke_api_key(scope, key_id) do
      {:ok, key} ->
        json(conn, %{id: key.id, revoked: true})

      {:error, :not_found} ->
        not_found(conn, "api key not found")

      {:error, :forbidden} ->
        forbidden(conn, "cannot revoke another user key")

      {:error, changeset} ->
        validation_error(conn, "api_key_error", changeset)
    end
  end
end
