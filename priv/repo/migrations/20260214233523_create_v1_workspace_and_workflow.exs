defmodule Valkyrie.Repo.Migrations.CreateV1WorkspaceAndWorkflow do
  use Ecto.Migration

  def change do
    create table(:projects, primary_key: false) do
      add :id, :binary_id, primary_key: true

      add :organization_id, references(:organizations, type: :binary_id, on_delete: :delete_all),
        null: false

      add :name, :string, null: false
      add :definition_of_done, :string, default: "", null: false
      add :deleted_at, :utc_datetime_usec

      timestamps(type: :utc_datetime_usec)
    end

    create index(:projects, [:organization_id])

    create table(:features, primary_key: false) do
      add :id, :binary_id, primary_key: true

      add :project_id, references(:projects, type: :binary_id, on_delete: :delete_all),
        null: false

      add :name, :string, null: false
      add :description, :string, default: "", null: false
      add :deleted_at, :utc_datetime_usec

      timestamps(type: :utc_datetime_usec)
    end

    create index(:features, [:project_id])

    create table(:stories, primary_key: false) do
      add :id, :binary_id, primary_key: true

      add :project_id, references(:projects, type: :binary_id, on_delete: :delete_all),
        null: false

      add :feature_id,
          references(:features, type: :binary_id, on_delete: :nilify_all)

      add :name, :string, null: false
      add :description, :string, default: "", null: false
      add :state, :string, default: "backlog", null: false
      add :deleted_at, :utc_datetime_usec

      timestamps(type: :utc_datetime_usec)
    end

    create index(:stories, [:project_id])
    create index(:stories, [:project_id, :inserted_at])
    create index(:stories, [:project_id, :state])

    create constraint(:stories, :stories_state_check,
             check: "state in ('backlog','ready','in_progress','review','done')"
           )

    create table(:story_claims, primary_key: false) do
      add :id, :binary_id, primary_key: true
      add :story_id, references(:stories, type: :binary_id, on_delete: :delete_all), null: false
      add :claimed_by_user_id, references(:users, on_delete: :delete_all), null: false
      add :released_at, :utc_datetime_usec

      timestamps(type: :utc_datetime_usec)
    end

    create unique_index(:story_claims, [:story_id], where: "released_at is null")

    create table(:story_comments, primary_key: false) do
      add :id, :binary_id, primary_key: true
      add :story_id, references(:stories, type: :binary_id, on_delete: :delete_all), null: false
      add :author_user_id, references(:users, on_delete: :delete_all), null: false
      add :body, :string, null: false

      timestamps(type: :utc_datetime_usec, updated_at: false)
    end

    create index(:story_comments, [:story_id, :inserted_at])

    create table(:story_history_events, primary_key: false) do
      add :id, :binary_id, primary_key: true
      add :story_id, references(:stories, type: :binary_id, on_delete: :delete_all), null: false
      add :actor_user_id, references(:users, on_delete: :delete_all), null: false
      add :event_type, :string, null: false
      add :payload, :map, null: false, default: %{}

      timestamps(type: :utc_datetime_usec, updated_at: false)
    end

    create index(:story_history_events, [:story_id, :inserted_at])
  end
end
