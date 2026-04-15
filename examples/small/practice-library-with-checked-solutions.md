# Requirement: "a library of small programming practice tasks with checked solutions"

A task registry and a grader. No UI, no runner; callers supply answers and receive a verdict.

std: (all units exist)

practice
  practice.new_catalog
    fn () -> task_catalog
    + creates an empty catalog
    # construction
  practice.add_task
    fn (catalog: task_catalog, id: string, prompt: string, expected: string) -> task_catalog
    + registers a task by id with its prompt and expected answer
    # catalog
  practice.get_prompt
    fn (catalog: task_catalog, id: string) -> optional[string]
    + returns the prompt for a task id if known
    - returns none for unknown ids
    # catalog
  practice.check
    fn (catalog: task_catalog, id: string, answer: string) -> result[bool, string]
    + returns true when the answer matches the expected value
    - returns error when the task id is unknown
    # grading
