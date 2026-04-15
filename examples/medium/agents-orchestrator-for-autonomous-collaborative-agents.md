# Requirement: "a framework for orchestrating collaborating autonomous agents"

Agents have a role and a task queue; a crew assigns tasks and collects results. The language model call is a thin std primitive so tests can stub it.

std
  std.llm
    std.llm.complete
      fn (prompt: string) -> result[string, string]
      + returns the completion text for the prompt
      - returns error when the backend rejects the request
      # llm

agents
  agents.new_agent
    fn (role: string, goal: string) -> agent_state
    + returns an agent with the given role description and goal
    # construction
  agents.new_crew
    fn () -> crew_state
    + returns an empty crew with no members or pending tasks
    # construction
  agents.add_member
    fn (crew: crew_state, agent: agent_state) -> crew_state
    + appends the agent to the crew membership
    # membership
  agents.assign_task
    fn (crew: crew_state, task: string, role: string) -> result[crew_state, string]
    + queues the task for the first agent whose role matches
    - returns error when no member has the requested role
    # scheduling
  agents.run_next
    fn (crew: crew_state) -> result[tuple[crew_state, string], string]
    + runs the head task by prompting the assigned agent and returns its output
    - returns error when the task queue is empty
    # execution
    -> std.llm.complete
  agents.run_all
    fn (crew: crew_state) -> result[tuple[crew_state, list[string]], string]
    + drains the task queue in order and returns every result
    - returns error on the first task that fails
    # execution
    -> std.llm.complete
