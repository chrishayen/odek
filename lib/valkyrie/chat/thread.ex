defmodule Valkyrie.Chat.Thread do
  use Ecto.Schema
  import Ecto.Changeset

  @primary_key {:id, :binary_id, autogenerate: true}
  @foreign_key_type :binary_id

  schema "chat_threads" do
    field :title, :string, default: ""
    field :deleted_at, :utc_datetime_usec

    belongs_to :organization, Valkyrie.Organizations.Organization
    belongs_to :project, Valkyrie.Workspace.Project
    belongs_to :story, Valkyrie.Workspace.Story
    belongs_to :created_by_user, Valkyrie.Accounts.User, type: :id

    timestamps(type: :utc_datetime_usec)
  end

  def changeset(thread, attrs) do
    thread
    |> cast(attrs, [
      :organization_id,
      :project_id,
      :story_id,
      :title,
      :created_by_user_id,
      :deleted_at
    ])
    |> validate_required([:organization_id, :project_id, :created_by_user_id])
  end
end
