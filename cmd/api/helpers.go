package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/OpenConnectOUSL/backend-api-v1/internal/utils"
	"github.com/OpenConnectOUSL/backend-api-v1/internal/validator"
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

func (app *application) processAndSavePDF(inputBase64 string, w http.ResponseWriter, r *http.Request) (string, error) {
	pdfData, err := base64.StdEncoding.DecodeString(inputBase64)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return "no key", err
	}

	if !app.isPDF(pdfData) {
		app.badRequestResponse(w, r, fmt.Errorf("invalid pdf file"))
		return "no key", err
	}

	const maxPDFSize = 5 * 1024 * 1024
	if len(pdfData) > maxPDFSize {
		app.badRequestResponse(w, r, fmt.Errorf("pdf file size must be less than 5MB"))
		return "no key", err
	}

	// save the pdf to a file
	uploadsDir := "uploads"
	err = os.MkdirAll(uploadsDir, 0755)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return "no key", err
	}

	// generate a random filename
	unique := utils.GenerateUUID()
	filenameWithId := unique + ".pdf"
	pdfPath := filepath.Join(uploadsDir, filenameWithId)

	err = os.WriteFile(pdfPath, pdfData, 0644)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return "no key", err
	}
	return unique, nil
}

func (app *application) readString(qs url.Values, key string, defaultValue string) string {
	s := qs.Get(key)

	if s == "" {
		return defaultValue
	}

	return s
}

func (app *application) readCSV(qs url.Values, key string, defaultValue []string) []string {
	csv := qs.Get(key)

	if csv == "" {
		return defaultValue
	}

	return strings.Split(csv, ",")
}

func (app *application) readInt(qs url.Values, key string, defaultValue int, v *validator.Validator) int {
	s := qs.Get(key)

	if s == "" {
		return defaultValue
	}
	i, err := strconv.Atoi(s)
	if err != nil {
		v.AddError(key, "must be an integer")
		return defaultValue
	}
	
	return i
}
