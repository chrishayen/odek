# Requirement: "an IoT device automation library"

Stores devices, evaluates rules over their current state, and fires triggers when conditions match.

std: (all units exist)

iot_auto
  iot_auto.new_hub
    @ () -> hub_state
    + creates an empty hub with no devices or rules
    # construction
  iot_auto.register_device
    @ (hub: hub_state, device_id: string, kind: string) -> result[void, string]
    + adds a device of a given kind ("sensor", "switch", "thermostat") to the hub
    - returns error when the device id is already registered
    # device_registry
  iot_auto.update_reading
    @ (hub: hub_state, device_id: string, key: string, value: f64) -> result[void, string]
    + records a new reading for a device under the named key
    - returns error when the device is unknown
    # state_update
  iot_auto.get_reading
    @ (hub: hub_state, device_id: string, key: string) -> optional[f64]
    + returns the most recent reading for the key
    # state_query
  iot_auto.add_rule
    @ (hub: hub_state, name: string, condition: rule_condition, action: rule_action) -> result[string, string]
    + adds a rule that fires the action whenever the condition becomes true
    - returns error when any device referenced in the rule is unknown
    # rule_registry
  iot_auto.evaluate_rules
    @ (hub: hub_state) -> list[fired_action]
    + re-evaluates all rules and returns the actions that should fire this tick
    ? edge-triggered: a rule fires only on the false-to-true transition
    # rule_engine
  iot_auto.apply_action
    @ (hub: hub_state, action: rule_action) -> result[void, string]
    + updates the hub state to reflect a rule action (e.g. switch on/off)
    - returns error when the target device does not support the action
    # actuation
