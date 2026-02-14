defmodule ValkyrieWeb.HealthController do
  use ValkyrieWeb, :controller

  def show(conn, _params) do
    json(conn, %{status: "ok"})
  end
end
