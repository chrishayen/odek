defmodule ValkyrieWeb.UserLive.LoginTest do
  use ValkyrieWeb.ConnCase, async: true

  import Phoenix.LiveViewTest
  import Valkyrie.AccountsFixtures

  describe "login page" do
    test "renders password-only login page", %{conn: conn} do
      {:ok, lv, _html} = live(conn, ~p"/users/log-in")

      assert has_element?(lv, "#login-page")
      assert has_element?(lv, ".login-left-panel")
      assert has_element?(lv, ".login-right-panel")
      assert has_element?(lv, "#login_form_password")
      assert has_element?(lv, "#login_form_password input[name='user[email]']")
      assert has_element?(lv, "#login_form_password input[name='user[password]']")
      refute has_element?(lv, "#login_form_magic")
      refute has_element?(lv, "a[href='/users/register']")
    end
  end

  describe "user login - password" do
    test "redirects if user logs in with valid credentials", %{conn: conn} do
      user = user_fixture() |> set_password()

      {:ok, lv, _html} = live(conn, ~p"/users/log-in")

      form =
        form(lv, "#login_form_password",
          user: %{email: user.email, password: valid_user_password()}
        )

      conn = submit_form(form, conn)

      assert redirected_to(conn) == ~p"/"
    end

    test "redirects to login page with a flash error if credentials are invalid", %{
      conn: conn
    } do
      {:ok, lv, _html} = live(conn, ~p"/users/log-in")

      form =
        form(lv, "#login_form_password", user: %{email: "test@email.com", password: "123456"})

      render_submit(form)

      conn = follow_trigger_action(form, conn)
      assert Phoenix.Flash.get(conn.assigns.flash, :error) == "Invalid email or password"
      assert redirected_to(conn) == ~p"/users/log-in"
    end
  end

  describe "re-authentication (sudo mode)" do
    setup %{conn: conn} do
      user = user_fixture()
      %{user: user, conn: log_in_user(conn, user)}
    end

    test "shows login page with email filled in", %{conn: conn, user: user} do
      {:ok, lv, _html} = live(conn, ~p"/users/log-in")

      assert has_element?(lv, ".login-form-subtitle")

      assert has_element?(
               lv,
               "#login_form_password input[name='user[email]'][value='#{user.email}']"
             )

      refute has_element?(lv, "a[href='/users/register']")
    end
  end
end
