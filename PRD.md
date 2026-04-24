# ODEK — Product Requirements Document
## "AI-Powered Software Decomposition into Rune Trees"

**Version:** 1.0
**Language:** Go 1.26.1+
**Target Time:** 3-hour hackathon
**Module:** `shotgun.dev/odek`

> **Known Issue:** `main.go` declares `var systemPrompt string` (never initialized) and passes the empty string to `decompose.DecomposeStructured()` when using the `-d` flag. The real system prompt lives in `internal/decomposer/decompose.md` (embedded via `//go:embed`). To fix the `-d` flag, load that embedded file or refactor `main.go` to use `decomposer.NewDecomposer` directly.

---

## 1. Product Overview

ODEK is a terminal application that uses an LLM (via an OpenAI-compatible API) to decompose software requirements into hierarchical function trees called **Runes**. The user describes a feature in a chat interface; the AI breaks it down into a `project_package` of feature-specific functions and a reusable `std_package` of generic library primitives. The result is browsed in a Miller-column (macOS Finder-style) TUI.

**Core flows:**
1. **Landing** → animated ODEK logo over scrolling kanji.
2. **Chat** → user describes a feature; AI discusses or decomposes it.
3. **Decomposition pane** → interactive column browser of the rune tree.
4. **Recursive expansion** → AI expands each rune into sub-runes level-by-level.

---

## 2. Architecture

```
shotgun.dev/odek
├── main.go                      # CLI entry: -p, -d, -j flags; launches TUI
├── go.mod                       # module + charm/bubbletea v2 deps
│
├── openai/                      # HTTP client for OpenAI-compatible APIs
│   ├── list_models.go           # Client constructor, ListModels, HealthCheck
│   └── chat.go                  # Chat, streaming, AskToolLoop, ToolHandler
│
├── decompose/                   # Lightweight structured-output helpers
│   ├── decompose.go             # Decompose() + DecomposeStructured()
│   └── types.go                 # ParseDecomposition, Validate, FormatJSON
│
├── internal/decomposer/         # CORE ENGINE
│   ├── decompose.md             # SYSTEM PROMPT (embedded via //go:embed)
│   ├── decomposer.go            # Tool schemas, multi-turn loop, parallel init, merge
│   ├── types.go                 # Session, Snapshot, tree state, RuneStatus
│   ├── expand.go                # BFS recursive expansion with event streaming
│   ├── events.go                # ExpansionEvent sum type
│   ├── config.go                # Config + ConfigForEffort(level)
│   └── normalize.go             # Strip "fn"/"@" from leaked signatures
│
├── internal/examples/           # Example corpus loader
│   └── examples.go              # LoadFromDir, Lookup, Manifest
│
├── internal/effort/             # Complexity estimator
│   └── effort.go                # Estimate via forced rate_effort tool call
│
├── internal/toollog/            # Telemetry
│   └── toollog.go               # JSONL logger for read_example calls
│
├── internal/tui/                # TERMINAL UI (Bubble Tea v2)
│   ├── landing.go               # Animated splash screen
│   ├── create_feature.go        # Chat page + decompose orchestration
│   ├── feature_decomp.go        # Miller-column rune browser + detail pane
│   ├── split_feature.go         # Split-pane wrapper (chat | runes)
│   ├── transition.go            # Slide/collapse page transition
│   ├── chat.go                  # Rich chat component (markdown, code, thinking)
│   ├── styles.go                # Colors, kanji pool, help bindings
│   └── decomp_render.go         # Render decomposition summary text
│
├── std/                         # STRUCTURED STDLIB RUNE REGISTRY
│   ├── concurrency/once/{do}/
│   ├── http/client_delete/{send}/
│   ├── http/client_options/{send}/
│   ├── http/client_patch/{send}/
│   └── regex/test/1.0.0.md
│
└── examples/                    # EXAMPLE CORPUS (~1,000 .md files)
    ├── trivial/*.md             # tiny programs (hello-world, add-two-integers)
    ├── small/*.md               # single-file libs
    ├── medium/*.md              # multi-module libs
    └── large/*.md               # subsystems / full stacks
```

---

## 3. External Dependencies (go.mod)

```go
module shotgun.dev/odek

go 1.22

require (
    charm.land/bubbles/v2 v2.1.0
    charm.land/bubbletea/v2 v2.0.4
    charm.land/lipgloss/v2 v2.0.3
    github.com/alecthomas/chroma/v2 v2.23.1
    github.com/charmbracelet/x/ansi v0.11.7
    github.com/lucasb-eyer/go-colorful v1.4.0
)
```

---

## 4. openai/ — API Client

### 4.1 Client (list_models.go)
```go
type Client struct{ baseURL, apiKey string; client *http.Client }
func NewClient(baseURL string, apiKey ...string) (*Client, error)
func (c *Client) ListModels(ctx context.Context) ([]ModelInfo, error)
func (c *Client) HealthCheck(ctx context.Context) error
```
- Base URL defaults to `http://127.0.0.1:1234`. Append `/v1` internally.
- Auth header `Authorization: Bearer <apiKey>` when key provided.

### 4.2 Chat Completions (chat.go)

**Core types:**
```go
type ChatMessage struct {
    Role             string
    Content          string
    ReasoningContent string // for streaming thinking blocks
    ToolCalls        []ToolCall
    ToolCallID       string
}
type Tool struct { Type string; Function *FunctionDefinition }
type FunctionDefinition struct { Name, Description string; Parameters any }
type ToolCall struct { ID string; Type string; Function ToolCallFunction }
type ToolCallFunction struct { Name, Arguments string }
type ChatCompletionRequest struct {
    Model, Messages, Temperature, MaxTokens, Tools, ToolChoice, Stream
}
type ChatCompletionResponse struct { Choices []Choice; Usage *Usage }
type Choice struct { Message ChatMessage; FinishReason string; Delta Delta }
type Delta struct { Content, ReasoningContent string }
type Usage struct { PromptTokens, CompletionTokens, TotalTokens int }
```

**Methods:**
```go
func (c *Client) Chat(ctx context.Context, req *ChatCompletionRequest) (*ChatCompletionResponse, error)
// If ctx carries a thinking callback OR req.Stream=true → uses SSE streaming path.
// Otherwise uses standard POST /chat/completions.

func (c *Client) Ask(ctx context.Context, systemPrompt, userMessage string) (string, error)
func (c *Client) AskMessages(ctx context.Context, messages []ChatMessage) (string, error)
func (c *Client) AskTool(ctx context.Context, systemPrompt, userMessage string, tool Tool) (string, error)
```

