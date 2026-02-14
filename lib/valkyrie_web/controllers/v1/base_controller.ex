defmodule ValkyrieWeb.V1.BaseController do
  @moduledoc false

  alias Valkyrie.Authorization

  defmacro __using__(_opts) do
    quote do
      use ValkyrieWeb, :controller
      import ValkyrieWeb.V1.BaseController
    end
  end

  def current_scope!(conn), do: conn.assigns.current_scope

  def authorize(conn, permission) do
    role = conn.assigns[:current_scope] && conn.assigns.current_scope.role

    if is_binary(role) and Authorization.can?(role, permission) do
      :ok
    else
      {:error, render_error(conn, :forbidden, "forbidden", "insufficient permissions")}
    end
  end

  def render_error(conn, status, code, message) do
    conn
    |> Plug.Conn.put_status(status)
    |> Phoenix.Controller.json(%{error: %{code: code, message: message}})
  end

  def parse_limit(params, fallback \\ 20) do
    case Map.get(params, "limit") do
      nil -> fallback
      value -> parse_int(value, fallback)
    end
  end

  def parse_int(value, _fallback) when is_integer(value) and value > 0, do: value

  def parse_int(value, fallback) when is_binary(value) do
    case Integer.parse(value) do
      {n, ""} when n > 0 -> n
      _ -> fallback
    end
  end

  def parse_int(_value, fallback), do: fallback

  def blank?(nil), do: true
  def blank?(""), do: true
  def blank?(_), do: false
end
