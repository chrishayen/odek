defmodule Valkyrie.Workflow.StoryComment do
  use Ecto.Schema
  import Ecto.Changeset

  @primary_key {:id, :binary_id, autogenerate: true}
  @foreign_key_type :binary_id

  schema "story_comments" do
    field :body, :string

    belongs_to :story, Valkyrie.Workspace.Story
    belongs_to :author_user, Valkyrie.Accounts.User, type: :id

    timestamps(type: :utc_datetime_usec, updated_at: false)
  end

  def changeset(comment, attrs) do
    comment
    |> cast(attrs, [:story_id, :author_user_id, :body])
    |> validate_required([:story_id, :author_user_id, :body])
    |> validate_length(:body, min: 1)
  end
end
