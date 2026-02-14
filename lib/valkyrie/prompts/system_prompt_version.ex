defmodule Valkyrie.Prompts.SystemPromptVersion do
  use Ecto.Schema
  import Ecto.Changeset

  @primary_key {:id, :binary_id, autogenerate: true}
  @foreign_key_type :binary_id

  schema "system_prompt_versions" do
    field :version, :integer
    field :body, :string

    belongs_to :system_prompt, Valkyrie.Prompts.SystemPrompt
    belongs_to :created_by_user, Valkyrie.Accounts.User, type: :id

    timestamps(type: :utc_datetime_usec, updated_at: false)
  end

  def changeset(prompt_version, attrs) do
    prompt_version
    |> cast(attrs, [:system_prompt_id, :version, :body, :created_by_user_id])
    |> validate_required([:system_prompt_id, :version, :body, :created_by_user_id])
    |> unique_constraint([:system_prompt_id, :version])
  end
end
