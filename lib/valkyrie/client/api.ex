defmodule Valkyrie.Client.API do
  @moduledoc """
  Req-based client for runtime API-key endpoints.
  """

  alias Valkyrie.Client.Config

  defstruct base_url: nil, api_key: nil

  def new(%Config{} = config) do
    %__MODULE__{base_url: config.base_url, api_key: config.api_key}
  end

  def poll_stories(%__MODULE__{} = client, project_id, opts \\ []) do
    limit = Keyword.get(opts, :limit, 20)
    cursor = Keyword.get(opts, :cursor)

    params =
      %{project_id: project_id, limit: limit}
      |> maybe_put(:cursor, cursor)

    request(client, method: :get, path: "/v1/stories/poll", params: params)
  end

  def claim_story(%__MODULE__{} = client, story_id) do
    request(client, method: :post, path: "/v1/stories/#{story_id}/claim", json: %{})
  end

  def resolve_prompt(%__MODULE__{} = client, project_id, prompt_key, user_prompt \\ nil) do
    request(
      client,
      method: :post,
      path: "/v1/projects/#{project_id}/prompt-preview",
      json: %{prompt_key: prompt_key, user_prompt: user_prompt}
    )
  end

  def update_story_state(%__MODULE__{} = client, story_id, state) do
    request(client, method: :post, path: "/v1/stories/#{story_id}/state", json: %{state: state})
  end

  def append_chat_message(%__MODULE__{} = client, thread_id, body, sender_kind \\ "runtime") do
    request(
      client,
      method: :post,
      path: "/v1/chat/threads/#{thread_id}/messages",
      json: %{body: body, sender_kind: sender_kind}
    )
  end

  defp request(%__MODULE__{} = client, opts) do
    method = Keyword.fetch!(opts, :method)
    path = Keyword.fetch!(opts, :path)
    params = Keyword.get(opts, :params, %{})
    json = Keyword.get(opts, :json)

    req =
      Req.new(
        base_url: client.base_url,
        url: path,
        method: method,
        params: params,
        auth: {:bearer, client.api_key},
        receive_timeout: 30_000
      )

    req = if is_nil(json), do: req, else: Req.merge(req, json: json)

    case Req.request(req) do
      {:ok, %Req.Response{status: status, body: body}} when status in 200..299 -> {:ok, body}
      {:ok, %Req.Response{status: status, body: body}} -> {:error, {:http_error, status, body}}
      {:error, reason} -> {:error, reason}
    end
  end

  defp maybe_put(map, _key, nil), do: map
  defp maybe_put(map, key, value), do: Map.put(map, key, value)
end
