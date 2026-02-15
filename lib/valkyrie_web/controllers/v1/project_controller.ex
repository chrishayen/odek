defmodule ValkyrieWeb.V1.ProjectController do
  use ValkyrieWeb.V1.BaseController

  alias Valkyrie.Workspace

  def create(conn, params) do
    scope = current_scope!(conn)

    case Workspace.create_project(scope, params) do
      {:ok, project} ->
        conn
        |> put_status(:created)
        |> json(%{
          id: project.id,
          organization_id: project.organization_id,
          name: project.name,
          definition_of_done: project.definition_of_done
        })

      {:error, :invalid_request} ->
        bad_request(conn, "invalid_request", "organization_id is required")

      {:error, :forbidden} ->
        forbidden(conn, "organization mismatch")

      {:error, changeset} ->
        validation_error(conn, "project_error", changeset)
    end
  end

  def show(conn, %{"project_id" => project_id}) do
    scope = current_scope!(conn)

    case Workspace.get_project(scope, project_id) do
      nil ->
        not_found(conn, "project not found")

      project ->
        json(conn, %{
          id: project.id,
          organization_id: project.organization_id,
          name: project.name,
          definition_of_done: project.definition_of_done,
          created_at: project.inserted_at,
          updated_at: project.updated_at
        })
    end
  end

  def update(conn, %{"project_id" => project_id} = params) do
    scope = current_scope!(conn)

    case Workspace.update_project(scope, project_id, params) do
      {:ok, project} ->
        json(conn, %{
          id: project.id,
          organization_id: project.organization_id,
          name: project.name,
          definition_of_done: project.definition_of_done,
          updated: true
        })

      {:error, :not_found} ->
        not_found(conn, "project not found")

      {:error, changeset} ->
        validation_error(conn, "project_error", changeset)
    end
  end
end
