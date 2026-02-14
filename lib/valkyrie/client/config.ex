defmodule Valkyrie.Client.Config do
  @moduledoc """
  Runtime client configuration loaded from environment variables.
  """

  @enforce_keys [:base_url, :api_key, :project_id]
  defstruct base_url: nil,
            api_key: nil,
            project_id: nil,
            poll_interval_sec: 15,
            runtime_system_prompt_key: "worker",
            runtime_user_prompt: "",
            runtime_done_state: "review"

  def load_from_env do
    with {:ok, base_url} <- required_env("SESSIONDB_BASE_URL"),
         {:ok, api_key} <- required_env("SESSIONDB_API_KEY"),
         {:ok, project_id} <- required_env("SESSIONDB_PROJECT_ID") do
      {:ok,
       %__MODULE__{
         base_url: String.trim_trailing(base_url, "/"),
         api_key: api_key,
         project_id: project_id,
         poll_interval_sec: int_env("SESSIONDB_POLL_INTERVAL_SEC", 15),
         runtime_system_prompt_key:
           env_or_default("SESSIONDB_RUNTIME_SYSTEM_PROMPT_KEY", "worker"),
         runtime_user_prompt: env_or_default("SESSIONDB_RUNTIME_USER_PROMPT", ""),
         runtime_done_state: env_or_default("SESSIONDB_RUNTIME_WORKER_DONE_STATE", "review")
       }}
    end
  end

  defp required_env(key) do
    case System.get_env(key) do
      value when is_binary(value) and value != "" -> {:ok, value}
      _ -> {:error, {:missing_env, key}}
    end
  end

  defp env_or_default(key, default) do
    case System.get_env(key) do
      value when is_binary(value) and value != "" -> value
      _ -> default
    end
  end

  defp int_env(key, default) do
    case System.get_env(key) do
      value when is_binary(value) ->
        case Integer.parse(value) do
          {n, ""} when n > 0 -> n
          _ -> default
        end

      _ ->
        default
    end
  end
end
