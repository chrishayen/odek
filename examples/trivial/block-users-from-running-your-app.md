# Requirement: "a guard that detects when the current process is running as root"

A single check the caller can use to refuse elevated execution.

std
  std.os
    std.os.effective_uid
      @ () -> i32
      + returns the effective user id of the current process
      # os

sudo_block
  sudo_block.is_root
    @ () -> bool
    + returns true when the effective user id is 0
    - returns false for any non-zero uid
    # privilege_check
    -> std.os.effective_uid
