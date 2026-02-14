defmodule Valkyrie.Prompts.RuleProfileItem do
  use Ecto.Schema
  import Ecto.Changeset

  @primary_key {:id, :binary_id, autogenerate: true}
  @foreign_key_type :binary_id

  schema "rule_profile_items" do
    belongs_to :profile, Valkyrie.Prompts.RuleProfile
    belongs_to :rule, Valkyrie.Prompts.Rule

    timestamps(type: :utc_datetime_usec, updated_at: false)
  end

  def changeset(item, attrs) do
    item
    |> cast(attrs, [:profile_id, :rule_id])
    |> validate_required([:profile_id, :rule_id])
    |> unique_constraint([:profile_id, :rule_id])
  end
end
