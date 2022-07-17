package httpflv

const (
	TagHeaderSize        int = 11
	PrevTagSizeFieldSize int = 4
	flvHeaderSize        int = 13
)

const (
	TagTypeMetadata uint8 = 18
	TagTypeVideo    uint8 = 9
	TagTypeAudio    uint8 = 8
)

var readBufSize = 256

const (
	SoundFormatAac uint8 = 10
)

const (
	AvcPacketTypeSeqHeader uint8 = 0
	AvcPacketTypeNalu      uint8 = 1

	HevcPacketTypeSeqHeader uint8 = 0
	HevcPacketTypeNalu      uint8 = 1

	AacPacketTypeSeqHeader uint8 = 0
	AacPacketTypeRaw       uint8 = 1
)

const (
	AvcKeyFrame   = frameTypeKey<<4 | codecIdAvc
	AvcInterFrame = frameTypeInter<<4 | codecIdAvc

	HevcKeyFrame   = frameTypeKey<<4 | codecIdHevc
	HevcInterFrame = frameTypeInter<<4 | codecIdHevc
)

const (
	codecIdAvc  uint8 = 7
	codecIdHevc uint8 = 12
)

const (
	frameTypeKey   uint8 = 1
	frameTypeInter uint8 = 2
)
