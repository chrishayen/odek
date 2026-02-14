defmodule Valkyrie.Workflow.StoryClaim do
  use Ecto.Schema
  import Ecto.Changeset

  @primary_key {:id, :binary_id, autogenerate: true}
  @foreign_key_type :binary_id

  schema "story_claims" do
    field :released_at, :utc_datetime_usec

    belongs_to :story, Valkyrie.Workspace.Story
    belongs_to :claimed_by_user, Valkyrie.Accounts.User, type: :id

    timestamps(type: :utc_datetime_usec)
  end

  def changeset(claim, attrs) do
    claim
    |> cast(attrs, [:story_id, :claimed_by_user_id, :released_at])
    |> validate_required([:story_id, :claimed_by_user_id])
    |> unique_constraint(:story_id)
  end
end
