package helpers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"runtime/debug"
)

func RespondWithError(w http.ResponseWriter, code int, err, reason string) {
	if code > 499 {
		slog.Error("Responding with 5xx error", "code", code, "err", err, "reason", reason, "stack", debug.Stack())
	}
	type errResponse struct {
		StatusMessage string `json:"statusMessage"`
		Error         string `json:"error"`
		Reason        string `json:"reason"`
	}
	RespondWithJSON(w, code, errResponse{http.StatusText(code), err, reason})
}

func RespondWithJSON(w http.ResponseWriter, code int, payload any) {
	data, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(500)
		slog.Error("Failed to marshal JSON response", "code", code, "paylod", payload, "stack", debug.Stack())
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(data)
}

// ParseJSON parses a json object from the stream r and stores it in the value pointed by v. It interprets the errors
// returned by Decode (if any), and converts it into human readable form if possible
func ParseJSON(r io.ReadCloser, v any, allowUnknownFields bool) error {
	d := json.NewDecoder(r)
	if !allowUnknownFields {
		d.DisallowUnknownFields()
	}
	err := d.Decode(v)
	if err != nil {
		var syntaxErr *json.SyntaxError
		var typeErr *json.UnmarshalTypeError
		var invalidUnmarshalErr *json.InvalidUnmarshalError

		switch {
		case errors.As(err, &syntaxErr):
			return fmt.Errorf("invalid JSON at offset %d", syntaxErr.Offset)
		case errors.Is(err, io.ErrUnexpectedEOF):
			return fmt.Errorf("invalid JSON syntax")
		case errors.As(err, &typeErr):
			if typeErr.Field != "" {
				return fmt.Errorf("invalid JSON type for field %q", typeErr.Field)
			}
			return fmt.Errorf("invalid JSON at offset %d", typeErr.Offset)
		case errors.Is(err, io.EOF):
			return fmt.Errorf("invalid JSON: body is empty")
		case errors.As(err, &invalidUnmarshalErr):
			slog.Error("json unmarshal failed", "reason", "invalid type was passed to it")
			panic(err)
		default:
			return err
		}
	}
	return nil
}
