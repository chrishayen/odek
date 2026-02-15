defmodule ValkyrieWeb.V1.FeatureController do
  use ValkyrieWeb.V1.BaseController

  alias Valkyrie.Workspace

  def create(conn, params) do
    scope = current_scope!(conn)

    case Workspace.create_feature(scope, params) do
      {:ok, feature} ->
        conn
        |> put_status(:created)
        |> json(%{
          id: feature.id,
          project_id: feature.project_id,
          name: feature.name,
          description: feature.description
        })

      {:error, :not_found} ->
        not_found(conn, "project not found")

      {:error, changeset} ->
        validation_error(conn, "feature_error", changeset)
    end
  end
end
