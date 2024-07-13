package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
)

type envelope map[string]any

func (app *application) readIDParam(r *http.Request) (uuid.UUID, error) {
	params := httprouter.ParamsFromContext(r.Context())

	idStr:= params.ByName("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return uuid.Nil, errors.New("invalid id parameter")
	}
	return id, nil
}

func (app *application) writeJSON(w http.ResponseWriter, status int, data envelope, headers http.Header) error {
    js, err := json.Marshal(data)
    if err != nil {
        return err
    }
    js = append(js, '\n')

    for key, value := range headers {
        w.Header()[key] = value
    }
    w.Header().Set("Content-Type", "application/json")

    // Use http.StatusText to get the status text for the given status code
    w.Header().Set("Status", http.StatusText(status))

    _, err = w.Write(js)
    return err
}

func (app *application) readJSON(w http.ResponseWriter, r *http.Request, dst any) error {

	maxBytes := 1_048_576 * 5 // 1MB
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError
		var maxBytesError *http.MaxBytesError

		switch {
		case errors.As(err, &syntaxError):
			return fmt.Errorf("body contains badly-formed JSON (at character %d)", syntaxError.Offset)

		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("body contains badly-formed JSON")

		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("body contains incorrect JSON type for field %q", unmarshalTypeError.Field)
			}
			return fmt.Errorf("body contains  incorrect JSON type (at character %d", unmarshalTypeError.Offset)

		case errors.Is(err, io.EOF):
			return errors.New("body must not be empty")

		case errors.As(err, &invalidUnmarshalError):
			panic(err)
		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			return fmt.Errorf("body contains unknown key %s", fieldName)
		case errors.As(err, &maxBytesError):
			return fmt.Errorf("body must not be larger than %d bytes", maxBytesError.Limit)
		case errors.As(err, &invalidUnmarshalError):
			panic(err)
		default:
			return err
		}
	}
	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must only contain a single JSON value")
	}
	return nil
}

func (app *application) isPDF (data []byte) bool {
	return bytes.HasPrefix(data, []byte("%PDF-"))
}

func (app *application) isPDFSizeValid (data []byte, maxSize int) bool {
	return len(data) <= maxSize
}