package main

import (
	"fmt"
	"github.com/OpenConnectOUSL/backend-api-v1/internal/data"
	"net/http"
	"time"
)

func (app *application) createIdeaHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "create a new idea")
}

func (app *application) showIdeaHandler(w http.ResponseWriter, r *http.Request) {

	id, err := app.readIDParam(r)
	if err != nil {
		http.NotFound(w, r)
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
		app.logger.Print(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
