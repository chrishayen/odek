defmodule ValkyrieWeb.V1.APIKeyController do
  use ValkyrieWeb.V1.BaseController

  alias Valkyrie.APIKeys
  alias Valkyrie.Accounts.Scope

  def create(conn, params) do
    with {:ok, scope} <- session_scope(conn) do
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
          render_error(conn, :forbidden, "forbidden", "organization mismatch")

        {:error, changeset} ->
          render_error(conn, :bad_request, "api_key_error", inspect(changeset.errors))
      end
    else
      {:error, error_conn} -> error_conn
    end
  end

  def index(conn, params) do
    with {:ok, scope} <- session_scope(conn) do
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
    else
      {:error, error_conn} -> error_conn
    end
  end

  def revoke(conn, %{"key_id" => key_id}) do
    with {:ok, scope} <- session_scope(conn) do
      case APIKeys.revoke_api_key(scope, key_id) do
        {:ok, key} ->
          json(conn, %{id: key.id, revoked: true})

        {:error, :not_found} ->
          render_error(conn, :not_found, "not_found", "api key not found")

        {:error, :forbidden} ->
          render_error(conn, :forbidden, "forbidden", "cannot revoke another user key")

        {:error, changeset} ->
          render_error(conn, :bad_request, "api_key_error", inspect(changeset.errors))
      end
    else
      {:error, error_conn} -> error_conn
    end
  end

  defp session_scope(conn) do
    case conn.assigns[:current_scope] do
      %Scope{user: %{}} = scope -> {:ok, scope}
      _ -> {:error, render_error(conn, :unauthorized, "unauthorized", "invalid session")}
    end
  end
end
