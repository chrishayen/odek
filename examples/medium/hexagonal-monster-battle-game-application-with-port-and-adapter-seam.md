# Requirement: "a sample hexagonal-architecture application for a creature-battle game, exposing use cases over a port-and-adapter seam"

Small domain around a monster battle with a port for persistence and a port for randomness; the core use cases only depend on the ports.

std: (all units exist)

monster_battle
  monster_battle.new_monster
    fn (id: string, name: string, hp: i32, attack: i32, defense: i32) -> monster
    + returns a monster record with full hp
    # domain
  monster_battle.new_battle
    fn (player: monster, enemy: monster) -> battle_state
    + pairs the two monsters into a fresh battle with the player acting first
    # construction
  monster_battle.apply_attack
    fn (attacker: monster, defender: monster) -> monster
    + returns the defender with hp reduced by max(1, attacker.attack - defender.defense)
    + clamps hp at zero
    # combat
  monster_battle.is_defeated
    fn (m: monster) -> bool
    + returns true when hp is zero or below
    - returns false otherwise
    # combat
  monster_battle.player_turn
    fn (state: battle_state) -> battle_state
    + applies the player's attack to the enemy and advances the turn
    # use_case
  monster_battle.enemy_turn
    fn (state: battle_state, rng: rng_port) -> battle_state
    + picks an enemy action using the randomness port and applies it to the player
    # use_case
  monster_battle.save
    fn (state: battle_state, repo: battle_repo_port) -> result[void, string]
    + persists the battle via the repository port
    - returns error when the port reports failure
    # persistence
  monster_battle.load
    fn (id: string, repo: battle_repo_port) -> result[battle_state, string]
    + restores a battle by id via the repository port
    - returns error when the id is unknown or the port fails
    # persistence
