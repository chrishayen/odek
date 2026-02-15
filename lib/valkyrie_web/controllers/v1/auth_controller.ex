defmodule ValkyrieWeb.V1.AuthController do
  use ValkyrieWeb.V1.BaseController

  alias Valkyrie.Accounts
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
    user = current_scope!(conn).user

    with current_password when is_binary(current_password) <- Map.get(params, "current_password"),
         new_password when is_binary(new_password) <- Map.get(params, "new_password"),
         %Accounts.User{} <- Accounts.get_user_by_email_and_password(user.email, current_password),
         {:ok, {_updated, _tokens}} <-
           Accounts.update_user_password(user, %{password: new_password}) do
      json(conn, %{ok: true})
    else
      {:error, changeset} ->
        validation_error(conn, "invalid_password", changeset)

      _ ->
        unauthorized(conn, "invalid credentials")
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
    user = current_scope!(conn).user
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
        forbidden(conn, "session mismatch")

      true ->
        bad_request(conn, "invalid_request", "organization_id is required")

      _ ->
        forbidden(conn, "user is not a member of target organization")
    end
  end

  def frontend_projects(conn, params) do
    scope = current_scope!(conn)
    user = scope.user
    requested_org = Map.get(params, "organization_id", scope.organization_id)

    cond do
      user.must_change_password ->
        render_error(conn, :forbidden, "password_change_required", "password change required")

      requested_org != scope.organization_id ->
        forbidden(conn, "session organization mismatch")

      true ->
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
    end
  end
end
