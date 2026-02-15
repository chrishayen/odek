defmodule ValkyrieWeb.V1.BootstrapController do
  use ValkyrieWeb.V1.BaseController

  alias Valkyrie.Accounts
  alias Valkyrie.Organizations
  alias Valkyrie.Repo

  def create_organization(conn, params) do
    case Organizations.create_organization(params) do
      {:ok, org} ->
        conn
        |> put_status(:created)
        |> json(%{id: org.id, name: org.name})

      {:error, changeset} ->
        render_error(conn, :bad_request, "invalid_request", inspect(changeset.errors))
    end
  end

  def create_user(conn, params) do
    name = Map.get(params, "name", "")
    password = Map.get(params, "password", "")

    with false <- blank?(Map.get(params, "email")),
         false <- blank?(password),
         {:ok, user} <- Accounts.register_user(%{email: params["email"], name: name}),
         {:ok, {user, _tokens}} <- Accounts.update_user_password(user, %{password: password}),
         {:ok, user} <-
           user |> Ecto.Changeset.change(name: name, must_change_password: true) |> Repo.update(),
         {:ok, membership} <- Organizations.ensure_user_membership(user.id) do
      conn
      |> put_status(:created)
      |> json(%{
        id: user.id,
        email: user.email,
        name: user.name,
        organization_id: membership.organization_id,
        role: membership.role
      })
    else
      true ->
        render_error(conn, :bad_request, "invalid_request", "email and password are required")

      {:error, changeset} ->
        render_error(conn, :bad_request, "invalid_request", inspect(changeset.errors))
    end
  end

  def create_membership(conn, params) do
    case Organizations.create_membership(params) do
      {:ok, membership} ->
        conn
        |> put_status(:created)
        |> json(%{
          id: membership.id,
          organization_id: membership.organization_id,
          user_id: membership.user_id,
          role: membership.role
        })

      {:error, changeset} ->
        render_error(conn, :bad_request, "invalid_request", inspect(changeset.errors))
    end
  end
end
