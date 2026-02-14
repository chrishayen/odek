defmodule Valkyrie.Client.Runner do
  @moduledoc """
  Minimal runtime work loop built on the Req API client.
  """

  alias Valkyrie.Client.{API, Config}

  def run_once(%Config{} = config) do
    client = API.new(config)

    with {:ok, %{"stories" => [story | _]}} <-
           API.poll_stories(client, config.project_id, limit: 20),
         {:ok, _claim} <- API.claim_story(client, story["id"]),
         {:ok, %{"prompt" => _prompt}} <-
           API.resolve_prompt(
             client,
             story["project_id"],
             config.runtime_system_prompt_key,
             config.runtime_user_prompt
           ) do
      API.update_story_state(client, story["id"], config.runtime_done_state)
    else
      {:ok, %{"stories" => []}} -> {:ok, :no_work}
      other -> other
    end
  end
end
