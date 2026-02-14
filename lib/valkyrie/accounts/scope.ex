defmodule Valkyrie.Accounts.Scope do
  @moduledoc """
  Defines the scope of the caller to be used throughout the app.

  The `Valkyrie.Accounts.Scope` allows public interfaces to receive
  information about the caller, such as if the call is initiated from an
  end-user, and if so, which user. Additionally, such a scope can carry fields
  such as "super user" or other privileges for use as authorization, or to
  ensure specific code paths can only be access for a given scope.

  It is useful for logging as well as for scoping pubsub subscriptions and
  broadcasts when a caller subscribes to an interface or performs a particular
  action.

  Feel free to extend the fields on this struct to fit the needs of
  growing application requirements.
  """

  alias Valkyrie.Accounts.User

  defstruct user: nil,
            organization_id: nil,
            role: nil,
            api_key_id: nil,
            auth_mode: :session

  @doc """
  Creates a scope for the given user.

  Returns nil if no user is given.
  """
  def for_user(%User{} = user) do
    %__MODULE__{user: user, auth_mode: :session}
  end

  def for_user(nil), do: nil

  @doc """
  Creates a scope resolved from API key authentication.
  """
  def for_api_key(%User{} = user, organization_id, role, api_key_id) do
    %__MODULE__{
      user: user,
      organization_id: organization_id,
      role: role,
      api_key_id: api_key_id,
      auth_mode: :api_key
    }
  end
end
