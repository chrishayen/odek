# Requirement: "a library exposing the core actions of a filesystem mounter"

GUI is the caller's concern; this library exposes a single callable to request a mount and one to release it.

std: (all units exist)

fs_mounter
  fs_mounter.mount
    @ (source: string, mount_point: string) -> result[void, string]
    + attaches the source filesystem at the given mount point
    - returns error when the source cannot be mounted
    # mounting
  fs_mounter.unmount
    @ (mount_point: string) -> result[void, string]
    + detaches a previously mounted filesystem
    - returns error when nothing is mounted at that point
    # unmounting
