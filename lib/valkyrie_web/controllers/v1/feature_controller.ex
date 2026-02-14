defmodule ValkyrieWeb.V1.FeatureController do
  use ValkyrieWeb.V1.BaseController

  alias Valkyrie.Workspace

  def create(conn, params) do
    scope = current_scope!(conn)

    with :ok <- authorize(conn, "projects.write"),
         {:ok, feature} <- Workspace.create_feature(scope, params) do
      conn
      |> put_status(:created)
      |> json(%{
        id: feature.id,
        project_id: feature.project_id,
        name: feature.name,
        description: feature.description
      })
    else
      {:error, error_conn} when is_struct(error_conn, Plug.Conn) ->
        error_conn

      {:error, :not_found} ->
        render_error(conn, :not_found, "not_found", "project not found")

      {:error, changeset} ->
        render_error(conn, :bad_request, "feature_error", inspect(changeset.errors))
    end
  end
end
