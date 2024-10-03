package utils

import (
	"encoding/json"
	"strconv"

	"github.com/google/uuid"
)

// ParseUnixTime parses a string into an int64 Unix epoch time in seconds.
func ParseUnixTime(timeStr string) (int64, error) {
	return strconv.ParseInt(timeStr, 10, 64)
}

// SerializeToJSON serializes a given struct to JSON bytes.
func SerializeToJSON(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

// DeserializeFromJSON deserializes JSON bytes into a given struct.
func DeserializeFromJSON(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

// IsValidUUID checks if a string is a valid UUID.
func IsValidUUID(u string) bool {
	_, err := uuid.Parse(u)
	return err == nil
}

// ErrorResponse represents the structure for error responses.
type ErrorResponse struct {
	Error string `json:"error"`
}
