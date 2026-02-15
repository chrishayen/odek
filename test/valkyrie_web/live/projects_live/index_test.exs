defmodule ValkyrieWeb.ProjectsLive.IndexTest do
  use ValkyrieWeb.ConnCase, async: true

  import Phoenix.LiveViewTest
  import Valkyrie.V1Fixtures

  describe "authentication" do
    test "redirects anonymous users", %{conn: conn} do
      assert {:error, {:redirect, %{to: "/users/log-in"}}} = live(conn, ~p"/projects")
    end
  end

  describe "project listing and search" do
    setup %{conn: conn} do
      %{user: user, organization: organization} = user_with_membership_fixture("member")
      scope = api_scope_fixture(user, organization, "member")

      _ = project_fixture(scope, %{name: "Zephyr"})
      _ = project_fixture(scope, %{name: "alpha"})
      _ = project_fixture(scope, %{name: "Beta"})

      %{conn: log_in_user(conn, user), scope: scope}
    end

    test "renders both / and /projects and sorts by alpha", %{conn: conn} do
      {:ok, projects_view, _html} = live(conn, ~p"/projects")
      {:ok, root_view, _html} = live(conn, ~p"/")

      assert has_element?(projects_view, "#projects-page")
      assert has_element?(root_view, "#projects-page")

      assert has_element?(
               projects_view,
               "#projects-grid article:nth-of-type(1) .project-card-title",
               "alpha"
             )

      assert has_element?(
               projects_view,
               "#projects-grid article:nth-of-type(2) .project-card-title",
               "Beta"
             )

      assert has_element?(
               projects_view,
               "#projects-grid article:nth-of-type(3) .project-card-title",
               "Zephyr"
             )
    end

    test "filters projects instantly by name", %{conn: conn} do
      {:ok, view, _html} = live(conn, ~p"/projects")

      _ =
        view
        |> form("#projects-search-form", search: %{q: "be"})
        |> render_change()

      assert has_element?(view, "#projects-grid article .project-card-title", "Beta")
      refute has_element?(view, "#projects-grid article .project-card-title", "alpha")
      refute has_element?(view, "#projects-grid article .project-card-title", "Zephyr")
    end
  end

  describe "project creation" do
    setup %{conn: conn} do
      %{user: user, organization: organization} = user_with_membership_fixture("member")
      scope = api_scope_fixture(user, organization, "member")

      %{conn: log_in_user(conn, user), scope: scope}
    end

    test "shows empty state and can create a project", %{conn: conn} do
      {:ok, view, _html} = live(conn, ~p"/projects")

      assert has_element?(view, "#projects-empty")
      assert has_element?(view, "#empty-state-new-project")

      _ = view |> element("#new-project-button") |> render_click()
      assert has_element?(view, "#project-create-modal")
      assert has_element?(view, "#project-create-form")

      submit_result =
        view
        |> form("#project-create-form",
          project: %{name: "Alpha Project", description: "first pass"}
        )
        |> render_submit()

      assert {:error, {:live_redirect, %{to: to}}} = submit_result
      assert String.match?(to, ~r|^/projects/[^/]+$|)

      {:ok, show_view, _html} = live(conn, to)
      assert has_element?(show_view, "#project-detail-page")
      assert has_element?(show_view, ".project-detail-title", "Alpha Project")
    end

    test "opens and closes the create modal", %{conn: conn} do
      {:ok, view, _html} = live(conn, ~p"/projects")

      refute has_element?(view, "#project-create-modal")
      _ = view |> element("#new-project-button") |> render_click()
      assert has_element?(view, "#project-create-modal")

      _ = view |> element("#cancel-project-create") |> render_click()
      refute has_element?(view, "#project-create-modal")
    end

    test "validates unique name inside workspace", %{conn: conn, scope: scope} do
      _ = project_fixture(scope, %{name: "Agent Sandbox"})
      {:ok, view, _html} = live(conn, ~p"/projects")

      _ = view |> element("#new-project-button") |> render_click()

      _ =
        view
        |> form("#project-create-form", project: %{name: "agent sandbox", description: "dupe"})
        |> render_submit()

      assert has_element?(view, "#project-create-form .text-error", "has already been taken")
    end
  end
end
