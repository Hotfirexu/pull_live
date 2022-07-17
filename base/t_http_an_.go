package base

type (
	SessionType int
)

type StatSession struct {
	SessionId string `json:"session_id"`
	Protocol  string `json:"protocol"`
	BaseType  string `json:"base_type"`

	StartTime string `json:"start_time"`

	RemoteAddr string `json:"remote_addr"`

	ReadBytesSum  uint64 `json:"read_bytes_sum"`
	WroteBytesSum uint64 `json:"wrote_bytes_sum"`
	Bitrate       int    `json:"bitrate"`
	ReadBitrate   int    `json:"read_bitrate"`
	WriteBitrate  int    `json:"write_bitrate"`

	typ SessionType
}
