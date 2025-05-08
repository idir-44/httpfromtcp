package headers

import (
	"bytes"
	"fmt"
	"strings"
	"unicode"
)

type Headers map[string]string

func NewHeaders() Headers {
	return Headers{}
}

const crlf = "\r\n"

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	idx := bytes.Index(data, []byte(crlf))
	if idx == -1 {
		return
	}

	rawHeader := string(data[:idx])

	if rawHeader == "" {
		return 2, true, nil
	}

	headerFields, err := sliceString(rawHeader, ":")
	if err != nil {
		return 0, false, fmt.Errorf("invalid format for field-line")
	}

	if strings.HasSuffix(headerFields[0], " ") {
		return 0, false, fmt.Errorf("field-name should not end with a space")
	}

	fieldName := strings.ToLower(strings.TrimSpace(headerFields[0]))
	if !isHeaderKeyValid(fieldName) {
		return 0, false, fmt.Errorf("invalid field-name")
	}

	fieldValue := strings.TrimSpace(headerFields[1])

	h.Set(fieldName, fieldValue)

	return idx + 2, false, nil
}

func (h Headers) Set(key, value string) {
	if _, ok := h[key]; ok {
		h[key] = strings.Join([]string{h[key], value}, ", ")
		return
	}
	h[key] = value
}

func (h Headers) Get(key string) (value string, ok bool) {
	value, ok = h[strings.ToLower(key)]

	return
}

func sliceString(s, sep string) ([]string, error) {
	idx := strings.Index(s, sep)
	if idx == -1 {
		return []string{}, fmt.Errorf("sep not found")
	}

	return []string{s[:idx], s[idx+1:]}, nil
}

func isHeaderKeyValid(key string) bool {
	symbols := []rune{'!', '#', '$', '%', '&', '\'', '*', '+', '-', '.', '^', '_', '`', '|', '~'}

	for _, s := range key {
		if !unicode.IsUpper(s) && !unicode.IsLower(s) && !unicode.IsDigit(s) && !contains(symbols, s) {
			return false
		}
	}
	return true
}

func contains(r []rune, sep rune) bool {
	for _, s := range r {
		if s == sep {
			return true
		}
	}

	return false
}
