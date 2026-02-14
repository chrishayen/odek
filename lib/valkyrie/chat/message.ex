defmodule Valkyrie.Chat.Message do
  use Ecto.Schema
  import Ecto.Changeset

  @primary_key {:id, :binary_id, autogenerate: true}
  @foreign_key_type :binary_id

  schema "chat_messages" do
    field :sender_kind, :string
    field :body, :string
    field :metadata, :map, default: %{}

    belongs_to :thread, Valkyrie.Chat.Thread
    belongs_to :sender_user, Valkyrie.Accounts.User, type: :id

    timestamps(type: :utc_datetime_usec, updated_at: false)
  end

  def changeset(message, attrs) do
    message
    |> cast(attrs, [:thread_id, :sender_kind, :sender_user_id, :body, :metadata])
    |> validate_required([:thread_id, :sender_kind, :body])
    |> validate_length(:body, min: 1)
  end
end
