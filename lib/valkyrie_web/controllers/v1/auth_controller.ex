defmodule ValkyrieWeb.V1.AuthController do
  use ValkyrieWeb.V1.BaseController

  alias Valkyrie.Accounts
  alias Valkyrie.Accounts.Scope
  alias Valkyrie.Organizations
  alias Valkyrie.Workspace

  @session_id_key :frontend_session_id
  @active_org_key :active_organization_id

  def login(conn, params) do
    email = Map.get(params, "email", "")
    password = Map.get(params, "password", "")
    organization_id = Map.get(params, "organization_id")

    with false <- blank?(email),
         false <- blank?(password),
         false <- blank?(organization_id),
         %Accounts.User{} = user <- Accounts.get_user_by_email_and_password(email, password),
         membership when not is_nil(membership) <-
           Organizations.get_membership(organization_id, user.id) do
      session_id = Ecto.UUID.generate()

      conn =
        conn
        |> put_session(:user_token, Accounts.generate_user_session_token(user))
        |> put_session(@session_id_key, session_id)
        |> put_session(@active_org_key, organization_id)

      conn
      |> put_status(:created)
      |> json(%{
        session_id: session_id,
        organization_id: organization_id,
        must_change_password: user.must_change_password,
        role: membership.role
      })
    else
      true ->
        render_error(
          conn,
          :bad_request,
          "invalid_request",
          "email, password and organization_id are required"
        )

      nil ->
        render_error(conn, :unauthorized, "unauthorized", "invalid credentials")

      _ ->
        render_error(conn, :unauthorized, "unauthorized", "invalid credentials")
    end
  end

  def change_password(conn, params) do
    user = conn.assigns.current_scope && conn.assigns.current_scope.user

    with %Accounts.User{} = user <- user,
         current_password when is_binary(current_password) <- Map.get(params, "current_password"),
         new_password when is_binary(new_password) <- Map.get(params, "new_password"),
         %Accounts.User{} <- Accounts.get_user_by_email_and_password(user.email, current_password),
         {:ok, {_updated, _tokens}} <-
           Accounts.update_user_password(user, %{password: new_password}) do
      json(conn, %{ok: true})
    else
      nil ->
        render_error(conn, :unauthorized, "unauthorized", "invalid session")

      {:error, changeset} ->
        render_error(conn, :bad_request, "invalid_password", inspect(changeset.errors))

      _ ->
        render_error(conn, :unauthorized, "unauthorized", "invalid credentials")
    end
  end

  def create_session(conn, _params) do
    render_error(conn, :not_found, "not_found", "endpoint disabled; use /v1/auth/login")
  end

  def revoke_session(conn, %{"session_id" => session_id}) do
    if get_session(conn, @session_id_key) == session_id do
      conn
      |> configure_session(drop: true)
      |> json(%{session_id: session_id, revoked: true})
    else
      render_error(conn, :not_found, "not_found", "session not found")
    end
  end

  def switch_org(conn, %{"session_id" => session_id} = params) do
    active_session_id = get_session(conn, @session_id_key)
    user = conn.assigns.current_scope && conn.assigns.current_scope.user
    organization_id = Map.get(params, "organization_id")

    with true <- active_session_id == session_id,
         %Accounts.User{} = user <- user,
         false <- blank?(organization_id),
         membership when not is_nil(membership) <-
           Organizations.get_membership(organization_id, user.id) do
      conn
      |> put_session(@active_org_key, organization_id)
      |> json(%{session_id: session_id, organization_id: organization_id, role: membership.role})
    else
      false ->
        render_error(conn, :forbidden, "forbidden", "session mismatch")

      nil ->
        render_error(conn, :unauthorized, "unauthorized", "invalid session")

      true ->
        render_error(conn, :bad_request, "invalid_request", "organization_id is required")

      _ ->
        render_error(conn, :forbidden, "forbidden", "user is not a member of target organization")
    end
  end

  def frontend_projects(conn, params) do
    user = conn.assigns.current_scope && conn.assigns.current_scope.user

    with %Accounts.User{} = user <- user,
         false <- user.must_change_password,
         active_org when is_binary(active_org) <- get_session(conn, @active_org_key),
         requested_org <- Map.get(params, "organization_id", active_org),
         true <- requested_org == active_org,
         membership when not is_nil(membership) <-
           Organizations.get_membership(active_org, user.id) do
      scope = %Scope{
        user: user,
        organization_id: active_org,
        role: membership.role,
        auth_mode: :session
      }

      projects =
        Workspace.list_projects(scope)
        |> Enum.map(fn project ->
          %{
            id: project.id,
            organization_id: project.organization_id,
            name: project.name,
            definition_of_done: project.definition_of_done
          }
        end)

      json(conn, %{projects: projects})
    else
      nil ->
        render_error(conn, :unauthorized, "unauthorized", "invalid session")

      true ->
        render_error(conn, :forbidden, "password_change_required", "password change required")

      false ->
        render_error(conn, :forbidden, "forbidden", "session organization mismatch")

      _ ->
        render_error(conn, :forbidden, "forbidden", "session organization mismatch")
    end
  end
end
