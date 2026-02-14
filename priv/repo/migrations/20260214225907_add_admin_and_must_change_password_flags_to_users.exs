defmodule Valkyrie.Repo.Migrations.AddAdminAndMustChangePasswordFlagsToUsers do
  use Ecto.Migration

  def change do
    alter table(:users) do
      add :is_admin, :boolean, default: false, null: false
      add :must_change_password, :boolean, default: false, null: false
    end
  end
end
