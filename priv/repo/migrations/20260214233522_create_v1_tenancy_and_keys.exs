defmodule Valkyrie.Repo.Migrations.CreateV1TenancyAndKeys do
  use Ecto.Migration

  def change do
    alter table(:users) do
      add :name, :string, default: "", null: false
    end

    create table(:organizations, primary_key: false) do
      add :id, :binary_id, primary_key: true
      add :name, :string, null: false
      add :deleted_at, :utc_datetime_usec

      timestamps(type: :utc_datetime_usec)
    end

    create unique_index(:organizations, [:name])

    create table(:organization_memberships, primary_key: false) do
      add :id, :binary_id, primary_key: true

      add :organization_id, references(:organizations, type: :binary_id, on_delete: :delete_all),
        null: false

      add :user_id, references(:users, on_delete: :delete_all), null: false
      add :role, :string, null: false

      timestamps(type: :utc_datetime_usec)
    end

    create unique_index(:organization_memberships, [:organization_id, :user_id])
    create index(:organization_memberships, [:user_id])

    create constraint(:organization_memberships, :organization_memberships_role_check,
             check: "role in ('owner','admin','member','viewer')"
           )

    create table(:api_keys, primary_key: false) do
      add :id, :binary_id, primary_key: true

      add :organization_id, references(:organizations, type: :binary_id, on_delete: :delete_all),
        null: false

      add :user_id, references(:users, on_delete: :delete_all), null: false
      add :key_prefix, :string, null: false
      add :key_hash, :binary, null: false
      add :revoked_at, :utc_datetime_usec

      timestamps(type: :utc_datetime_usec)
    end

    create unique_index(:api_keys, [:key_hash])
    create index(:api_keys, [:organization_id])
    create index(:api_keys, [:user_id])
    create index(:api_keys, [:organization_id, :revoked_at])
  end
end
