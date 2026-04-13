# Requirement: "a deep reinforcement learning algorithms library"

Agents are opaque state handles; the project layer exposes construction, action selection, and update steps. Tensor and network primitives live in std.

std
  std.tensor
    std.tensor.zeros
      @ (shape: list[i32]) -> tensor
      + returns a tensor of the given shape filled with zeros
      # tensor
    std.tensor.from_list
      @ (values: list[f64]) -> tensor
      + returns a 1-D tensor holding the given values
      # tensor
    std.tensor.add
      @ (a: tensor, b: tensor) -> tensor
      + returns element-wise sum
      # tensor
    std.tensor.scale
      @ (t: tensor, factor: f64) -> tensor
      + returns t multiplied by factor
      # tensor
    std.tensor.argmax
      @ (t: tensor) -> i32
      + returns the index of the largest element
      # tensor
  std.nn
    std.nn.mlp
      @ (layer_sizes: list[i32]) -> network_handle
      + builds a fully-connected network with the given layer widths
      # network
    std.nn.forward
      @ (net: network_handle, input: tensor) -> tensor
      + returns the output tensor for the given input
      # network
    std.nn.backward
      @ (net: network_handle, loss: f64) -> network_handle
      + applies one gradient step against the accumulated loss
      # network
  std.random
    std.random.uniform
      @ () -> f64
      + returns a random f64 in [0, 1)
      # random

rl
  rl.new_replay_buffer
    @ (capacity: i32) -> replay_state
    + creates a fixed-capacity experience buffer
    # buffer
  rl.store
    @ (buffer: replay_state, obs: tensor, action: i32, reward: f64, next_obs: tensor, done: bool) -> replay_state
    + appends a transition, dropping the oldest when full
    # buffer
  rl.sample_batch
    @ (buffer: replay_state, batch_size: i32) -> list[transition]
    + returns a random batch of transitions
    - returns an empty list when the buffer has fewer items than batch_size
    # buffer
    -> std.random.uniform
  rl.new_dqn
    @ (obs_dim: i32, action_dim: i32, hidden: list[i32]) -> agent_state
    + constructs a DQN agent with a Q-network and target network
    # agent
    -> std.nn.mlp
  rl.new_ppo
    @ (obs_dim: i32, action_dim: i32, hidden: list[i32]) -> agent_state
    + constructs a PPO agent with policy and value networks
    # agent
    -> std.nn.mlp
  rl.select_action
    @ (agent: agent_state, obs: tensor, epsilon: f64) -> i32
    + returns the greedy action with probability 1 - epsilon
    + otherwise returns a uniformly random action
    # policy
    -> std.nn.forward
    -> std.tensor.argmax
    -> std.random.uniform
  rl.update
    @ (agent: agent_state, batch: list[transition], learning_rate: f64) -> agent_state
    + performs one optimization step against the batch
    # training
    -> std.nn.forward
    -> std.nn.backward
  rl.sync_target
    @ (agent: agent_state) -> agent_state
    + copies the online network weights into the target network
    # training
  rl.discount_returns
    @ (rewards: list[f64], gamma: f64) -> list[f64]
    + returns discounted cumulative returns in reverse order
    ? gamma in [0, 1]
    # advantage
  rl.advantage
    @ (rewards: list[f64], values: list[f64], gamma: f64, lambda_val: f64) -> list[f64]
    + returns GAE advantage estimates
    # advantage
  rl.save
    @ (agent: agent_state, path: string) -> result[void, string]
    + serializes the agent's parameters to disk
    - returns error on write failure
    # persistence
  rl.load
    @ (path: string) -> result[agent_state, string]
    + loads an agent from disk
    - returns error when the file is missing or malformed
    # persistence
