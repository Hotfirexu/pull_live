package httpflv

import (
	"io"
	"pull_live/pkg/bele"
)

type TagHeader struct {
	Type      uint8  // type
	DataSize  uint32 // body大小  不包含header和prev tag size字段
	TimeStamp uint32 // 绝对时间戳  单位是毫秒
	StreamId  uint32 // always 0
}

type Tag struct {
	Header TagHeader
	Raw    []byte // 结构为:(11字节的tag header) + (body) + (4字节的 prev tag size)
}

func (tag *Tag) Payload() []byte {
	return tag.Raw[TagHeaderSize : len(tag.Raw)-PrevTagSizeFieldSize]
}

func (tag *Tag) IsAacSeqHeader() bool {
	return tag.Header.Type == TagTypeAudio && tag.Raw[TagHeaderSize]>>4 == SoundFormatAac && tag.Raw[TagHeaderSize+1] == AacPacketTypeSeqHeader
}

func ReadTag(rd io.Reader) (tag Tag, err error) {
	rawHeader := make([]byte, TagHeaderSize)
	if _, err = io.ReadAtLeast(rd, rawHeader, TagHeaderSize); err != nil {
		return
	}

	header := parseTagHeader(rawHeader)

	needed := int(header.DataSize) + PrevTagSizeFieldSize
	tag.Header = header
	tag.Raw = make([]byte, TagHeaderSize+needed)
	copy(tag.Raw, rawHeader)

	if _, err = io.ReadAtLeast(rd, tag.Raw[TagHeaderSize:], needed); err != nil {
		return
	}

	return
}

func (tag *Tag) IsVideoKeySeqHeader() bool {
	return tag.IsAvcKeySeqHeader() || tag.IsHevcKeySeqHeader()
}

func (tag *Tag) IsAvcKeySeqHeader() bool {
	return tag.Header.Type == TagTypeVideo && tag.Raw[TagHeaderSize] == AvcKeyFrame && tag.Raw[TagHeaderSize+1] == AvcPacketTypeSeqHeader
}

func (tag *Tag) IsHevcKeySeqHeader() bool {
	return tag.Header.Type == TagTypeVideo && tag.Raw[TagHeaderSize] == HevcKeyFrame && tag.Raw[TagHeaderSize+1] == HevcPacketTypeSeqHeader
}

func parseTagHeader(rawHeader []byte) TagHeader {
	var h TagHeader
	h.Type = rawHeader[0]
	h.DataSize = bele.BeUint24(rawHeader[1:])
	h.TimeStamp = (uint32(rawHeader[7]) << 24) + bele.BeUint24(rawHeader[4:])
	return h
}
