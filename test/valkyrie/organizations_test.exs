defmodule Valkyrie.OrganizationsTest do
  use Valkyrie.DataCase, async: true

  import Valkyrie.AccountsFixtures
  import Valkyrie.V1Fixtures

  alias Valkyrie.Organizations

  describe "ensure_user_membership/1" do
    test "creates a personal owner membership when the user has none" do
      user = user_fixture()
      assert [] == Organizations.list_user_memberships(user.id)

      assert {:ok, membership} = Organizations.ensure_user_membership(user.id)
      assert membership.user_id == user.id
      assert membership.role == "owner"
      assert is_binary(membership.organization_id)
      org_id = membership.organization_id
      assert %{id: ^org_id} = Organizations.get_organization(org_id)
      assert [_only_membership] = Organizations.list_user_memberships(user.id)
    end

    test "returns an existing membership without creating another organization" do
      user = user_fixture()
      organization = organization_fixture()
      existing_membership = membership_fixture(user, organization, "member")

      assert {:ok, membership} = Organizations.ensure_user_membership(user.id)
      assert membership.id == existing_membership.id
      assert membership.organization_id == organization.id
      assert membership.role == "member"
      assert [_only_membership] = Organizations.list_user_memberships(user.id)
    end
  end
end
