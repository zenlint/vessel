package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
)

type StreamOutput struct {
	Events []map[string]interface{}
}

func NewStreamOutput() *StreamOutput {
	return &StreamOutput{
		Events: make([]map[string]interface{}, 0),
	}
}

func (so *StreamOutput) Write(p []byte) (int, error) {
	e := make(map[string]interface{})
	for _, line := range bytes.Split(p, []byte("\n")) {
		if len(line) == 0 {
			continue
		}
		if err := json.Unmarshal(line, &e); err != nil {
			return 0, fmt.Errorf("unmarshal '%s': %v", string(line), err)
		}
		if e["stream"] != nil {
			fmt.Print(e["stream"].(string))
		} else if e["error"] != nil {
			fmt.Println(e["error"].(string))
		} else if e["status"] != nil {
		} else {
			fmt.Println(e)
		}
		so.Events = append(so.Events, e)
	}
	return len(p), nil
}

type PrefixWriter struct {
	prefix []byte
	length int
	writer io.Writer
}

func NewPrefixWriter(prefix string, w io.Writer) *PrefixWriter {
	byts := []byte(prefix)
	return &PrefixWriter{byts, len(byts), w}
}

func (w *PrefixWriter) Write(p []byte) (int, error) {
	n, err := w.writer.Write(w.prefix)
	if err != nil {
		return n, err
	}
	n, err = w.writer.Write(p)
	return n + w.length, err
}

// MapToStrings convert a string map to a string slice.
func MapToStrings(m map[string]bool) []string {
	strs := make([]string, 0, len(m))
	for s := range m {
		strs = append(strs, s)
	}
	return strs
}
