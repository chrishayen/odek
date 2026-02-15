defmodule Valkyrie.Organizations do
  @moduledoc """
  Organization tenancy and membership management.
  """

  import Ecto.Query, warn: false

  alias Valkyrie.Authorization
  alias Valkyrie.Accounts.User
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

  def ensure_user_membership(user_id, opts \\ []) when is_integer(user_id) do
    role = Keyword.get(opts, :role, "owner")
    organization_name = Keyword.get(opts, :organization_name, unique_organization_name(user_id))

    Repo.transact(fn ->
      user_exists? =
        case Repo.one(
               from u in User,
                 where: u.id == ^user_id,
                 lock: "FOR UPDATE",
                 select: u.id
             ) do
          nil -> false
          _id -> true
        end

      if user_exists? do
        case Repo.one(
               from m in Membership,
                 where: m.user_id == ^user_id,
                 order_by: [asc: m.inserted_at, asc: m.id],
                 limit: 1,
                 preload: [:organization]
             ) do
          nil ->
            with {:ok, organization} <- create_organization(%{name: organization_name}),
                 {:ok, membership} <-
                   create_membership(%{
                     organization_id: organization.id,
                     user_id: user_id,
                     role: role
                   }) do
              {:ok, membership}
            else
              {:error, changeset} -> {:error, changeset}
            end

          membership ->
            {:ok, membership}
        end
      else
        {:error, :not_found}
      end
    end)
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

  defp unique_organization_name(user_id) do
    suffix =
      Ecto.UUID.generate()
      |> String.split("-", parts: 2)
      |> List.first()

    "org-#{user_id}-#{suffix}"
  end
end
