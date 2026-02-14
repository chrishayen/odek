defmodule Valkyrie.Prompts.ProjectRuleProfile do
  use Ecto.Schema
  import Ecto.Changeset

  @primary_key {:id, :binary_id, autogenerate: true}
  @foreign_key_type :binary_id

  schema "project_rule_profiles" do
    belongs_to :project, Valkyrie.Workspace.Project
    belongs_to :profile, Valkyrie.Prompts.RuleProfile

    timestamps(type: :utc_datetime_usec, updated_at: false)
  end

  def changeset(link, attrs) do
    link
    |> cast(attrs, [:project_id, :profile_id])
    |> validate_required([:project_id, :profile_id])
    |> unique_constraint([:project_id, :profile_id])
  end
end