**AskToolLoop** — the heart of multi-turn tool usage:
```go
type ToolHandler func(ctx context.Context, call ToolCall) (result string, terminal bool, err error)

func (c *Client) AskToolLoop(
    ctx context.Context,
    messages []ChatMessage,
    tools []Tool,
    handler ToolHandler,
    maxIterations int,
    toolChoice any, // nil → "auto"
) (final ChatMessage, history []ChatMessage, err error)
```
- Loop posts history + tools to `/chat/completions`.
- If assistant message has no `ToolCalls`, return it.
- If it has tool calls, dispatch each to `handler` in order.
  - `terminal=false` → append tool result, continue loop.
  - `terminal=true` → append tool result, return assistant message.
- Exceeding `maxIterations` without terminal → error.
- Returned `history` always contains full accumulated messages (including on error).

**Streaming implementation details:**
- Set `Stream=true`, `Accept: text/event-stream`.
- Parse SSE lines prefixed with `data:`.
- Accumulate `Delta.Content` and `Delta.ReasoningContent`.
- Accumulate `Delta.ToolCalls` by `Index` (fragments arrive across chunks).
- Fire thinking callback on every `reasoning_content` delta.
- Reassemble into a single `ChatCompletionResponse` so callers never care which path ran.

**Context key for thinking:**
```go
func WithThinkingCallback(ctx context.Context, cb func(string)) context.Context
```

---

## 5. internal/decomposer/ — Decomposition Engine

### 5.1 System Prompt (decompose.md)

This file is embedded via `//go:embed decompose.md`. It is **the single most important artifact**.

**Full prompt text (exactly as shipped):**

