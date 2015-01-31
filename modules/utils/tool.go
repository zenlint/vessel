package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
)

type StreamOutput struct {
	Events []map[string]string
}

func NewStreamOutput() *StreamOutput {
	return &StreamOutput{
		Events: make([]map[string]string, 0),
	}
}

func (so *StreamOutput) Write(p []byte) (int, error) {
	e := make(map[string]string)
	for _, line := range bytes.Split(p, []byte("\n")) {
		if len(line) == 0 {
			continue
		}
		if err := json.Unmarshal(line, &e); err != nil {
			return 0, err
		}
		fmt.Print(string(e["stream"]))
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
