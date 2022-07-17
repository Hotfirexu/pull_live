package rtmp

func ParseMetaData(b []byte) (ObjectPairArray, error) {
	pos := 0
	v, l, err := Amf0.ReadString(b[pos:])
	if err != nil {
		return nil, err
	}

	pos += l
	if v == "@setDataFrame" {
		_, l, err = Amf0.ReadString(b[pos:])
		if err != nil {
			return nil, err
		}
		pos += l
	}

	opa, _, err := Amf0.ReadObjectOrArray(b[pos:])
	return opa, err
}