```markdown
You are a software architect that decomposes requirements into composition trees using a stdlib-first strategy.

You are building **libraries** — packages of reusable, importable functions that consumers call from their own code. The output is never an executable, CLI, or binary. Do not generate main() functions, argument parsing, process lifecycle management (signal handling, graceful shutdown), or any other binary-level concerns. Every root node is a library entry point: a function with typed inputs and typed outputs.

# What is a composition tree?

A composition tree decomposes software into a hierarchy where the dot-separated path IS the structure:

- The **application** is composed of **slices** (vertical functional areas)
- Each **slice** is composed of **components** (modules within that area)
- Each **component** is composed of **units** (individual functions)

Every node is a real, testable piece of code. Parent nodes wire their children together. Leaf nodes are isolated functions with clear inputs and outputs.

# Stdlib-first strategy

When decomposing a requirement, your FIRST job is to identify what generic capabilities it needs and build those as standard library units (std.*). Your SECOND job is to decompose the feature into its own package tree that composes those stdlib units.

**The question to ask**: "What reusable library would I build so that this feature — and future features — can just compose it?"

**Two namespaces:**
- std.* — the standard library. Generic, reusable, feature-agnostic. This is where reusable functionality lives.
- project_name.* — the feature package. Fully decomposed into its own tree of components and units. References std.* units where appropriate via -> links.

# Complexity is bad

Decompose like a thoughtful engineer. Restraint is a feature. Specifically:

- **No invented variants.** If the user asks for "hello world", you produce ONE greeter, not `say_hello`, `say_hello_with_name`, and `say_hello_multiple_times`. Do not add parameters, variants, or scope the user did not request. The user gets exactly what they asked for.
- **No identity functions.** A rune like `format_world(s: string) -> s` is noise. Anything with no real body is not a rune — drop it.
- **No duplicate definitions.** Every rune name appears in exactly ONE place. If `std.io.print_string` lives in std, the project package never redefines `print_string`. It uses `-> std.io.print_string` to reference it. Listing the same unit in both packages is a bug, not a "reference".
- **No within-package synonyms.** `create_user` and `user_registration` are the same thing. Pick one and drop the other. Same for `authenticate_user` / `user_login`.

Thin general-purpose stdlib wrappers are GOOD. `std.io.print_string` is a legitimate std unit even though it's a thin wrapper around the language's print — it's reusable across projects and describes a generic capability. What matters is that std units are generic and reusable, not that they add code complexity. Thin is fine. Feature-specific filler is not.

Before emitting any rune, ask: "Would a thoughtful engineer actually write this function?" If the honest answer is no, drop it.

# Answer shape (for the decompose tool)

Your answer has three parts, which you encode as arguments to the `decompose` tool call — not as prose in the chat:

**Part 1: `summary`** — a 1-2 sentence narrative. On a fresh decomposition, describe what the feature is and the approach you took. On a refinement pass (when a prior decomposition is included above), describe what you changed in response to the user's latest feedback and why. Explain; do not list rune names.

**Part 2: `std_package` (new stdlib units)** — the std.* composition tree. Only include NEW std.* units not already in the existing stdlib. If all needed stdlib already exists, omit `std_package` entirely.

**Part 3: `project_package` (the feature)** — the feature package tree. Decompose it into components and units just like std. Where the feature uses a stdlib unit, show it inline as a child with the prefix "-> " to indicate a reference (not a new definition). These references do not need test cases since they are already tested in std.

The tree format below describes the conceptual shape of std_package and project_package. You do not print trees — you encode them as the `runes` maps inside each package argument to the `decompose` tool.

Tree shape:

    std
      std.slice
        std.slice.component
          std.slice.component.leaf_unit
            fn (input: type) -> return_type
            + unit test
            - unit failure
            ? default value chosen since unspecified
            # responsibility_label

    project_name
      project_name.slice
        -> std.some.unit
        project_name.slice.leaf_unit
          fn (input: type) -> return_type
          + test for this unit
          - failure case
          ? errors logged to stderr, not a file
          # responsibility_label

The trees above are for your reference only. Your actual output is the `decompose` tool call, whose JSON arguments encode these trees into `project_package` / `std_package` / `summary`.

# Type system

Use these types for signatures:
- Signed integers: i8, i16, i32, i64
- Unsigned integers: u8, u16, u32, u64
- Floating point: f32, f64
- Primitives: string, bool, bytes
- Collections: list[T], map[K, V]
- Nullable: optional[T]
- Fallible: result[T, E]
- No return value: void

Types can be nested: result[list[i32], string]

# Rules

1. STDLIB FIRST. Decompose generic capabilities before the feature. Ask: could another feature use this without modification? If yes → std.*
2. The tree structure IS the composition. Parent nodes compose their children. Indentation shows nesting.
3. Every leaf node MUST have a signature (fn line) and test cases (except -> references). Package (non-leaf) nodes are organizational groupings and MUST NOT have signatures or test cases. Include every meaningful positive (+) and negative (-) test case — not just one of each. Cover edge cases, boundary values, and error variants.
4. Every node SHOULD list assumptions (? lines) — decisions you made to fill gaps the user didn't specify. These are behaviors, defaults, scope choices, or strategies you chose on their behalf. The user will review these and refine. Examples: "default port 8080", "graceful shutdown with 5s timeout", "serves index.html at directory root", "plaintext logging, not JSON". Omit ? lines only when the requirement fully specifies the node's behavior.
5. Every leaf node MUST have a responsibility tag (# line) — a short snake_case label for the functional concern it serves (e.g., "input_handling", "rendering", "network_io"). Nodes that work together toward the same concern share the same label. Choose labels that describe what the user would evaluate as a group when reviewing completeness.
6. std.* test cases must be feature-agnostic. Never mention a feature name or use feature-specific example values in std tests. The FIRST positive test case becomes the rune's description — it must describe the generic capability.
   - BAD: `+ writes "hello world" to stdout` (uses content from the specific feature)
   - GOOD: `+ writes provided string to stdout` (describes the generic capability)
7. Do not emit nodes for constants, config values, or type definitions — only executable behavior.
8. Use canonical verbs in leaf names: create, read, update, delete, validate, send, resolve, parse, serve, listen, handle, shutdown, filter, sort, transform, log, hash, generate, verify, etc.
9. Normalize verb synonyms (e.g., "show" → "read", "remove" → "delete").
10. snake_case everything. Subjects are domain entities, not UI elements.
11. If existing units are provided as context, you have three options for each:
   - -> path.to.unit — reference it as-is (it already does what you need)
   - ~> path.to.unit — extend it (it partially does what you need; include only the NEW +/- test cases to add)
   - Define a new node — when nothing existing covers the capability
12. The output is always a **library**. Do not generate CLI entry points, main functions, argument parsing, or binary-level concerns (signal handling, process exit codes, graceful shutdown). The feature root node is a library entry point — a callable function with typed parameters and return values that consumers import and call.
13. The root nodes (std, project_name) are package containers — they MUST have at least one child unit. Do not put executable behavior directly on a root node. Decompose into the minimum necessary child units — avoid unnecessary nesting depth.
14. NO DUPLICATION. Every rune name appears in exactly ONE package. If a unit lives in std, the project package does NOT redefine it — use `-> std.path.to.unit` to reference. Same applies within a package: pick one name per function, never two synonyms for the same thing.
15. NO INVENTED SCOPE. Decompose exactly what the user asked for. If they asked for "hello world", do not also invent `say_hello_with_name` or `say_hello_multiple_times`. If they asked for "add two integers", do not also invent `add_signed` and `add_unsigned`. The user gets what they specified — no more, no less.
16. The `function_signature` JSON field in the `decompose` tool call contains ONLY the type signature — e.g. `(a: i32, b: i32) -> i32` or `() -> string`. The `fn` marker shown in the tree format above is a visual marker for the human reader of this prompt and the example corpus; it is NOT part of the signature string. Do not include `fn`, `@`, or any other prefix in `function_signature`.

# Refinement

When refining a previous decomposition based on user feedback:
- If the feedback changes the feature's core behavior or purpose, update the feature root name and ALL child node names to match. Names must always reflect what the code actually does.
- Do not preserve old names when the behavior they describe has changed.
- Std unit names should also be updated if their behavior changes, but leave unchanged std units alone.

# How you deliver answers

You have two tools available on every decompose call:

1. **`read_example(paths)`** — load the full text of one or more example decompositions from the corpus. Pass a list of handles (each handle is `tier/slug` as listed in the `Example index` section at the bottom of this system message). Reading 2–5 handles in one call is normal.

2. **`decompose(project_package, std_package?)`** — submit your final answer. Call this last, exactly once per requirement.

### Typical flow

1. Read the user's requirement.
2. Scan the `Example index` section at the end of this system message and pick the handles that look most relevant to what the requirement asks for.
3. Call `read_example` with those handles (pass 2–5 in one call). Read the returned files carefully.
4. If the first batch doesn't cover the concepts you need, call `read_example` again with different handles.
5. Design your decomposition by matching the style of the examples — the same level of restraint on trivial requirements, the same substantive std primitives when real subsystems are needed, the same rule that every rune appears in exactly one package.
6. Call `decompose` with your answer.

Never skip `read_example` on a non-trivial requirement. The corpus is how you anchor your taste to a known-good baseline.

# Shape anchors

These three anchors show the minimum viable format in case the example index is empty. They are NOT a replacement for calling `read_example` — always prefer fresh retrieval against the current corpus.

---

Requirement: "hello world"

A library that returns the canonical greeting. Printing is the caller's responsibility.

std: (all units exist)

greeter
  greeter.greet
    fn () -> string
    + returns the string "Hello, world!"
    ? greeting is hardcoded; no parameter
    # greeting

---

Requirement: "a JWT signer and verifier (HS256)"

The project layer is just two entry points; the real work lives in std primitives (encoding, HMAC, JSON, time) that any crypto/auth project would reuse.

std
  std.encoding
    std.encoding.base64url_encode
      fn (data: bytes) -> string
      + encodes bytes to base64url without padding
      + returns "" for empty input
      # encoding
    std.encoding.base64url_decode
      fn (encoded: string) -> result[bytes, string]
      + decodes base64url with or without padding
      - returns error on characters outside the base64url alphabet
      # encoding
  std.crypto
    std.crypto.hmac_sha256
      fn (key: bytes, data: bytes) -> bytes
      + computes HMAC-SHA256 of data under key
      + returns 32 bytes
      # cryptography
  std.json
    std.json.parse_object
      fn (raw: string) -> result[map[string, string], string]
      + parses a JSON object into a string-to-string map
      - returns error on invalid JSON or non-object root
      # serialization
    std.json.encode_object
      fn (obj: map[string, string]) -> string
      + encodes a string-to-string map as JSON
      # serialization

jwt
  jwt.sign
    fn (payload: map[string, string], secret: string) -> result[string, string]
    + returns a JWT in "header.payload.signature" format
    - returns error when secret is empty
    # token_signing
    -> std.json.encode_object
    -> std.encoding.base64url_encode
    -> std.crypto.hmac_sha256
  jwt.verify
    fn (token: string, secret: string) -> result[map[string, string], string]
    + returns the payload when signature is valid
    - returns error when the token does not have exactly three segments
    - returns error when the signature does not match
    # token_verification
    -> std.encoding.base64url_decode
    -> std.crypto.hmac_sha256
    -> std.json.parse_object

---

Requirement: "a chat application backend"

Six project operations at the feature boundary; all substantive plumbing — bcrypt, JWT, WebSocket, SQL — lives in std as genuinely reusable subsystems. Notice that `bcrypt.hash`, `jwt.sign`, `sql.query`, etc. appear exactly once — only under std_package. The project runes mention them as dependencies in their descriptions.

std
  std.bcrypt
    std.bcrypt.hash
      fn (password: string) -> result[string, string]
      + returns a bcrypt hash of the password with a random salt
      - returns error when password exceeds 72 bytes
      # cryptography
    std.bcrypt.verify
      fn (password: string, hash: string) -> bool
      + returns true when the password matches the hash
      # cryptography
  std.jwt
    std.jwt.sign
      fn (payload: map[string, string], secret: string) -> result[string, string]
      + signs a payload with HS256 and returns a JWT
      # token_signing
    std.jwt.verify
      fn (token: string, secret: string) -> result[map[string, string], string]
      + verifies a JWT and returns its payload
      - returns error on bad signature or expired token
      # token_verification
  std.websocket
    std.websocket.upgrade
      fn (req: http_request) -> result[websocket_conn, string]
      + upgrades an HTTP request to a WebSocket connection
      # networking
    std.websocket.send
      fn (c: websocket_conn, data: bytes) -> result[void, string]
      + sends a binary or text frame on the connection
      # networking
  std.sql
    std.sql.query
      fn (conn: db_conn, sql: string, args: list[any]) -> result[list[row], string]
      + executes a parameterized query and returns the result rows
      - returns error on syntax error or constraint violation
      # persistence

chat
  chat.create_user
    fn (creds: credentials) -> result[user_id, string]
    + registers a new user with the password stored as a bcrypt hash
    - returns error when the username is already taken
    # account_management
    -> std.bcrypt.hash
    -> std.sql.query
  chat.authenticate
    fn (creds: credentials) -> result[session_token, string]
    + verifies the password and returns a signed session token
    - returns error on bad password or unknown user
    # account_management
    -> std.bcrypt.verify
    -> std.jwt.sign
    -> std.sql.query
  chat.create_room
    fn (token: session_token, name: string) -> result[room_id, string]
    + creates a new chat room owned by the authenticated user
    # room_management
    -> std.jwt.verify
    -> std.sql.query
  chat.join_room
    fn (token: session_token, room_id: room_id) -> result[void, string]
    + adds the authenticated caller to the room
    # room_management
    -> std.jwt.verify
    -> std.sql.query
  chat.send_message
    fn (token: session_token, room_id: room_id, body: string) -> result[message_id, string]
    + stores the message and broadcasts it to connected room members
    # messaging
    -> std.jwt.verify
    -> std.sql.query
    -> std.websocket.send
  chat.fetch_messages
    fn (token: session_token, room_id: room_id, since: i64) -> result[list[message], string]
    + returns messages posted to the room since the given unix timestamp
    # messaging
    -> std.jwt.verify
    -> std.sql.query

---

**Reminder — how to deliver your answer:**

Pick handles from the `Example index` at the end of this system message, call `read_example` with those handles to load their full text, then call `decompose` with your final answer as JSON: a `summary` string (1-2 sentence narrative), a `project_package` object, and when appropriate a `std_package` object.

**NEVER reply in plain text. NEVER paste the tree format into your message body. Your ONLY valid outputs are `read_example` tool calls and one final `decompose` tool call.** Trees in this prompt describe shape; they are not the answer format.

Now decompose the following requirement:
```

