defmodule Valkyrie.Repo do
  use Ecto.Repo,
    otp_app: :valkyrie,
    adapter: Ecto.Adapters.Postgres
end
