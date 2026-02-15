defmodule ValkyrieWeb.WorkspaceScope do
  @moduledoc false

  alias Valkyrie.Accounts.Scope
  alias Valkyrie.Organizations

  def resolve(nil), do: nil

  def resolve(%Scope{user: nil}), do: nil

  def resolve(%Scope{} = scope) do
    membership =
      scope.user.id
      |> Organizations.list_user_memberships()
      |> pick_membership(scope.organization_id)
      |> case do
        nil ->
          case Organizations.ensure_user_membership(scope.user.id) do
            {:ok, ensured_membership} -> ensured_membership
            {:error, _reason} -> nil
          end

        selected_membership ->
          selected_membership
      end

    case membership do
      nil ->
        nil

      membership ->
        %Scope{scope | organization_id: membership.organization_id, role: membership.role}
    end
  end

  defp pick_membership([], _organization_id), do: nil

  defp pick_membership(memberships, organization_id) when is_binary(organization_id) do
    Enum.find(memberships, &(&1.organization_id == organization_id)) ||
      sort_memberships(memberships)
  end

  defp pick_membership(memberships, _organization_id), do: sort_memberships(memberships)

  defp sort_memberships(memberships) do
    memberships
    |> Enum.sort_by(fn membership ->
      organization_name =
        case membership.organization do
          %{name: name} when is_binary(name) -> String.downcase(name)
          _ -> ""
        end

      {organization_name, membership.organization_id}
    end)
    |> List.first()
  end
end
