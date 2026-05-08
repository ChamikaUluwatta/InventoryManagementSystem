package testutil

import (
	"bytes"
	"encoding/json"
	"testing"
)

func MarshalBody(t *testing.T, v any) *bytes.Buffer {
	t.Helper()
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(v); err != nil {
		t.Fatalf("failed to marshal body: %v", err)
	}
	return &buf
}
