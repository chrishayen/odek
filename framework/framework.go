package framework

import "embed"

//go:embed go/dispatch.go
var GoDispatch string

// FS provides access to all framework files.
//
//go:embed go/*
var FS embed.FS
