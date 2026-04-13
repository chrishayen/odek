# Requirement: "a video platform library where membership and content decisions are governed by users through proposals and votes"

Core backend: accounts, video catalog with owners, and a governance module where members file proposals and votes decide outcomes.

std
  std.id
    std.id.new_uuid
      @ () -> string
      + returns a random v4 UUID string
      # identifiers
  std.time
    std.time.now_seconds
      @ () -> i64
      + returns current unix time in seconds
      # time
  std.hash
    std.hash.sha256_hex
      @ (data: bytes) -> string
      + returns the lowercase hex SHA-256 of data
      # hashing

video_platform
  video_platform.new
    @ () -> platform_state
    + creates an empty platform with no members, videos, or proposals
    # construction
  video_platform.enroll_member
    @ (state: platform_state, handle: string) -> result[tuple[member_id, platform_state], string]
    + creates a new member with zero voting weight
    - returns error when the handle is already enrolled
    # membership
    -> std.id.new_uuid
  video_platform.grant_voting_weight
    @ (state: platform_state, member: member_id, weight: i64) -> result[platform_state, string]
    + sets the voting weight for a member
    - returns error when weight is negative
    # membership
  video_platform.upload_video
    @ (state: platform_state, owner: member_id, title: string, content_hash: string) -> result[tuple[video_id, platform_state], string]
    + records a video owned by the member with the given content hash
    - returns error when the owner is not enrolled
    # catalog
    -> std.id.new_uuid
  video_platform.compute_content_hash
    @ (data: bytes) -> string
    + returns the canonical hex hash of video content for deduplication
    # catalog
    -> std.hash.sha256_hex
  video_platform.list_videos
    @ (state: platform_state) -> list[video]
    + returns all cataloged videos in upload order
    # catalog
  video_platform.file_proposal
    @ (state: platform_state, proposer: member_id, body: string, deadline: i64) -> result[tuple[proposal_id, platform_state], string]
    + creates an open proposal with the given deadline in unix seconds
    - returns error when the proposer is not a member
    - returns error when the deadline is not in the future
    # governance
    -> std.id.new_uuid
    -> std.time.now_seconds
  video_platform.cast_vote
    @ (state: platform_state, voter: member_id, proposal: proposal_id, approve: bool) -> result[platform_state, string]
    + records a vote weighted by the voter's voting weight
    - returns error when the proposal is closed
    - returns error when the voter has already voted on this proposal
    # governance
  video_platform.tally
    @ (state: platform_state, proposal: proposal_id) -> result[tally, string]
    + returns the current approve and reject totals
    - returns error when the proposal does not exist
    # governance
  video_platform.close_proposal
    @ (state: platform_state, proposal: proposal_id) -> result[tuple[proposal_outcome, platform_state], string]
    + closes the proposal after its deadline and returns approved or rejected
    - returns error when the deadline has not yet passed
    # governance
    -> std.time.now_seconds
  video_platform.takedown_video
    @ (state: platform_state, video: video_id, proposal: proposal_id) -> result[platform_state, string]
    + removes a video from the catalog when the referenced proposal was approved
    - returns error when the proposal was not approved
    # moderation
