defmodule Valkyrie.Prompts.SystemPrompt do
  use Ecto.Schema
  import Ecto.Changeset

  @primary_key {:id, :binary_id, autogenerate: true}
  @foreign_key_type :binary_id

  schema "system_prompts" do
    field :prompt_key, :string
    field :name, :string
    field :active_version, :integer

    belongs_to :organization, Valkyrie.Organizations.Organization

    timestamps(type: :utc_datetime_usec)
  end

  def changeset(prompt, attrs) do
    prompt
    |> cast(attrs, [:organization_id, :prompt_key, :name, :active_version])
    |> validate_required([:organization_id, :prompt_key, :name])
    |> unique_constraint([:organization_id, :prompt_key])
  end
end
