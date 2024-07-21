package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/OpenConnectOUSL/backend-api-v1/internal/data"
	"github.com/OpenConnectOUSL/backend-api-v1/internal/validator"
	"github.com/google/uuid"
)

func (app *application) createIdeaHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title       string   `json:"title"`
		Description string   `json:"description"`
		PDFBase64   string   `json:"pdfBase64"`
		Category    string   `json:"category"`
		Tags        []string `json:"tags"`
		SubmittedBy string `json:"submitted_by"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// Validate the input data 

	if len(input.SubmittedBy) != 36 {
		app.badRequestResponse(w, r, fmt.Errorf("submitted_by must be a valid UUID of length 36"))
		return
	}

	submittedBy, err := uuid.Parse(input.SubmittedBy)
	if err != nil {
		app.badRequestResponse(w, r, fmt.Errorf("submitted_by must be a valid UUID"))
		return
	}

	uniqueID, err := app.processAndSavePDF(input.PDFBase64, w, r)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	idea := &data.Idea{
		Title:       input.Title,
		Description: input.Description,
		Pdf:         uniqueID,
		Category:    input.Category,
		Tags:        input.Tags,
		SubmittedBy: submittedBy,
	}

	v := validator.New()

	if data.ValidateIdea(v, idea); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Ideas.Insert(idea)
 	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
 	}

 	headers := make(http.Header)
 	headers.Set("Location", fmt.Sprintf("/v1/ideas/%d", idea.ID))
 	// Save the idea to the database or perform other necessary operations
// 	// ...

	err = app.writeJSON(w, http.StatusCreated, envelope{"idea": idea}, headers)

 	if err != nil {
 		app.serverErrorResponse(w, r, err)
 	}

 	fmt.Fprintf(w, "%v\n", "Idea submitted successfully")

}

// func (app *application) createIdeaHandler(w http.ResponseWriter, r *http.Request) {
// 	// Parse the multipart form
// 	err := r.ParseMultipartForm(128 * 1024) // Adjusted to 128MB
//     if err != nil {
//         log.Printf("Error parsing multipart form: %v", err) // Added logging for debugging
//         if err == http.ErrMissingFile {
//             app.badRequestResponse(w, r, errors.New("no file uploaded"))
//             return
//         }
//         app.serverErrorResponse(w, r, err)
//         return
//     }

// 	// Get the form values
// 	title := r.FormValue("title")
// 	description := r.FormValue("description")
// 	category := r.FormValue("category")
// 	submittedByStr := r.FormValue("submitted_by")

// 	// Get the tags from the form
// 	tags := r.Form["tags"]

// 	if errMap := validator.ValidateRequiredFields(title, description, category, tags, submittedByStr); len(errMap) > 0 {
// 		err := ValidationError{Errors: errMap}
// 		app.badRequestResponse(w, r, err)
// 		return
// 	}
	
// 	// Convert the submitted_by string to int
// 	submittedBy, err := strconv.Atoi(submittedByStr)
// 	if err != nil {
// 		app.badRequestResponse(w, r, err)
// 		return
// 	}

// 	// Get the PDF file from the form
// 	pdfFile, header, err := r.FormFile("pdf")
// 	if err != nil {
// 		app.badRequestResponse(w, r, err)
// 		return
// 	}
// 	defer pdfFile.Close()

// 	// Create the uploads directory if it doesn't exist
// 	uploadsDir := "uploads"
// 	err = os.MkdirAll(uploadsDir, 0755)
// 	if err != nil {
// 		app.serverErrorResponse(w, r, err)
// 		return
// 	}

// 	uniqueID := utils.GenerateUUID()
// 	filenameWithID := uniqueID + "_" + filepath.Base(header.Filename)

// 	if err := validator.ValidatePDFFile(header); err != nil {
// 		app.badRequestResponse(w, r, err)
// 		return
// 	}

// 	// Create a new file in the uploads directory
// 	pdfPath := filepath.Join(uploadsDir, filenameWithID)
// 	if err != savePDFFile(pdfFile, pdfPath) {
// 		app.serverErrorResponse(w, r, err)
// 		return
// 	}

// 	idea := &data.Idea{
// 		Title:           title,
// 		Description:     description,
// 		Category:        category,
// 		Pdf:         	 uniqueID,
// 		Tags:            tags,
// 		SubmittedBy:     submittedBy,
// 	}

// 	// fmt.Println(idea.Pdf)
// 	v := validator.New()

// 	if data.ValidateIdea(v, idea); !v.Valid() {
// 		app.failedValidationResponse(w, r, v.Errors)
// 		return
// 	}


// 	err = app.models.Ideas.Insert(idea)
// 	if err != nil {
// 		app.serverErrorResponse(w, r, err)
// 		return
// 	}

// 	headers := make(http.Header)
// 	headers.Set("Location", fmt.Sprintf("/v1/ideas/%d", idea.ID))
// 	// Save the idea to the database or perform other necessary operations
// 	// ...

// 	err = app.writeJSON(w, http.StatusCreated, envelope{"idea": idea}, headers)

// 	if err != nil {
// 		app.serverErrorResponse(w, r, err)
// 	}

// 	fmt.Fprintf(w, "%v\n", "Idea submitted successfully")

// }



func (app *application) showIdeaHandler(w http.ResponseWriter, r *http.Request) {

	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	idea, err := app.models.Ideas.Get(id)
	if err != nil {
		switch {
			case errors.Is(err, data.ErrRecordNotFound):
				app.notFoundResponse(w, r)
			default:
				app.serverErrorResponse(w, r, err)
		}
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"idea": idea}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) updateIdeaHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	idea, err := app.models.Ideas.Get(id)
	if err != nil {
		switch {
			case errors.Is(err, data.ErrRecordNotFound):
				app.notFoundResponse(w, r)
			default:
				app.serverErrorResponse(w, r, err)
		}
		return
	}

	var input struct {
		Title       *string   `json:"title"`
		Description *string   `json:"description"`
		Category    *string   `json:"category"`
		Tags        []string `json:"tags"`
		PdfBase64   *string   `json:"pdfBase64"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}



	uniqueID, err := app.processAndSavePDF(*input.PdfBase64, w, r)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	if input.Title != nil {
		idea.Title = *input.Title
	}

	if input.Description != nil {
		idea.Description = *input.Description
	}

	if input.Category != nil {
		idea.Category = *input.Category
	}

	if input.Tags != nil {
		idea.Tags = input.Tags
	}

	if input.PdfBase64 != nil {
		idea.Pdf = uniqueID
	}

	v := validator.New()

	if data.ValidateIdea(v, idea); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Ideas.Update(idea)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"idea": idea}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deleteIdeaHandler(w http.ResponseWriter, r *http.Request) {

	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.models.Ideas.Delete(id)
	if err != nil {
		switch {
			case errors.Is(err, data.ErrRecordNotFound):
				app.notFoundResponse(w, r)
			default:
				app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "idea deleted successfully"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}