package tui

import "shotgun.dev/odek/decompose"

func mockDecomposition() (*decompose.Decomposition, []responsibilityGroup) {
	// networking leaves
	listen := &decompose.Rune{
		Path:        "hello_world/server/network/listen",
		Version:     "1",
		Signature:   "(port: u16) -> result[socket, string]",
		Description: "Binds a TCP socket on the given port and returns a listener ready to accept connections.",
	}
	accept := &decompose.Rune{
		Path:        "hello_world/server/network/accept",
		Version:     "1",
		Signature:   "(l: socket) -> result[conn, string]",
		Description: "Blocks until a client connects, then yields a per-request connection handle.",
	}
	network := &decompose.Rune{
		Path:        "hello_world/server/network",
		Version:     "1",
		Signature:   "(port: u16) -> result[conn_stream, string]",
		Description: "Low level socket I/O: binds the port and streams incoming connections.",
		Children:    []*decompose.Rune{listen, accept},
	}

	// http parser leaves
	parseRequestLine := &decompose.Rune{
		Path:        "hello_world/server/http/parser/parse_request_line",
		Version:     "1",
		Signature:   "(line: bytes) -> result[request_line, string]",
		Description: "Parses the first line of an HTTP request into method, path, and version.",
	}
	parseHeaders := &decompose.Rune{
		Path:        "hello_world/server/http/parser/parse_headers",
		Version:     "1",
		Signature:   "(raw: bytes) -> result[map[string, string], string]",
		Description: "Reads CRLF-delimited header lines into a map until the blank line.",
	}
	parser := &decompose.Rune{
		Path:        "hello_world/server/http/parser",
		Version:     "1",
		Signature:   "(c: conn) -> result[request, string]",
		Description: "Parses an HTTP/1.1 request off the connection.",
		Children:    []*decompose.Rune{parseRequestLine, parseHeaders},
	}

	// http writer leaves
	writeStatus := &decompose.Rune{
		Path:        "hello_world/server/http/writer/write_status",
		Version:     "1",
		Signature:   "(c: conn, code: u16) -> result[void, string]",
		Description: "Writes the status line, e.g. 'HTTP/1.1 200 OK'.",
	}
	writeBody := &decompose.Rune{
		Path:        "hello_world/server/http/writer/write_body",
		Version:     "1",
		Signature:   "(c: conn, body: bytes) -> result[void, string]",
		Description: "Writes the content-length header and the response body bytes.",
	}
	writer := &decompose.Rune{
		Path:        "hello_world/server/http/writer",
		Version:     "1",
		Signature:   "(c: conn, resp: response) -> result[void, string]",
		Description: "Serializes an HTTP response back onto the connection.",
		Children:    []*decompose.Rune{writeStatus, writeBody},
	}

	httpComp := &decompose.Rune{
		Path:        "hello_world/server/http",
		Version:     "1",
		Signature:   "(c: conn) -> result[void, string]",
		Description: "HTTP/1.1 codec: parses inbound requests, serializes outbound responses.",
		Children:    []*decompose.Rune{parser, writer},
	}

	handleRequest := &decompose.Rune{
		Path:        "hello_world/server/handle_request",
		Version:     "1",
		Signature:   "(c: conn) -> result[void, string]",
		Description: "Glues the parser, greeter, and writer together for a single request.",
	}

	server := &decompose.Rune{
		Path:        "hello_world/server",
		Version:     "1",
		Signature:   "(port: u16) -> result[void, string]",
		Description: "Wires networking, HTTP codec, and request handling into a running server.",
		Children:    []*decompose.Rune{network, httpComp, handleRequest},
	}

	// greeter sub-tree
	loadTemplate := &decompose.Rune{
		Path:        "hello_world/greeter/template/load",
		Version:     "1",
		Signature:   "(name: string) -> result[template, string]",
		Description: "Loads a named greeting template from the embedded assets.",
	}
	renderTemplate := &decompose.Rune{
		Path:        "hello_world/greeter/template/render",
		Version:     "1",
		Signature:   "(t: template, vars: map[string, string]) -> string",
		Description: "Substitutes the template's placeholders with the given variables.",
	}
	template := &decompose.Rune{
		Path:        "hello_world/greeter/template",
		Version:     "1",
		Signature:   "(name: string, vars: map[string, string]) -> result[string, string]",
		Description: "Small templating layer used by the greeter to build the response body.",
		Children:    []*decompose.Rune{loadTemplate, renderTemplate},
	}
	greet := &decompose.Rune{
		Path:        "hello_world/greeter/greet",
		Version:     "1",
		Signature:   "(name: string) -> string",
		Description: "Returns a friendly hello for the caller.",
	}
	greeter := &decompose.Rune{
		Path:        "hello_world/greeter",
		Version:     "1",
		Signature:   "(name: string) -> string",
		Description: "Produces the response body that the server hands back to each client.",
		Children:    []*decompose.Rune{template, greet},
	}

	// logging sub-tree
	formatJSON := &decompose.Rune{
		Path:        "hello_world/logging/formatter/format_json",
		Version:     "1",
		Signature:   "(fields: map[string, string]) -> string",
		Description: "Renders a log record as a single JSON object.",
	}
	formatLine := &decompose.Rune{
		Path:        "hello_world/logging/formatter/format_line",
		Version:     "1",
		Signature:   "(fields: map[string, string]) -> string",
		Description: "Renders a log record as a key=value logfmt line.",
	}
	formatter := &decompose.Rune{
		Path:        "hello_world/logging/formatter",
		Version:     "1",
		Signature:   "(fields: map[string, string]) -> string",
		Description: "Chooses a format strategy and renders structured log records.",
		Children:    []*decompose.Rune{formatJSON, formatLine},
	}

	sinkStderr := &decompose.Rune{
		Path:        "hello_world/logging/sink/stderr",
		Version:     "1",
		Signature:   "(line: string) -> void",
		Description: "Writes a formatted log line to stderr with a trailing newline.",
	}
	sinkFile := &decompose.Rune{
		Path:        "hello_world/logging/sink/file",
		Version:     "1",
		Signature:   "(path: string, line: string) -> result[void, string]",
		Description: "Appends a formatted log line to the file at the given path.",
	}
	sink := &decompose.Rune{
		Path:        "hello_world/logging/sink",
		Version:     "1",
		Signature:   "(line: string) -> result[void, string]",
		Description: "Destinations that consume formatted log lines.",
		Children:    []*decompose.Rune{sinkStderr, sinkFile},
	}

	logging := &decompose.Rune{
		Path:        "hello_world/logging",
		Version:     "1",
		Signature:   "(method: string, path: string) -> void",
		Description: "Observability helpers shared by the request path.",
		Children:    []*decompose.Rune{formatter, sink},
	}

	root := &decompose.Rune{
		Path:        "hello_world",
		Version:     "1",
		Signature:   "(port: u16) -> result[void, string]",
		Description: "Root library entry point that starts the hello world server.",
		Children:    []*decompose.Rune{server, greeter, logging},
	}

	d := &decompose.Decomposition{
		FeatureName: "hello world http server",
		Description: "A tiny HTTP server that greets each visitor.",
		RuneTree:    root,
	}

	groups := []responsibilityGroup{
		{
			name:        "networking",
			description: "Accepting connections and routing incoming requests.",
			runePaths: []string{
				server.Path,
				network.Path,
				listen.Path,
				accept.Path,
				handleRequest.Path,
			},
		},
		{
			name:        "http_protocol",
			description: "Parsing inbound requests and serializing outbound responses.",
			runePaths: []string{
				httpComp.Path,
				parser.Path,
				parseRequestLine.Path,
				parseHeaders.Path,
				writer.Path,
				writeStatus.Path,
				writeBody.Path,
			},
		},
		{
			name:        "greeting",
			description: "Producing the response body for each request.",
			runePaths: []string{
				greeter.Path,
				template.Path,
				loadTemplate.Path,
				renderTemplate.Path,
				greet.Path,
			},
		},
		{
			name:        "observability",
			description: "Logging requests for debugging and post-hoc inspection.",
			runePaths: []string{
				logging.Path,
				formatter.Path,
				formatJSON.Path,
				formatLine.Path,
				sink.Path,
				sinkStderr.Path,
				sinkFile.Path,
			},
		},
	}

	return d, groups
}
