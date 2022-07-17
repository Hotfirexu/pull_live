package http

import "errors"

var (
	ErrHttpHeader   = errors.New("nazahttp: read http header failed")
	ErrFirstLine    = errors.New("nazahttp: parse first line failed")
	ErrParamMissing = errors.New("nazahttp: param missing")
)

const (
	HeaderFieldContentLength = "Content-Length"
	HeaderFieldContentType   = "application/json"
)
