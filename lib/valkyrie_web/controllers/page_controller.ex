defmodule ValkyrieWeb.PageController do
  use ValkyrieWeb, :controller

  def home(conn, _params) do
    render(conn, :home)
  end
end
