defmodule Valkyrie.Prompts.SkillProfileItem do
  use Ecto.Schema
  import Ecto.Changeset

  @primary_key {:id, :binary_id, autogenerate: true}
  @foreign_key_type :binary_id

  schema "skill_profile_items" do
    belongs_to :profile, Valkyrie.Prompts.SkillProfile
    belongs_to :skill, Valkyrie.Prompts.Skill

    timestamps(type: :utc_datetime_usec, updated_at: false)
  end

  def changeset(item, attrs) do
    item
    |> cast(attrs, [:profile_id, :skill_id])
    |> validate_required([:profile_id, :skill_id])
    |> unique_constraint([:profile_id, :skill_id])
  end
end
