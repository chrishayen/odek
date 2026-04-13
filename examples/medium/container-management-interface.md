# Requirement: "a container management interface"

Lists, inspects, and controls containers and images through a pluggable container runtime client.

std: (all units exist)

container_manager
  container_manager.list_containers
    @ (client: runtime_client) -> result[list[container_info], string]
    + returns running and stopped containers with id, name, image, state
    - returns error when the runtime is unreachable
    # query
  container_manager.list_images
    @ (client: runtime_client) -> result[list[image_info], string]
    + returns local images with id, repo tags, and size
    - returns error when the runtime is unreachable
    # query
  container_manager.inspect_container
    @ (client: runtime_client, id: string) -> result[container_detail, string]
    + returns full config and runtime state for a container
    - returns error when id is unknown
    # inspection
  container_manager.start_container
    @ (client: runtime_client, id: string) -> result[void, string]
    + starts a stopped container
    - returns error when id is unknown
    - returns error when container is already running
    # lifecycle
  container_manager.stop_container
    @ (client: runtime_client, id: string, timeout_secs: i32) -> result[void, string]
    + sends stop and waits up to timeout_secs before forcing
    - returns error when id is unknown
    # lifecycle
  container_manager.remove_container
    @ (client: runtime_client, id: string, force: bool) -> result[void, string]
    + removes a container; force removes even when running
    - returns error when a running container is removed without force
    # lifecycle
  container_manager.stream_logs
    @ (client: runtime_client, id: string, on_line: log_callback) -> result[void, string]
    + invokes on_line for each log line until the container exits
    - returns error when id is unknown
    # logging
