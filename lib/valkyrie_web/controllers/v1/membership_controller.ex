defmodule ValkyrieWeb.V1.MembershipController do
  use ValkyrieWeb.V1.BaseController

  alias Valkyrie.Organizations

  def update_role(conn, params) do
    scope = current_scope!(conn)
    organization_id = Map.get(params, "organization_id")
    user_id = parse_int(Map.get(params, "user_id"), -1)
    role = Map.get(params, "role")

    cond do
      not is_binary(organization_id) or user_id <= 0 or not is_binary(role) ->
        bad_request(conn, "invalid_request", "organization_id, user_id and role are required")

      organization_id != scope.organization_id ->
        forbidden(conn, "cross-organization role updates are forbidden")

      true ->
        case Organizations.update_membership_role(organization_id, user_id, role) do
          {:ok, membership} ->
            json(conn, %{
              organization_id: membership.organization_id,
              user_id: membership.user_id,
              role: membership.role
            })

          {:error, :not_found} ->
            not_found(conn, "membership not found")

          {:error, changeset} ->
            validation_error(conn, "invalid_role", changeset)
        end
    end
  end
end