### 5.2 Tool Schemas (decomposer.go)

**`read_example` tool:**
```go
readExampleTool = openai.Tool{
    Type: openai.ToolTypeFunction,
    Function: &openai.FunctionDefinition{
        Name:        "read_example",
        Description: "Read the full contents of one or more example decompositions from the corpus. ...",
        Parameters: map[string]any{
            "type": "object",
            "properties": map[string]any{
                "paths": map[string]any{
                    "type": "array", "items": map[string]any{"type": "string"},
                    "minItems": 1, "maxItems": 5,
                    "description": "list of example handles to load, e.g. ['medium/csv-reader', 'trivial/hello-world']",
                },
            },
            "required": []string{"paths"},
        },
    },
}
```

**`decompose` tool:**
```go
runeSchema := map[string]any{
    "type": "object",
    "properties": map[string]any{
        "description":        map[string]any{"type": "string"},
        "function_signature": map[string]any{
            "type":        "string",
            "description": "Bare type signature only, e.g. '(a: i32, b: i32) -> result[i32, string]'. Do NOT include any marker prefix like 'fn' or '@'; those are visual markers in the tree format, not part of the value.",
        },
        "positive_tests":     map[string]any{"type": "array", "items": map[string]any{"type": "string"}},
        "negative_tests":     map[string]any{"type": "array", "items": map[string]any{"type": "string"}},
        "assumptions":        map[string]any{"type": "array", "items": map[string]any{"type": "string"}},
    },
    "required": []string{"description", "function_signature", "positive_tests", "negative_tests", "assumptions"},
}
packageSchema := map[string]any{
    "type": "object",
    "properties": map[string]any{
        "name":  map[string]any{"type": "string"},
        "runes": map[string]any{"type": "object", "additionalProperties": runeSchema},
    },
    "required": []string{"name", "runes"},
}
decomposeTool = openai.Tool{
    Type: openai.ToolTypeFunction,
    Function: &openai.FunctionDefinition{
        Name:        "decompose",
        Description: "Submit a rune decomposition. Provide a 1-2 sentence summary, a project_package, and optionally a std_package of reusable utilities.",
        Parameters: map[string]any{
            "type": "object",
            "properties": map[string]any{
                "summary": map[string]any{
                    "type":        "string",
                    "description": "A 1-2 sentence narrative shown to the user in the chat. On a fresh decomposition, describe what the feature is and the approach you took. On a refinement pass (when a prior decomposition is included), describe what you changed in response to the user's latest feedback and why. Explain, do not list rune names.",
                },
                "project_package": packageSchema,
                "std_package":     packageSchema,
            },
            "required": []string{"summary", "project_package"},
        },
    },
}
```

### 5.3 Data Types (types.go)

```go
type Rune struct {
    Description   string   `json:"description"`
    FunctionSig   string   `json:"function_signature"`
    PositiveTests []string `json:"positive_tests"`
    NegativeTests []string `json:"negative_tests"`
    Assumptions   []string `json:"assumptions"`
}
type PackageNode struct {
    Name  string          `json:"name"`
    Runes map[string]Rune `json:"runes"`
}
type DecompositionResponse struct {
    Summary        string       `json:"summary"`
    ProjectPackage PackageNode  `json:"project_package"`
    StdPackage     *PackageNode `json:"std_package,omitempty"`
}
type ClarificationRequest struct{ Message string }
type ClarificationNeeded struct{ Message string }
func (e *ClarificationNeeded) Error() string { return "clarification needed: " + e.Message }

type AutoDecomposition struct {
    Path       string
    Depth      int
    Response   *DecompositionResponse
    ParentPath string
    ChildPaths []string
}
```

### 5.4 Session State (types.go)

```go
type Session struct {
    Requirement  string
    EffortLevel  int
    EffortReason string
    Root         *AutoDecomposition
    BaseMessages []openai.ChatMessage
    // guarded by mu:
    tree       map[string]*AutoDecomposition
    treeOrder  []string
    status     map[string]RuneStatus
    totalRunes int
    maxDepth   int
    expanding  bool
    Events     <-chan ExpansionEvent
    Cancel     func()
}
```

