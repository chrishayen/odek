# Requirement: "an abstraction layer that deploys applications to multiple cloud providers through a unified interface"

The library models a deployment plan and executes it against a pluggable provider driver. Providers are injected; the library itself does not talk to any specific API.

std: (all units exist)

deploy
  deploy.app_spec
    @ (name: string, image: string, replicas: i32, env: map[string, string]) -> app_spec
    + constructs a normalized application spec
    ? image is an opaque reference; the provider interprets its format
    # spec
  deploy.register_provider
    @ (name: string, apply_fn: provider_apply_fn, destroy_fn: provider_destroy_fn, status_fn: provider_status_fn) -> result[void, string]
    + registers a named provider driver
    - returns error when a provider with the same name is already registered
    # registry
  deploy.plan
    @ (current: list[app_spec], desired: list[app_spec]) -> deployment_plan
    + returns a plan containing create, update, and delete actions
    + apps are matched by name between current and desired lists
    # planning
  deploy.apply
    @ (provider_name: string, plan: deployment_plan) -> result[deployment_result, string]
    + executes the plan against the named provider, stopping on first error
    - returns error when the provider is unknown
    - returns a partial result when one action fails, listing completed actions
    # execution
    -> deploy.register_provider
  deploy.destroy
    @ (provider_name: string, app_name: string) -> result[void, string]
    + tears down a previously deployed application
    - returns error when the provider is unknown
    # execution
    -> deploy.register_provider
  deploy.status
    @ (provider_name: string, app_name: string) -> result[app_status, string]
    + queries the provider for the current status of an app
    - returns error when the provider is unknown
    # observation
    -> deploy.register_provider
