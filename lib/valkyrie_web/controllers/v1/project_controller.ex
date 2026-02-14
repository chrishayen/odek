defmodule ValkyrieWeb.V1.ProjectController do
  use ValkyrieWeb.V1.BaseController

  alias Valkyrie.Workspace

  def create(conn, params) do
    scope = current_scope!(conn)

    with :ok <- authorize(conn, "projects.write"),
         {:ok, project} <- Workspace.create_project(scope, params) do
      conn
      |> put_status(:created)
      |> json(%{
        id: project.id,
        organization_id: project.organization_id,
        name: project.name,
        definition_of_done: project.definition_of_done
      })
    else
      {:error, error_conn} when is_struct(error_conn, Plug.Conn) ->
        error_conn

      {:error, :invalid_request} ->
        render_error(conn, :bad_request, "invalid_request", "organization_id is required")

      {:error, :forbidden} ->
        render_error(conn, :forbidden, "forbidden", "organization mismatch")

      {:error, changeset} ->
        render_error(conn, :bad_request, "project_error", inspect(changeset.errors))
    end
  end

  def show(conn, %{"project_id" => project_id}) do
    scope = current_scope!(conn)

    with :ok <- authorize(conn, "projects.read"),
         project when not is_nil(project) <- Workspace.get_project(scope, project_id) do
      json(conn, %{
        id: project.id,
        organization_id: project.organization_id,
        name: project.name,
        definition_of_done: project.definition_of_done,
        created_at: project.inserted_at,
        updated_at: project.updated_at
      })
    else
      {:error, error_conn} when is_struct(error_conn, Plug.Conn) -> error_conn
      _ -> render_error(conn, :not_found, "not_found", "project not found")
    end
  end

  def update(conn, %{"project_id" => project_id} = params) do
    scope = current_scope!(conn)

    with :ok <- authorize(conn, "projects.write"),
         {:ok, project} <- Workspace.update_project(scope, project_id, params) do
      json(conn, %{
        id: project.id,
        organization_id: project.organization_id,
        name: project.name,
        definition_of_done: project.definition_of_done,
        updated: true
      })
    else
      {:error, error_conn} when is_struct(error_conn, Plug.Conn) ->
        error_conn

      {:error, :not_found} ->
        render_error(conn, :not_found, "not_found", "project not found")

      {:error, changeset} ->
        render_error(conn, :bad_request, "project_error", inspect(changeset.errors))
    end
  end
end
