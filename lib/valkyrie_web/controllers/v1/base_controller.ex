defmodule ValkyrieWeb.V1.BaseController do
  @moduledoc false

  defmacro __using__(_opts) do
    quote do
      use ValkyrieWeb, :controller
      import ValkyrieWeb.V1.BaseController
    end
  end

  def current_scope!(conn), do: conn.assigns.current_scope

  def render_error(conn, status, code, message) do
    conn
    |> Plug.Conn.put_status(status)
    |> Phoenix.Controller.json(%{error: %{code: code, message: message}})
  end

  def unauthorized(conn, message \\ "authentication required"),
    do: render_error(conn, :unauthorized, "unauthorized", message)

  def forbidden(conn, message \\ "insufficient permissions"),
    do: render_error(conn, :forbidden, "forbidden", message)

  def bad_request(conn, code, message), do: render_error(conn, :bad_request, code, message)
  def not_found(conn, message), do: render_error(conn, :not_found, "not_found", message)
  def conflict(conn, code, message), do: render_error(conn, :conflict, code, message)

  def validation_error(conn, code, changeset),
    do: render_error(conn, :bad_request, code, inspect(changeset.errors))

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
