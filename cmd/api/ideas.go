package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/OpenConnectOUSL/backend-api-v1/internal/data"
	"github.com/OpenConnectOUSL/backend-api-v1/internal/utils"
	"github.com/OpenConnectOUSL/backend-api-v1/internal/validator"
)

// func (app *application) createIdeaHandler(w http.ResponseWriter, r *http.Request) {
// 	var input struct {
// 		Title       string   `json:"title"`
// 		Description string   `json:"description"`
// 		Pdf         string   `json:"pdf"`
// 		Category    string   `json:"category"`
// 		Tags        []string `json:"tags"`
// 		SubmittedBy int      `json:"submitted_by"`
// 	}

// 	err := app.readJSON(w, r, &input)
// 	if err != nil {
// 		app.badRequestResponse(w, r, err)
// 		return
// 	}

// 	idea := &data.Idea{
// 		Title:       input.Title,
// 		Description: input.Description,
// 		Pdf:         input.Pdf,
// 		Category:    input.Category,
// 		Tags:        input.Tags,
// 		SubmittedBy: input.SubmittedBy,
// 	}

// 	v := validator.New()

// 	if data.ValidateIdea(v, idea); !v.Valid() {
// 		app.failedValidationResponse(w, r, v.Errors)
// 		return
// 	}

// 	fmt.Fprintf(w, "%v\n", input)
// }

func (app *application) createIdeaHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the multipart form
	err := r.ParseMultipartForm(128 * 1024) // Adjusted to 128MB
    if err != nil {
        log.Printf("Error parsing multipart form: %v", err) // Added logging for debugging
        if err == http.ErrMissingFile {
            app.badRequestResponse(w, r, errors.New("no file uploaded"))
            return
        }
        app.serverErrorResponse(w, r, err)
        return
    }

	// Get the form values
	title := r.FormValue("title")
	description := r.FormValue("description")
	category := r.FormValue("category")
	submittedByStr := r.FormValue("submitted_by")

	// Get the tags from the form
	tags := r.Form["tags"]

	if errMap := validator.ValidateRequiredFields(title, description, category, tags, submittedByStr); len(errMap) > 0 {
		err := ValidationError{Errors: errMap}
		app.badRequestResponse(w, r, err)
		return
	}
	
	// Convert the submitted_by string to int
	submittedBy, err := strconv.Atoi(submittedByStr)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// Get the PDF file from the form
	pdfFile, header, err := r.FormFile("pdf")
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	defer pdfFile.Close()

	// Create the uploads directory if it doesn't exist
	uploadsDir := "uploads"
	err = os.MkdirAll(uploadsDir, 0755)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	uniqueID := utils.GenerateUUID()
	filenameWithID := uniqueID + "_" + filepath.Base(header.Filename)

	if err := validator.ValidatePDFFile(header); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// Create a new file in the uploads directory
	pdfPath := filepath.Join(uploadsDir, filenameWithID)
	if err != savePDFFile(pdfFile, pdfPath) {
		app.serverErrorResponse(w, r, err)
		return
	}

	idea := &data.Idea{
		Title:           title,
		Description:     description,
		Category:        category,
		Pdf:         	 pdfPath,
		Tags:            tags,
		SubmittedBy:     submittedBy,
	}

	fmt.Println(pdfPath)
	v := validator.New()

	if data.ValidateIdea(v, idea); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	// Save the idea to the database or perform other necessary operations
	// ...

	fmt.Fprintf(w, "%v\n", "Idea submitted successfully")

}



func (app *application) showIdeaHandler(w http.ResponseWriter, r *http.Request) {

	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	idea := data.Idea{
		ID:              id,
		CreatedAt:       time.Now(),
		Title:           "New Research Idea",
		Description:     "Detailed description of the new research idea.",
		Category:        "Technology",
		Tags:            []string{"AI", "Machine Learning"},
		SubmittedBy:     101,
		SubmittedAt:     time.Now(),
		Upvotes:         10,
		Downvotes:       2,
		Status:          "pending",
		Comments:        []data.Comment{},
		InterestedUsers: []int{102, 103},
		Version:         1,
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"idea": idea}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
