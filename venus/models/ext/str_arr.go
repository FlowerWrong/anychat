package ext

import (
	"bytes"
	"strings"
)

// StrArray ...
type StrArray []string

// FromDB {barton.com}
func (s *StrArray) FromDB(bts []byte) error {
	if len(bts) == 0 {
		return nil
	}

	str := string(bts)
	if strings.HasPrefix(str, "{") && strings.HasSuffix(str, "}") {
		str = str[1 : len(str)-1]
	}

	var out []string
	for _, el := range strings.Split(str, ",") {
		out = append(out, strings.TrimSpace(el))
	}

	*s = StrArray(out)
	return nil
}

// ToDB ...
func (s *StrArray) ToDB() ([]byte, error) {
	return serializeStrArray(*s, "{", "}"), nil
}

// serializeStrArray ...
func serializeStrArray(s []string, prefix string, suffix string) []byte {
	var buffer bytes.Buffer

	buffer.WriteString(prefix)

	for idx, val := range s {
		if idx > 0 {
			buffer.WriteString(",")
		}
		buffer.WriteString(val)
	}

	buffer.WriteString(suffix)

	return buffer.Bytes()
}
