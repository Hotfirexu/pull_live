package base

import "errors"

var (
	ErrAddrEmpty               = errors.New("http server addr empty")
	ErrMultiRegisterForPattern = errors.New("http server multiple registrations for pattern")

	ErrSessionNotStarted = errors.New("session has not been started yet")

	ErrInvalidUrl = errors.New("invalid url")
)

var (
	ErrShortBuffer  = errors.New("buffer too short")
	ErrFileNotExist = errors.New("file not exist")
)

var ErrAvc = errors.New("avc: fxxk")
