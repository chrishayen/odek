defmodule ValkyrieWeb.UserSessionController do
  use ValkyrieWeb, :controller

  alias Valkyrie.Accounts
  alias ValkyrieWeb.UserAuth

  def create(conn, params) do
    log_in_with_password(conn, params["user"] || %{}, "Welcome back!")
  end

  def update_password(conn, %{"user" => user_params}) do
    user = conn.assigns.current_scope.user
    true = Accounts.sudo_mode?(user)
    {:ok, {_user, expired_tokens}} = Accounts.update_user_password(user, user_params)

    # disconnect all existing LiveViews with old sessions
    UserAuth.disconnect_sessions(expired_tokens)

    conn
    |> put_session(:user_return_to, ~p"/users/settings")
    |> log_in_with_password(
      %{"email" => user.email, "password" => user_params["password"]},
      "Password updated successfully!"
    )
  end

  def delete(conn, _params) do
    conn
    |> put_flash(:info, "Logged out successfully.")
    |> UserAuth.log_out_user()
  end

  defp log_in_with_password(conn, user_params, info) do
    email = Map.get(user_params, "email", "")
    password = Map.get(user_params, "password", "")

    if user = Accounts.get_user_by_email_and_password(email, password) do
      conn
      |> put_flash(:info, info)
      |> UserAuth.log_in_user(user, user_params)
    else
      # In order to prevent user enumeration attacks, don't disclose whether the email is registered.
      conn
      |> put_flash(:error, "Invalid email or password")
      |> put_flash(:email, String.slice(email, 0, 160))
      |> redirect(to: ~p"/users/log-in")
    end
  end
end
