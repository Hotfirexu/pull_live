package http

import (
	"net/http"
	"strings"
)

type LineReader interface {
	ReadLine() (line []byte, isPrefix bool, err error)
}

func ReadHttpHeader(r LineReader) (firstLine string, headers http.Header, err error) {
	headers = make(http.Header)

	readLineFn := func() (string, error) {
		var line string
		var bLine []byte
		var isPrefix bool
		for {
			bLine, isPrefix, err = r.ReadLine()
			if err != nil {
				return "", err
			}

			line += string(bLine)
			if !isPrefix {
				break
			}
		}
		return line, nil
	}

	firstLine, err = readLineFn()
	if err != nil {
		return
	}

	if len(firstLine) == 0 {
		err = ErrHttpHeader
		return
	}

	var lastKey string
	for {
		var l string
		l, err = readLineFn()
		if err != nil {
			return
		}

		if len(l) == 0 {
			break
		}

		pos := strings.Index(l, ":")
		if pos == -1 {
			if lastKey != "" {
				vs := headers.Values(lastKey)
				vs[len(vs)-1] = vs[len(vs)-1] + l
			}
			continue
		}
		lastKey = strings.Trim(l[0:pos], " ")
		headers.Add(strings.Trim(l[0:pos], " "), strings.Trim(l[pos+1:], " "))
	}
	return
}

func ParseHttpStatusLine(line string) (version string, statusCode string, reason string, err error) {
	return parseFirstLine(line)
}

func parseFirstLine(line string) (item1, item2, item3 string, err error) {
	f := strings.Index(line, " ")
	if f == -1 {
		err = ErrFirstLine
		return
	}

	s := strings.Index(line[f+1:], " ")
	if s == -1 {
		return line[0:f], line[f+1:], "", nil
	}
	if f+1+s+1 == len(line) {
		return line[0:f], line[f+1 : f+1+s], "", nil
	}
	return line[0:f], line[f+1 : f+1+s], line[f+1+s+1:], nil
}