**Snapshot (read-only render view):**
```go
type Snapshot struct {
    HasSession        bool
    PackageName       string
    TopLevelNames     []string
    RunesByName       map[string]Rune
    StatusByName      map[string]RuneStatus
    ChildrenByName    map[string][]string
    Requirement       string
    Summary           string
    TotalRunes        int
    MaxDepthReached   int
    Expanding         bool
    InFlightCount     int
    ErrorCount        int
    PackagePaths      []string
    RuneByPath        map[string]Rune
    DisplayNameByPath map[string]string
}
```

### 5.5 Expansion Events (events.go)

```go
type ExpansionEvent interface{ isExpansionEvent() }
type EventLevelStarted   struct{ Depth, Count int }
type EventRuneStarted    struct{ Path string; Depth int }
type EventRuneExpanded   struct{ Path, ParentPath string; Depth int; Response *DecompositionResponse; ElapsedMs int64; ChildCount int }
type EventRuneError      struct{ Path string; Depth int; Err string; ElapsedMs int64 }
type EventLevelCompleted struct{ Depth int; WallClockMs, SumRequestMs int64 }
type EventReadExample    struct{ Paths, Found []string }
type EventCapReached     struct{ TotalRunes, Cap int }
type EventCancelled      struct{}
type EventDone           struct{ TotalDecompositions, TotalRunes, MaxDepth int }
```

### 5.6 Config (config.go)

```go
type Config struct {
    ParallelInitial int
    MaxDepth        int
    RuneCap         int
    Recurse         bool
}
func ConfigForEffort(level int) Config {
    switch level {
    case 1: return Config{1, 0, 10, false}
    case 2: return Config{1, 10, 25, true}
    case 3: return Config{3, 10, 50, true}
    case 4: return Config{5, 10, 100, true}
    case 5: return Config{5, 10, 200, true}
    default: return Config{3, 10, 50, true}
    }
}
```

### 5.7 Decomposer Methods (decomposer.go)

```go
func NewDecomposer(api *openai.Client, examplesDir, toolLogPath string) (*Decomposer, error)
func (d *Decomposer) NewSession(ctx context.Context, req string, effortLevel int, effortReason string, cfg Config, sessCtx SessionContext) (*Session, error)
func (d *Decomposer) Decompose(ctx context.Context, messages []openai.ChatMessage, emit func(ExpansionEvent)) (any, []openai.ChatMessage, error)
func (d *Decomposer) MergeAttempts(ctx context.Context, req string, attempts []DecompositionResponse) (DecompositionResponse, []openai.ChatMessage, error)
func (d *Decomposer) ExpandStreaming(ctx context.Context, sess *Session, cfg Config) <-chan ExpansionEvent
```

**NewSession behavior:**
1. If `cfg.ParallelInitial <= 1`: single Decompose call.
2. If `cfg.ParallelInitial > 1`: run N concurrent initial decomposes, then merge successes via MergeAttempts.
3. On ClarificationRequest → return *ClarificationNeeded error.
4. On success → create Session with root AutoDecomposition{Path:"root", Depth:0}.

**Decompose (tool loop) behavior:**
1. Posts system prompt + user message + tools to API.
2. Model may call read_example any number of times (up to maxToolIterations=6).
3. Each read_example is resolved against the example index; results appended as tool messages.
4. Model finally calls decompose; arguments parsed into DecompositionResponse.
5. If model replies in plain text → retry once with forced tool call; if still plain text → return ClarificationRequest.
6. emit is called for every EventReadExample.

**ExpandStreaming behavior:**
1. BFS level-by-level expansion of every rune in the session tree.
2. Each rune gets an isolated prompt (exact text):
```go
fmt.Sprintf(`Forget the prior decomposition. Imagine you are seeing "%s" for the first time, in isolation, as a function you have to implement.

The user is browsing this decomposition as an interactive hierarchy: each rune is a column in a Miller-column (macOS-Finder-style) view, and the user will drill from parent to child to child. Your job is to continue that hierarchical breakdown by one more level beneath "%s".

Question: what 0–3 child units make up "%s"'s implementation? Each child should be a self-contained step the user would naturally drill into — a private helper, a distinct pipeline stage, or an internal subsystem, depending on the parent's granularity. They will appear as the next column to the right of "%s".

Call the decompose tool. The runes map keys must be of the form "%s.<child_name>". Example, for a different rune: if you were expanding "image.compress", reasonable children would be "image.compress.detect_format", "image.compress.choose_quality", "image.compress.encode_bytes". Each is a verb-phrase describing one internal step.

If "%s" is a single primitive operation (like an arithmetic op or a single syscall) and has no meaningful children, return an empty runes map ({}). That is the correct answer for leaves.

Hard rules:
- Reply ONLY by calling the decompose tool.
- Children exist only to serve "%s"; never include sibling-level functions, never repeat existing names, never include "%s" itself.
- A good child is one the user would click on to see its own next-level breakdown. Prefer 0–3 meaningful children over padding.
- At most 3 children.`, ri.FullPath, ri.FullPath, ri.FullPath, ri.FullPath, ri.FullPath, ri.FullPath, ri.FullPath, ri.FullPath)
```
3. Keys must be of form "<path>.<child_name>".
4. Empty runes map = leaf node.
5. Runs up to cfg.MaxDepth deep, stops at cfg.RuneCap total runes.
6. Events are emitted via buffered channel; sess.Apply updates state in-line.

---

## 6. internal/examples/ — Example Corpus

### 6.1 File Format

Each `.md` file in `examples/{trivial,small,medium,large}/` follows this exact shape:

```markdown
# Requirement: "the exact user requirement as a quoted string"

Free-form decomposition tree in the same text-DSL used in decompose.md.
```

**Parsing rules:**
- Extract requirement via regex: `(?m)^\s*#?\s*Requirement:\s*"([^"]+)"`
- Slug = filename without `.md`
- Tier = parent directory name
- Skip files without a parseable requirement header
- Skip `README.md`

### 6.2 Index API

```go
func LoadFromDir(root string) (*Index, error)
func (idx *Index) Lookup(ref string) LookupResult
func (idx *Index) Manifest() string
```

**Manifest shape (inlined into system prompt):**
```markdown
## trivial (74)
- trivial/add-two-integers
- trivial/hello-world
...

## small (465)
- small/ansi-color-library-with-sgr-code-wrapping
...
```

### 6.4 The `std/` Structured Rune Registry

Separate from the `examples/` corpus, `std/` holds canonical stdlib rune definitions in a YAML-frontmatter format:

```yaml
---
version: 1.0.0
signature: '(pattern: string, s: string) -> result[bool, string]'
dependencies:
  - std.regex.compile
  - std.regex.match
responsibility: regex_matching
---

# std.regex.test

Convenience: compiles a pattern, tests it against the input, returns the boolean result.

## Signature

(pattern: string, s: string) -> result[bool, string]

## Behavior
...
```

