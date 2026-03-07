package handlers

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
)

// UUIDRequest is the JSON body for UUID tool endpoints that require input.
type UUIDRequest struct {
	Value string `json:"value"`
}

// UUIDGenerateResponse is the JSON response for single UUID generation.
type UUIDGenerateResponse struct {
	Result string `json:"result"`
}

// BulkGenerateRequest is the JSON body for bulk UUID generation.
type BulkGenerateRequest struct {
	Version int `json:"version"`
	Count   int `json:"count"`
}

// BulkGenerateResponse is the JSON response for bulk UUID generation.
type BulkGenerateResponse struct {
	Results []string `json:"results"`
}

// UUIDValidateResponse is the JSON response for UUID validation.
type UUIDValidateResponse struct {
	Valid   bool   `json:"valid"`
	Version int    `json:"version,omitempty"`
	Variant string `json:"variant,omitempty"`
	Error   string `json:"error,omitempty"`
}

// UUIDParseResponse is the JSON response for UUID parsing with metadata.
type UUIDParseResponse struct {
	Valid     bool    `json:"valid"`
	Version   int     `json:"version,omitempty"`
	Variant   string  `json:"variant,omitempty"`
	Timestamp *string `json:"timestamp,omitempty"`
	ClockSeq  *int    `json:"clockSeq,omitempty"`
	Node      *string `json:"node,omitempty"`
	Error     string  `json:"error,omitempty"`
}

// FormatUUIDRequest is the JSON body for UUID formatting.
type FormatUUIDRequest struct {
	Value  string `json:"value"`
	Format string `json:"format"`
}

// GenerateUUIDv4 generates a random UUID v4.
func GenerateUUIDv4(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	id := uuid.New()
	writeJSON(w, UUIDGenerateResponse{Result: id.String()})
}

// GenerateUUIDv7 generates a time-ordered UUID v7.
func GenerateUUIDv7(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	id, err := uuid.NewV7()
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to generate UUID v7: %v", err), http.StatusInternalServerError)
		return
	}
	writeJSON(w, UUIDGenerateResponse{Result: id.String()})
}

// BulkGenerateUUID generates multiple UUIDs of a specified version.
func BulkGenerateUUID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req BulkGenerateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}
	if req.Count < 1 {
		req.Count = 1
	}
	if req.Count > 1000 {
		req.Count = 1000
	}
	if req.Version != 4 && req.Version != 7 {
		req.Version = 4
	}

	results := make([]string, 0, req.Count)
	for i := 0; i < req.Count; i++ {
		var id uuid.UUID
		var err error
		if req.Version == 7 {
			id, err = uuid.NewV7()
			if err != nil {
				http.Error(w, fmt.Sprintf("failed to generate UUID v7: %v", err), http.StatusInternalServerError)
				return
			}
		} else {
			id = uuid.New()
		}
		results = append(results, id.String())
	}
	writeJSON(w, BulkGenerateResponse{Results: results})
}

// variantName returns a human-readable name for the UUID variant.
func variantName(v uuid.Variant) string {
	switch v {
	case uuid.RFC4122:
		return "RFC 4122"
	case uuid.Reserved:
		return "Reserved"
	case uuid.Microsoft:
		return "Microsoft"
	case uuid.Future:
		return "Future"
	default:
		return "Unknown"
	}
}

// ValidateUUID checks if a string is a valid UUID and returns version/variant info.
func ValidateUUID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req UUIDRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}
	input := strings.TrimSpace(req.Value)
	id, err := uuid.Parse(input)
	if err != nil {
		writeJSON(w, UUIDValidateResponse{
			Valid: false,
			Error: err.Error(),
		})
		return
	}
	writeJSON(w, UUIDValidateResponse{
		Valid:   true,
		Version: int(id.Version()),
		Variant: variantName(id.Variant()),
	})
}

// ParseUUID extracts detailed metadata from a UUID.
func ParseUUID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req UUIDRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}
	input := strings.TrimSpace(req.Value)
	id, err := uuid.Parse(input)
	if err != nil {
		writeJSON(w, UUIDParseResponse{
			Valid: false,
			Error: err.Error(),
		})
		return
	}

	resp := UUIDParseResponse{
		Valid:   true,
		Version: int(id.Version()),
		Variant: variantName(id.Variant()),
	}

	version := id.Version()
	if version == 1 {
		ts := id.Time()
		t := time.Unix(int64(ts-0x01b21dd213814000)/1e7, 0)
		tsStr := t.UTC().Format(time.RFC3339)
		resp.Timestamp = &tsStr
		clockSeq := int(id.ClockSequence())
		resp.ClockSeq = &clockSeq
		nodeID := id.NodeID()
		if len(nodeID) == 6 {
			nodeStr := fmt.Sprintf("%02x:%02x:%02x:%02x:%02x:%02x",
				nodeID[0], nodeID[1], nodeID[2], nodeID[3], nodeID[4], nodeID[5])
			resp.Node = &nodeStr
		}
	} else if version == 7 {
		bytes := id[:]
		msTimestamp := binary.BigEndian.Uint64(append([]byte{0, 0}, bytes[:6]...)) >> 16
		t := time.UnixMilli(int64(msTimestamp))
		tsStr := t.UTC().Format(time.RFC3339Nano)
		resp.Timestamp = &tsStr
	}

	writeJSON(w, resp)
}

// FormatUUID converts a UUID to different formats.
func FormatUUID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req FormatUUIDRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}
	input := strings.TrimSpace(req.Value)
	id, err := uuid.Parse(input)
	if err != nil {
		http.Error(w, fmt.Sprintf("invalid UUID: %v", err), http.StatusBadRequest)
		return
	}

	var result string
	switch strings.ToLower(req.Format) {
	case "no-hyphens", "nohyphens":
		result = strings.ReplaceAll(id.String(), "-", "")
	case "uppercase", "upper":
		result = strings.ToUpper(id.String())
	case "braces":
		result = "{" + id.String() + "}"
	case "urn":
		result = "urn:uuid:" + id.String()
	case "standard", "":
		result = id.String()
	default:
		http.Error(w, fmt.Sprintf("unknown format: %s", req.Format), http.StatusBadRequest)
		return
	}
	writeJSON(w, UUIDGenerateResponse{Result: result})
}
