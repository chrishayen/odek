# Requirement: "a library that records API calls at runtime and later replays them as test cases with mocked dependencies"

Records outbound and inbound traffic, groups it into cases, then renders replayable test fixtures and dependency mocks.

std
  std.fs
    std.fs.write_all
      @ (path: string, content: string) -> result[void, string]
      + creates or overwrites the file
      # filesystem
  std.json
    std.json.encode_value
      @ (value: json_value) -> string
      + serializes a json value to text
      # serialization
  std.time
    std.time.now_millis
      @ () -> i64
      + returns current unix time in milliseconds
      # time
  std.hash
    std.hash.sha256_hex
      @ (data: bytes) -> string
      + returns a lowercase hex sha-256 digest
      # hashing

keploy
  keploy.record_inbound
    @ (rec: recorder_state, req: http_request, resp: http_response) -> recorder_state
    + appends an inbound request/response pair with a timestamp
    # recording
    -> std.time.now_millis
  keploy.record_outbound
    @ (rec: recorder_state, dep: string, req: http_request, resp: http_response) -> recorder_state
    + appends an outbound dependency call keyed on dep
    # recording
    -> std.time.now_millis
  keploy.case_id
    @ (inbound: http_request) -> string
    + returns a stable id for an inbound request so repeated samples merge into one case
    # identity
    -> std.hash.sha256_hex
  keploy.group_into_cases
    @ (rec: recorder_state) -> list[test_case]
    + returns one test_case per unique inbound id, collecting its outbound calls
    # grouping
  keploy.render_case_fixture
    @ (case: test_case) -> string
    + emits a serialized fixture with inbound request, expected response, and recorded outbound calls
    # rendering
    -> std.json.encode_value
  keploy.render_dependency_mock
    @ (case: test_case) -> string
    + emits a mock that answers each recorded outbound call with its recorded response
    # rendering
  keploy.write_suite
    @ (cases: list[test_case], out_dir: string) -> result[void, string]
    + writes one fixture and one mock file per case under out_dir
    - returns error on the first file write that fails
    # io
    -> std.fs.write_all