**Rules:**
- File path encodes hierarchy: `std/<package>/<unit>/<version>.md`
- `{}` in paths denotes literal braces in the filesystem (e.g., `std/concurrency/once/{do}/`)
- Frontmatter fields: `version`, `signature`, `dependencies` (list of fully-qualified paths), `responsibility`
- Body is free-form markdown following the same DSL as `decompose.md`

For a 3-hour hackathon, create at least one file (e.g., `std/regex/test/1.0.0.md`) to validate the concept.

### 6.5 Corpus Strategy for Hackathon

**Minimum viable:** Create `examples/trivial/` with ~5 files:
- `hello-world.md`
- `add-two-integers.md`
- `reverse-string.md`
- `fibonacci-sequence.md`
- `temperature-converter.md`

Use the shape anchor examples from decompose.md as content. This is enough to boot the system. The full corpus (~1,000 files) improves quality but is not required to make the app work.

---

## 7. internal/effort/ — Complexity Estimator

```go
func Estimate(ctx context.Context, client *openai.Client, requirement string) (Result, error)
```

**System prompt:** `You are a software-complexity estimator. Given a software requirement, rate it 1-5 by calling the rate_effort tool. Reply only via the tool call.`

**Tool schema (rate_effort):**
```json
{
  "level": {"type":"integer","minimum":1,"maximum":5,"description":"1=trivial; 2=small; 3=medium; 4=large; 5=very large"},
  "reason": {"type":"string","description":"One short sentence justifying the level."}
}
```

---

## 8. internal/tui/ — Terminal UI Specification

### 8.1 Design System (styles.go)

```go
var (
    accent     = lipgloss.Color("212")
    accentSoft = lipgloss.Color("99")
    bgMain     = lipgloss.Color("#171717")
    fgBright   = lipgloss.Color("15")
    fgBody     = lipgloss.Color("245")
    fgDim      = lipgloss.Color("241")
    mockSep    = lipgloss.Color("#3e3e3e")
    mockHot    = lipgloss.Color("#f74e82")
    mockFocus  = lipgloss.Color("#9d7cf3")
)
var kanjiPool = []rune("日月火水木金土山川風花雪心愛空東西南北春夏秋冬父母兄弟姉妹雨雲雷電林森竹松梅桜龍虎鳥魚馬犬猫")
func kanjiAt(row, col int) rune
```

### 8.2 Landing Page (landing.go)

- Full-screen alternate buffer, `#171717` background.
- ODEK logo mask: 10-row ASCII art of the word "ODEK" at kanji resolution (each `#` = one kanji slot, 2 cells wide).
- Bouncing animation: Logo drifts inside terminal bounds, bouncing off edges. Tick every 80ms.
- Scrolling kanji field: Rows drift horizontally at alternating directions. Kanji under logo mask get warm VHS gradient (yellow to amber to orange to red to brown). Other kanji muted gray.
- Help bar near bottom center: `n  new  •  q  quit`
- Keybindings: `n` → new feature; `q` / `esc` / `ctrl+c` → quit.

**Additional rendering helpers:**
```go
func renderGradientOnBg(text string, stops []colorful.Color, bg string, totalWidth int) string
    // Paints a per-character horizontal gradient over text on a solid background.
    // Used for help-line chars that fall under the logo mask.
```

### 8.3 Chat Page (create_feature.go + chat.go)

**Layout (single-pane mode):**
```
[ODEK logo]  [status]  [up-arrow if scrollable]
[kanji line 1]
[kanji line 2]
[viewport: chat history]
[input textarea]
[help bar: enter send • alt+enter new line • up/down scroll]
```

**Chat component features:**
- textarea for input (alt+enter for newline).
- viewport for history with soft wrap.
- User messages: left thick border in accent pink.
- Assistant messages: left thick border in fgDim; label "clank".
- Assistant headlines: small inverted pill (e.g., "Effort: 2/5").
- Markdown rendering with Chroma syntax highlighting (dracula theme, terminal16m formatter).
- Inline code: backtick spans get bright-on-dark background.
- Code blocks: padded background `#1e1e1e`, 3-space left margin.
- Thinking box: When LLM streams reasoning_content, a bordered panel labeled "thinking" appears under latest user message, showing live deltas (max 8 wrapped lines, tail-truncated).

**SendHandler flow:**
1. User presses Enter → message appended, input cleared, status = chatSending.
2. Chat history sent to chat classifier LLM (see 8.5).
3. If classifier calls decompose tool → run runDecompose().
4. If classifier replies in text → show as normal assistant reply.
5. decomposeDoneMsg carries reply + optional expansion event channel.

### 8.4 Decomposition Browser (feature_decomp.go)

**Layout:**
```
[ODEK logo]
[kanji lines]
[summary text]
[rune count + status]
[Miller columns + detail pane]
[input line if active]
[help bar]
```

**Miller-column navigation:**
- Column 0: sectioned by package (std header + its runes, then project header + its runes).
- Columns 1+: flat list of children for selected rune from previous column.
- Navigation: arrows/hjkl move and drill, blink cursor, steady-after-move.
- Column width auto-sizing based on content (min 14, nominal 22).
- Horizontal scroll for deep trees (max 5 visible columns).
- Detail pane (rightmost): shows selected rune's description, signature, tests, assumptions, children.

**Empty state:**
```
◆ root
─────────────────────

send a message to start
```

**Loading state:** When state.decomposing or sess.Snapshot().Expanding, kanji scroll faster and placeholder text changes to "decomposing...".

**Missing function signatures (must be implemented for compilation):**
```go
func renderDecompTop(width, height, kanjiOffset int, snap decomposer.Snapshot) string
func wrapDecompText(text string, width int) []string
func statusTag(st decomposer.RuneStatus) string
func statusGlyph(selected, active, blinkOn bool) (string, lipgloss.Style)
func renderRuneInfo(name string, r decomposer.Rune, status decomposer.RuneStatus, children []string, maxW int) string
func renderDecompHelp(width int, inputActive, showTabSwitch bool) string
func renderFeaturePin() string
func renderDetailPane(selPath []string, snap decomposer.Snapshot, w int, titled func(string, string, int, bool) string) string
func displayName(snap decomposer.Snapshot, path string) string
```

**Session navigation methods on `featureDecompModel`:**
```go
func (m *featureDecompModel) parentOfCol(colIdx int) string
    // Returns "root" for column 0, otherwise the selection from the previous column.

func (m *featureDecompModel) normalizeSelection()
    // Truncates stale selPath entries, clamps focusedCol/colScroll into range,
    // and ensures the selection always points to a valid child.

func (m *featureDecompModel) ensureColVisible()
    // Adjusts colScroll so focusedCol is within [colScroll, colScroll+decompMaxVisCols).

func (m *featureDecompModel) steadyCursor()
    // Sets blinkOn=true and a 500ms grace period so rapid navigation stays visible.

func (m *featureDecompModel) moveSelection(delta int)
    // Shifts selection in the focused column by delta, truncates deeper selections.
```

