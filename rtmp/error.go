package rtmp

import (
	"errors"
	"fmt"
)

var (
	ErrAmfInvalidType = errors.New("rtmp: invalid amf0 type")
	ErrAmfTooShort    = errors.New("rtmp: too short to unmarshal amf0 data")
	ErrAmfNotExist    = errors.New("rtmp: not exist")

	ErrRtmpShortBuffer   = errors.New("rtmp: buffer too short")
	ErrRtmpUnexpectedMsg = errors.New("rtmp: unexpected msg")
)

func NewErrAmfInvalidType(b byte) error {
	return fmt.Errorf("%w. b=%d", ErrAmfInvalidType, b)
}
