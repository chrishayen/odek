defmodule ValkyrieWeb.V1.StoryWorkflowTest do
  use ValkyrieWeb.ConnCase, async: true

  import Valkyrie.V1Fixtures

  test "poll, claim, state validation, and soft delete behavior", %{conn: conn} do
    %{user: user, organization: organization} = user_with_membership_fixture("member")

    key_conn =
      conn
      |> log_in_user(user)
      |> Plug.Conn.put_session(:active_organization_id, organization.id)
      |> post("/v1/api-keys", %{"organization_id" => organization.id})

    %{"api_key" => raw_key} = json_response(key_conn, 201)

    api_conn =
      build_conn()
      |> put_req_header("authorization", "Bearer #{raw_key}")

    project_conn = post(api_conn, "/v1/projects", %{"name" => "Project A"})
    %{"id" => project_id} = json_response(project_conn, 201)

    ready_conn =
      post(api_conn, "/v1/stories", %{
        "project_id" => project_id,
        "name" => "Ready story",
        "description" => "work item",
        "state" => "ready"
      })

    %{"id" => ready_story_id} = json_response(ready_conn, 201)

    backlog_conn =
      post(api_conn, "/v1/stories", %{
        "project_id" => project_id,
        "name" => "Backlog story",
        "description" => "later",
        "state" => "backlog"
      })

    %{"id" => backlog_story_id} = json_response(backlog_conn, 201)

    poll_conn = get(api_conn, "/v1/stories/poll", %{"project_id" => project_id, "limit" => 10})
    %{"stories" => stories} = json_response(poll_conn, 200)
    assert Enum.any?(stories, &(&1["id"] == ready_story_id))
    refute Enum.any?(stories, &(&1["id"] == backlog_story_id))

    claim_conn = post(api_conn, "/v1/stories/#{ready_story_id}/claim")
    assert %{"claimed" => true} = json_response(claim_conn, 200)

    conflict_conn = post(api_conn, "/v1/stories/#{ready_story_id}/claim")
    assert %{"error" => %{"code" => "claim_conflict"}} = json_response(conflict_conn, 409)

    invalid_state_conn =
      post(api_conn, "/v1/stories/#{ready_story_id}/state", %{"state" => "blocked"})

    assert %{"error" => %{"code" => "invalid_state"}} = json_response(invalid_state_conn, 400)

    delete_conn = delete(api_conn, "/v1/stories/#{backlog_story_id}")
    assert %{"deleted" => true} = json_response(delete_conn, 200)

    list_conn = get(api_conn, "/v1/stories", %{"project_id" => project_id})
    %{"stories" => listed_stories} = json_response(list_conn, 200)

    refute Enum.any?(listed_stories, &(&1["id"] == backlog_story_id))
  end
end