**`buildColumns` exact signature:**
```go
func buildColumns(innerW, innerH int, selPath []string, focusedCol, colScroll int, active, blinkOn, decomposing bool, snap decomposer.Snapshot) string
```

**Internal closures inside `buildColumns`:**
- `renderRuneRow(path, selected string, focused bool, contentW int) string`
  - Produces: status glyph (`• `) + rune name + optional `›` drill-in hint.
  - Truncates names with `…` if they exceed content width.
- `applyScroll(parts []string, selRow, contentW int) string`
  - Windows `parts` around `selRow` to fit in `innerH` rows.
  - Adds `↑`/`↓` chrome rows when content is clipped.
- `titled(name, tag string, w int, focused bool) string`
  - Renders a column header: `◆ <name>` left-aligned, `# <tag>` right-aligned, horizontal rule below.

### 8.5 Split Pane (split_feature.go)

- Activates when terminal width >= 150 columns.
- Left = chat pane (~1/3 width). Right = decomposition browser (~2/3).
- tab / shift+tab switches focus.
- Unfocused pane dimmed (logo + help bar rows only).
- ctrl+]/[ resizes left pane by 4 columns.
- enter in right pane (on a rune) copies `path: ` into chat input and focuses left pane.

### 8.6 Chat Classifier System Prompt

This is the prompt given to the chat LLM on every user turn:

```markdown
You are Odek, a software library design collaborator. The user is iterating on a library spec, and on every turn you either discuss it or update it.

Classify the user's latest message by intent:

1. **Spec change** — the user is revising what the library should do: adding/removing/renaming capabilities, tightening or relaxing scope, changing the library's purpose. Call the `decompose` tool. Do not reply with text before calling it; the tool's output becomes the reply.

2. **Discussion** — the user is asking a question, exploring tradeoffs, clarifying how something works, or giving feedback that doesn't change what the library does. Reply in plain text. Be concise and practical. Do not paste rune trees yourself.

Examples:
- "make it a scientific calculator" → spec change, call the tool.
- "also support matrices" → spec change, call the tool.
- "drop the divide function" → spec change, call the tool.
- "how does logarithm work?" → discussion, reply in text.
- "what should divide-by-zero return?" → discussion, reply in text (until the user picks a behavior).
- "what's the tradeoff between X and Y?" → discussion, reply in text.

When in doubt, prefer discussion: it's cheap to follow up with a spec change, but a spurious rewrite costs the user their prior structure.
```

**Chat decompose tool schema:**
```json
{
  "name": "decompose",
  "description": "Apply a scope or requirement change to the feature spec. Call this whenever the user's latest message revises what the library should do...",
  "parameters": {
    "type": "object",
    "properties": {
      "levels": {"type":"integer","minimum":1,"maximum":10,"description":"Recursion depth. 1 = top-level only."},
      "effort": {"type":"integer","minimum":1,"maximum":5,"description":"Effort level 1-5. Default 2."}
    }
  }
}
```

### 8.7 Page Transition (transition.go)

- Duration: 320ms at 60fps.
- Easing: easeOutCubic.
- Outgoing view shrinks to small rectangle anchored at feature pin row while incoming view slides in from right.
- Pin label: " feature " rendered in inverted accent style.
- Any keypress during transition immediately completes it.

---

## 9. Alternate Entry Points (cmd/)

The `cmd/` directory contains standalone binaries separate from the main TUI:

### 9.1 `cmd/auto_recurse.go`
- `package main` — build with `go run ./cmd/auto_recurse.go ./cmd/print.go`
- Non-TUI CLI flow:
  1. Prompts for a requirement via stdin.
  2. Calls `effort.Estimate()` to get complexity level.
  3. Calls `decomposer.NewSession()` for initial decomposition.
  4. Prints the initial tree via `printInitialDecomposition()`.
  5. Prompts `[y/N]` to proceed with recursion.
  6. Runs `ExpandStreaming` and prints events + final tree.

### 9.2 `cmd/tui/main.go`
- `package main` — a **hardcoded prototype** of the decomposition browser.
- Contains static `tomlItems` and `runeInfos` (9 TOML runes) baked into the binary.
- Demonstrates the 2-column layout + detail pane but does NOT talk to an LLM.
- Useful as a UI reference, but the real app uses `internal/tui/`.

### 9.3 `cmd/print.go`
- Shared print helpers for `cmd/auto_recurse.go`:
  - `printBanner()`, `printInitialDecomposition()`, `printCompleteTree()`, `printRunesIndented()`, `printExpansionEvent()`, `wrapText()`, `plural()`

### 9.4 `cmd/recurse_detailed.go`
- Entirely commented-out old prototype. Can be ignored.

## 10. main.go — Entry Point

```go
func main() {
    // 1. API_BASE_URL env var or default http://localhost:8080
    // 2. Parse -p=<prompt> or -p <prompt> for direct chat
    // 3. Parse -d=<requirement> for direct structured decompose
    // 4. Parse -j or --json for structured output flag
    // 5. If -d: call decompose.DecomposeStructured, print JSON + token usage
    // 6. If -p: call client.Chat with user prompt, print response + token usage
    // 7. Otherwise: launch TUI
    dec, _ := decomposer.NewDecomposer(client, "examples", "/tmp/odek-example-log.jsonl")
    tui.Run(ctx, client, dec)
}
```

---

## 11. Example Corpus File (Template)

Create this file at `examples/trivial/hello-world.md`:

```markdown
# Requirement: "hello world"

A library that returns the canonical greeting. Printing is the caller's responsibility.

std: (all units exist)

greeter
  greeter.greet
    fn () -> string
    + returns the string "Hello, world!"
    ? greeting is hardcoded; no parameter
    # greeting
```

---

## 12. 3-Hour Hackathon Implementation Prompts

Use these prompts in order with an AI coding agent to rebuild the project from scratch.

### Prompt 1: Foundation (30 min)
```
Create a Go module shotgun.dev/odek with go.mod using Go 1.22.
Create package openai/ with:
- Client struct with NewClient(baseURL, apiKey...), ListModels, HealthCheck
- All chat types: ChatMessage, Tool, ToolCall, FunctionDefinition, ChatCompletionRequest/Response, Choice, Delta, Usage
- Chat() method with SSE streaming support (stream=true, Accept text/event-stream, parse data: lines, accumulate deltas and tool call fragments)
- Ask(), AskMessages(), AskTool(), and AskToolLoop() exactly as described in the PRD
- WithThinkingCallback context key
Implement decompose/ package with Decompose(), DecomposeStructured(), ParseDecomposition(), Validate, FormatJSON.
Write a small test in openai/ that mocks the server.
```

### Prompt 2: Core Engine (45 min)
```
Create internal/decomposer/ with:
- decompose.md as an embedded file (//go:embed) containing the EXACT system prompt from the PRD section 5.1
- Types: Rune, PackageNode, DecompositionResponse, ClarificationRequest, ClarificationNeeded, AutoDecomposition, RuneExpansionInfo, RuneStatus constants, Snapshot, Session
- Config and ConfigForEffort
- NewDecomposer that loads examples and sets up tool schemas
- Decompose() method: multi-turn tool loop with read_example and decompose tools (max 6 iterations)
- Handle plain-text fallback: retry once with forced tool call, then return ClarificationRequest
- NewSession with optional parallel initial attempts + MergeAttempts
- expand.go: ExpandStreaming with BFS level-by-level expansion, concurrent rune expansion, event emission
- events.go: all ExpansionEvent variants
- normalize.go: NormalizeFunctionSig and normalizePackageSignatures
Make sure the decompose tool schema matches the PRD exactly.
```

### Prompt 3: Examples & Utilities (15 min)
```
Create internal/examples/ with LoadFromDir, Lookup, Manifest.
Create internal/effort/ with Estimate via forced rate_effort tool call.
Create internal/toollog/ with JSONL append-only logger.
Create the examples/ directory with at least 5 trivial .md files following the corpus format from the PRD.
```

### Prompt 4: TUI Styles & Landing (20 min)
```
Create internal/tui/styles.go with the exact color palette, kanjiPool, kanjiAt, help bindings.
Create internal/tui/landing.go with:
- Bouncing ODEK logo mask (10 rows, kanji-resolution ASCII art)
- Scrolling kanji field with warm VHS gradient inside logo mask
- Help bar at bottom center
- Key handling: n -> new feature, q/esc/ctrl+c -> quit
- 80ms tick for animation
```

### Prompt 5: Chat Component (25 min)
```
Create internal/tui/chat.go with:
- chatModel using bubbles/v2 textarea and viewport
- Render user messages with pink left border
- Render assistant messages with gray left border, headline pills
- Markdown parsing with Chroma syntax highlighting (dracula theme, terminal16m, #1e1e1e code background)
- Inline code backtick styling
- Code blocks: padded background #1e1e1e, 3-space left margin
- Thinking box with "thinking" label, max 8 lines, live update via SetPendingThinking
- SendHandler callback interface
- History() filtering out system notes
```

### Prompt 6: Create Feature & Orchestration (25 min)
```
Create internal/tui/create_feature.go with:
- decomposeState heap-allocated holder with mutex (session, decomposing flag, thinking buffer)
- createFeatureModel: chat + kanji ticker
- makeFeatureSendHandler: builds chat messages, calls chat classifier LLM
- chatHandler: posts to LLM with chatDecomposeTool; if tool called, runs runDecompose
- runDecompose: builds SessionContext from discussion+prior, calls dec.NewSession, starts ExpandStreaming if levels>1, returns decomposeDoneMsg
- buildDiscussion and extractRequirement helpers
- pumpExpansionCmd to feed events into Bubble Tea loop
```

### Prompt 7: Decomposition Browser (30 min)
```
Create internal/tui/feature_decomp.go with:
- Miller-column layout: sectioned column 0 (std + project), flat columns 1+
- Navigation: arrows/hjkl move and drill, blink cursor, steady-after-move
- Column width auto-sizing based on content (min 14, nominal 22)
- Horizontal scroll for deep trees (max 5 visible columns)
- Detail pane showing description, signature, tests, assumptions, children
- Empty state and loading state
- Input mode on space (textinput for future use)
- decomp_render.go: renderDecompositionSummary
```

### Prompt 8: Split Pane & Transitions (15 min)
```
Create internal/tui/split_feature.go with:
- Split pane at >=150 columns, leftW starts at 1/3
- Tab switches focus; unfocused pane dimmed
- Ctrl+]/[ resizes; Enter on rune copies path to chat
- clipToBox and dimUnfocused helpers
Create internal/tui/transition.go with:
- 320ms slide-collapse transition
- easeOutCubic easing
- composeSlideCollapse with ANSI-aware overlay
```

### Prompt 9: main.go & Wiring (15 min)
```
Write main.go that:
- Reads API_BASE_URL env var (default http://localhost:8080)
- Parses -p, -d, -j flags
- Direct prompt mode: chat and print response
- Direct decompose mode: DecomposeStructured and print JSON
- Default mode: NewDecomposer + tui.Run
Build with go build and verify no compile errors.
```

### Prompt 10: Polish & Test (10 min)
```
Run go test ./... and fix any failures.
Run the TUI with a local LLM server (e.g., lmstudio on port 8080).
Verify:
- Landing screen animates
- n key enters chat
- Typing a requirement triggers decomposition
- Ctrl+Enter shows the rune tree
- Tab switches between chat and runes
- Arrow keys navigate the Miller columns
```

**Existing test files to preserve/regenerate:**
- `internal/decomposer/normalize_test.go` — 17 cases for `NormalizeFunctionSig` (strips `fn`, `@`, whitespace)
- `internal/examples/examples_test.go` — verifies corpus loads all 4 tiers with >=100 entries
- `internal/tui/split_feature_test.go` — verifies split pane renders at exact dimensions
- `internal/tui/transition_test.go` — verifies transition animation frame dimensions

**CI note:** `.github/workflows/test.yml` runs `go build ./...` and `go test ./e2e/ -v`. The `e2e/` directory does not currently exist; either create it or remove that step from CI.

---

## 13. Environment Setup Checklist

- [ ] Go 1.22+ installed
- [ ] Local OpenAI-compatible server running (e.g., LM Studio, Ollama, llama.cpp server) on localhost:8080
- [ ] API_BASE_URL env var set if using a remote server
- [ ] Terminal supports true color and alternate screen (any modern terminal)

---

## 14. Key Design Decisions (Do Not Change)

1. **Library output only.** Never generate main(), CLI args, or process lifecycle. The product is always a reusable library.
2. **Stdlib-first.** Generic capabilities go in std.* before feature-specific code.
3. **Tool calls only.** The decompose agent must never output plain text trees. It uses read_example then decompose tools.
4. **No duplicates.** Every rune name exists in exactly one package. References use -> std.path.unit.
5. **Session owns state.** The Session struct is heap-allocated and shared across chat + decomposition panes via pointer.
6. **Bubble Tea v2.** The entire TUI uses charm.land/bubbletea/v2 (not v1) and charm.land/lipgloss/v2.
7. **Effort drives config.** Complexity estimation maps 1-5 to parallel attempts, depth, and rune cap.

---

## 15. Future Enhancements (Post-Hackathon)

- Export decomposition to JSON / Markdown / code stubs
- Persistent project history (save/load sessions)
- Edit runes inline in the detail pane
- Generate actual implementation code from rune trees
- Search/filter example corpus
- Multi-model support (switch between local and remote APIs)
