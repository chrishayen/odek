defmodule ValkyrieWeb.Plugs.LoadAPIPrincipal do
  @moduledoc """
  Resolves API request principal from bearer API key or Phoenix session.
  """

  import Plug.Conn

  alias Valkyrie.APIKeys
  alias Valkyrie.Accounts.Scope
  alias Valkyrie.Organizations

  @active_org_key :active_organization_id

  def init(opts), do: opts

  def call(conn, opts) do
    allow_api_key = Keyword.get(opts, :allow_api_key, true)
    allow_session = Keyword.get(opts, :allow_session, true)
    require_org_context = Keyword.get(opts, :require_org_context, false)

    with {:ok, scope} <- resolve_scope(conn, allow_api_key, allow_session, require_org_context) do
      assign(conn, :current_scope, scope)
    else
      _ -> unauthorized(conn)
    end
  end

  defp resolve_scope(conn, allow_api_key, allow_session, require_org_context) do
    case get_req_header(conn, "authorization") do
      [header] when allow_api_key ->
        resolve_api_key_scope(header)

      [_header] when allow_session ->
        resolve_session_scope(conn, true, require_org_context)

      [_header] ->
        {:error, :unauthorized}

      _ when allow_session ->
        resolve_session_scope(conn, true, require_org_context)

      _ ->
        {:error, :unauthorized}
    end
  end

  defp resolve_api_key_scope(header) do
    with {:ok, token} <- bearer_token(header),
         {:ok, scope} <- APIKeys.authenticate_api_key(token) do
      {:ok, scope}
    end
  end

  defp resolve_session_scope(conn, true, false) do
    case conn.assigns[:current_scope] do
      %Scope{user: %{}} = scope -> {:ok, %Scope{scope | auth_mode: :session}}
      _ -> {:error, :invalid_session}
    end
  end

  defp resolve_session_scope(conn, true, true) do
    with %Scope{user: %{id: user_id}} = scope <- conn.assigns[:current_scope],
         organization_id when is_binary(organization_id) <- get_session(conn, @active_org_key),
         membership when not is_nil(membership) <-
           Organizations.get_membership(organization_id, user_id) do
      {:ok,
       %Scope{
         scope
         | organization_id: organization_id,
           role: membership.role,
           auth_mode: :session
       }}
    else
      _ -> {:error, :invalid_session_scope}
    end
  end

  defp bearer_token(header) do
    case String.split(header, " ", parts: 2) do
      [scheme, token] when scheme in ["Bearer", "bearer"] and byte_size(token) > 0 ->
        {:ok, token}

      _ ->
        {:error, :invalid_authorization}
    end
  end

  defp unauthorized(conn) do
    body = Jason.encode!(%{error: %{code: "unauthorized", message: "authentication required"}})

    conn
    |> put_resp_content_type("application/json")
    |> send_resp(:unauthorized, body)
    |> halt()
  end
end
