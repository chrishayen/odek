defmodule Valkyrie.V1Fixtures do
  @moduledoc false

  alias Valkyrie.Accounts.Scope
  alias Valkyrie.AccountsFixtures
  alias Valkyrie.Organizations
  alias Valkyrie.Workspace

  def organization_fixture(attrs \\ %{}) do
    defaults = %{name: "Org #{System.unique_integer([:positive])}"}
    {:ok, org} = Organizations.create_organization(Map.merge(defaults, attrs))
    org
  end

  def membership_fixture(user, organization, role \\ "member") do
    {:ok, membership} =
      Organizations.create_membership(%{
        organization_id: organization.id,
        user_id: user.id,
        role: role
      })

    membership
  end

  def user_with_membership_fixture(role \\ "member") do
    user = AccountsFixtures.user_fixture() |> AccountsFixtures.set_password()
    organization = organization_fixture()
    membership = membership_fixture(user, organization, role)

    %{user: user, organization: organization, membership: membership}
  end

  def api_scope_fixture(user, organization, role \\ "member") do
    %Scope{user: user, organization_id: organization.id, role: role, auth_mode: :api_key}
  end

  def project_fixture(scope, attrs \\ %{}) do
    defaults = %{name: "Project #{System.unique_integer([:positive])}"}
    {:ok, project} = Workspace.create_project(scope, Map.merge(defaults, attrs))
    project
  end
end
