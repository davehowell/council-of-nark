package mediator

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
)

// Serve processes strict LF-delimited JSON requests. A malformed frame receives
// an error response and does not terminate the channel; an oversized frame or
// transport failure does.
func (s *Session) Serve(reader io.Reader, writer io.Writer) error {
	scanner := bufio.NewScanner(reader)
	scanner.Buffer(make([]byte, 64*1024), 1024*1024)
	encoder := json.NewEncoder(writer)
	for scanner.Scan() {
		var request Request
		if err := json.Unmarshal(scanner.Bytes(), &request); err != nil {
			s.mu.Lock()
			s.sequence++
			response := Response{Sequence: s.sequence, OK: false, Error: fmt.Sprintf("invalid request frame: %v", err)}
			s.mu.Unlock()
			if err := encoder.Encode(response); err != nil {
				return err
			}
			continue
		}
		if err := encoder.Encode(s.Handle(request)); err != nil {
			return err
		}
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("read mediator channel: %w", err)
	}
	return nil
}
