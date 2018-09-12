package constants

import "errors"

// Errors that could be occurred
var (
	ErrNoOutputDocForMethod = errors.New("failed to get output documentation for given method")
	ErrNoInputDocForMethod  = errors.New("failed to get input documentation for given method")
	ErrProtoDoc             = errors.New("for using RPC you must have autodoc implemented on given server type")
	ErrNoServerType         = errors.New("no server type chosen for getting autodoc")
	ErrGzip                 = errors.New("failed to unzip gzip binary")
	ErrEmptyRequest         = errors.New("empty request body")
)
