package utils

import (
	"bytes"
	"encoding/json"
	"time"
)

// Parse time format from RFC3339 to YYYY-MM-DD
func ParseTimeFormat(t *time.Time) *string {
	if t == nil {
		return nil
	}
	res := t.Format("2006-01-02")
	return &res
}

// Convert bytes to buffer helper
func AnyToBytesBuffer(i interface{}) (*bytes.Buffer, error) {
	buf := new(bytes.Buffer)
	err := json.NewEncoder(buf).Encode(i)
	if err != nil {
		return buf, err
	}
	return buf, nil
}
