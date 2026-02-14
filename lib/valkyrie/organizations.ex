defmodule Valkyrie.Organizations do
  @moduledoc """
  Organization tenancy and membership management.
  """

  import Ecto.Query, warn: false

  alias Valkyrie.Authorization
  alias Valkyrie.Organizations.{Membership, Organization}
  alias Valkyrie.Repo

  def list_roles, do: Authorization.roles()

  def create_organization(attrs) do
    %Organization{}
    |> Organization.changeset(attrs)
    |> Repo.insert()
  end

  def get_organization(id), do: Repo.get(Organization, id)

  def create_membership(attrs) do
    %Membership{}
    |> Membership.changeset(attrs)
    |> Repo.insert()
  end

  def get_membership(organization_id, user_id) do
    Repo.one(
      from m in Membership,
        where: m.organization_id == ^organization_id and m.user_id == ^user_id,
        preload: [:organization, :user]
    )
  end

  def list_user_memberships(user_id) do
    Repo.all(
      from m in Membership,
        where: m.user_id == ^user_id,
        preload: [:organization]
    )
  end

  def update_membership_role(organization_id, user_id, role) do
    case get_membership(organization_id, user_id) do
      nil ->
        {:error, :not_found}

      membership ->
        membership
        |> Membership.changeset(%{role: role})
        |> Repo.update()
    end
  end

  def member?(organization_id, user_id) do
    not is_nil(get_membership(organization_id, user_id))
  end
end
