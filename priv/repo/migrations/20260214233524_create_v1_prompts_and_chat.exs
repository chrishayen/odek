defmodule Valkyrie.Repo.Migrations.CreateV1PromptsAndChat do
  use Ecto.Migration

  def change do
    create table(:skills, primary_key: false) do
      add :id, :binary_id, primary_key: true

      add :organization_id, references(:organizations, type: :binary_id, on_delete: :delete_all),
        null: false

      add :owner_user_id, references(:users, on_delete: :delete_all), null: false
      add :name, :string, null: false
      add :body, :string, default: "", null: false
      add :deleted_at, :utc_datetime_usec

      timestamps(type: :utc_datetime_usec)
    end

    create index(:skills, [:organization_id, :owner_user_id])

    create table(:rules, primary_key: false) do
      add :id, :binary_id, primary_key: true

      add :organization_id, references(:organizations, type: :binary_id, on_delete: :delete_all),
        null: false

      add :owner_user_id, references(:users, on_delete: :delete_all), null: false
      add :name, :string, null: false
      add :body, :string, default: "", null: false
      add :deleted_at, :utc_datetime_usec

      timestamps(type: :utc_datetime_usec)
    end

    create index(:rules, [:organization_id, :owner_user_id])

    create table(:skill_profiles, primary_key: false) do
      add :id, :binary_id, primary_key: true

      add :organization_id, references(:organizations, type: :binary_id, on_delete: :delete_all),
        null: false

      add :owner_user_id, references(:users, on_delete: :delete_all), null: false
      add :name, :string, null: false
      add :deleted_at, :utc_datetime_usec

      timestamps(type: :utc_datetime_usec)
    end

    create index(:skill_profiles, [:organization_id, :owner_user_id])

    create table(:rule_profiles, primary_key: false) do
      add :id, :binary_id, primary_key: true

      add :organization_id, references(:organizations, type: :binary_id, on_delete: :delete_all),
        null: false

      add :owner_user_id, references(:users, on_delete: :delete_all), null: false
      add :name, :string, null: false
      add :deleted_at, :utc_datetime_usec

      timestamps(type: :utc_datetime_usec)
    end

    create index(:rule_profiles, [:organization_id, :owner_user_id])

    create table(:skill_profile_items, primary_key: false) do
      add :id, :binary_id, primary_key: true

      add :profile_id, references(:skill_profiles, type: :binary_id, on_delete: :delete_all),
        null: false

      add :skill_id, references(:skills, type: :binary_id, on_delete: :delete_all), null: false

      timestamps(type: :utc_datetime_usec, updated_at: false)
    end

    create unique_index(:skill_profile_items, [:profile_id, :skill_id])

    create table(:rule_profile_items, primary_key: false) do
      add :id, :binary_id, primary_key: true

      add :profile_id, references(:rule_profiles, type: :binary_id, on_delete: :delete_all),
        null: false

      add :rule_id, references(:rules, type: :binary_id, on_delete: :delete_all), null: false

      timestamps(type: :utc_datetime_usec, updated_at: false)
    end

    create unique_index(:rule_profile_items, [:profile_id, :rule_id])

    create table(:project_skill_profiles, primary_key: false) do
      add :id, :binary_id, primary_key: true

      add :project_id, references(:projects, type: :binary_id, on_delete: :delete_all),
        null: false

      add :profile_id, references(:skill_profiles, type: :binary_id, on_delete: :delete_all),
        null: false

      timestamps(type: :utc_datetime_usec, updated_at: false)
    end

    create unique_index(:project_skill_profiles, [:project_id, :profile_id])

    create table(:project_rule_profiles, primary_key: false) do
      add :id, :binary_id, primary_key: true

      add :project_id, references(:projects, type: :binary_id, on_delete: :delete_all),
        null: false

      add :profile_id, references(:rule_profiles, type: :binary_id, on_delete: :delete_all),
        null: false

      timestamps(type: :utc_datetime_usec, updated_at: false)
    end

    create unique_index(:project_rule_profiles, [:project_id, :profile_id])

    create table(:system_prompts, primary_key: false) do
      add :id, :binary_id, primary_key: true

      add :organization_id, references(:organizations, type: :binary_id, on_delete: :delete_all),
        null: false

      add :prompt_key, :string, null: false
      add :name, :string, null: false
      add :active_version, :integer

      timestamps(type: :utc_datetime_usec)
    end

    create unique_index(:system_prompts, [:organization_id, :prompt_key])

    create table(:system_prompt_versions, primary_key: false) do
      add :id, :binary_id, primary_key: true

      add :system_prompt_id,
          references(:system_prompts, type: :binary_id, on_delete: :delete_all),
          null: false

      add :version, :integer, null: false
      add :body, :string, null: false
      add :created_by_user_id, references(:users, on_delete: :delete_all), null: false

      timestamps(type: :utc_datetime_usec, updated_at: false)
    end

    create unique_index(:system_prompt_versions, [:system_prompt_id, :version])

    create table(:chat_threads, primary_key: false) do
      add :id, :binary_id, primary_key: true

      add :organization_id, references(:organizations, type: :binary_id, on_delete: :delete_all),
        null: false

      add :project_id, references(:projects, type: :binary_id, on_delete: :delete_all),
        null: false

      add :story_id, references(:stories, type: :binary_id, on_delete: :nilify_all)
      add :title, :string, default: "", null: false
      add :created_by_user_id, references(:users, on_delete: :delete_all), null: false
      add :deleted_at, :utc_datetime_usec

      timestamps(type: :utc_datetime_usec)
    end

    create index(:chat_threads, [:organization_id, :project_id])

    create table(:chat_messages, primary_key: false) do
      add :id, :binary_id, primary_key: true

      add :thread_id, references(:chat_threads, type: :binary_id, on_delete: :delete_all),
        null: false

      add :sender_kind, :string, null: false
      add :sender_user_id, references(:users, on_delete: :nilify_all)
      add :body, :string, null: false
      add :metadata, :map, default: %{}, null: false

      timestamps(type: :utc_datetime_usec, updated_at: false)
    end

    create index(:chat_messages, [:thread_id, :inserted_at])
  end
end
