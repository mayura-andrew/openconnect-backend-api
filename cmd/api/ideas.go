package main

import (
	"fmt"
	"github.com/OpenConnectOUSL/backend-api-v1/internal/data"
	"github.com/OpenConnectOUSL/backend-api-v1/internal/validator"
	"net/http"
	"time"
)

func (app *application) createIdeaHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title       string   `json:"title"`        // Title of the idea
		Description string   `json:"description"`  // Detailed description of the idea
		Category    string   `json:"category"`     // Category of the idea
		Tags        []string `json:"tags"`         // List of tags associated with the idea
		SubmittedBy int      `json:"submitted_by"` // User ID of the person who submitted the idea
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	idea := &data.Idea{
		Title:       input.Title,
		Description: input.Description,
		Category:    input.Category,
		Tags:        input.Tags,
		SubmittedBy: input.SubmittedBy,
	}

	v := validator.New()

	if data.ValidateIdea(v, idea); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	fmt.Fprintf(w, "%v\n", input)
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
