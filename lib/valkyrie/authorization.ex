defmodule Valkyrie.Authorization do
  @moduledoc """
  Role-based permission policy for v1 API endpoints.
  """

  @role_rank %{
    "viewer" => 1,
    "member" => 2,
    "admin" => 3,
    "owner" => 4
  }

  @roles ["owner", "admin", "member", "viewer"]

  def roles, do: @roles

  def can?(role, permission) when is_binary(role) and is_binary(permission) do
    case permission do
      "projects.read" -> at_least?(role, "viewer")
      "stories.read" -> at_least?(role, "viewer")
      "chat.read" -> at_least?(role, "viewer")
      "projects.write" -> at_least?(role, "member")
      "stories.write" -> at_least?(role, "member")
      "chat.write" -> at_least?(role, "member")
      "keys.write" -> at_least?(role, "member")
      "membership.write" -> at_least?(role, "owner")
      _ -> false
    end
  end

  defp at_least?(role, baseline) do
    Map.get(@role_rank, role, 0) >= Map.get(@role_rank, baseline, 0)
  end
end
