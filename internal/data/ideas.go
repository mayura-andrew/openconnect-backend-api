package data

import (
	"github.com/OpenConnectOUSL/backend-api-v1/internal/validator"
	"time"
)

type Idea struct {
	ID              int64     `json:"id"` // Unique identifier for the idea
	CreatedAt       time.Time `json:"created_at"`
	Title           string    `json:"title"`            // Title of the idea
	Description     string    `json:"description"`      // Detailed description of the idea
	Category        string    `json:"category"`         // Category of the idea
	Tags            []string  `json:"tags"`             // List of tags associated with the idea
	SubmittedBy     int       `json:"submitted_by"`     // User ID of the person who submitted the idea
	SubmittedAt     time.Time `json:"submitted_at"`     // Timestamp of when the idea was submitted
	Upvotes         int       `json:"upvotes"`          // Number of upvotes received
	Downvotes       int       `json:"downvotes"`        // Number of downvotes received
	Status          string    `json:"status"`           // Current status of the idea (e.g., pending, approved, rejected)
	Comments        []Comment `json:"comments"`         // List of comments on the idea
	InterestedUsers []int     `json:"interested_users"` // List of user IDs who are interested in the idea
	Version         int32     `json:"version"`
}

type Comment struct {
	ID          int       `json:"id"`           // Unique identifier for the comment
	IdeaID      int       `json:"idea_id"`      // ID of the idea the comment is related to
	CommentedBy int       `json:"commented_by"` // User ID of the person who made the comment
	Content     string    `json:"content"`      // Content of the comment
	CreatedAt   time.Time `json:"created_at"`   // Timestamp of when the comment was created
}

func ValidateIdea(v *validator.Validator, idea *Idea) {
	v.Check(idea.Title != "", "title", "must be provided")
	v.Check(len(idea.Title) <= 100, "title", "must not be more than 100 bytes long")

	v.Check(idea.Description != "", "description", "must be provided")
	v.Check(len(idea.Description) <= 1000, "description", "must not be more than 1000 bytes long")

	v.Check(idea.Category != "", "category", "must be provided")
	v.Check(len(idea.Category) <= 50, "category", "must not be more than 50 bytes long")

	v.Check(idea.Tags != nil, "tags", "must be provided")
	v.Check(len(idea.Tags) >= 1, "tags", "must contain at least one tag")
	v.Check(validator.Unique(idea.Tags), "tags", "must not contain duplicate values")

	v.Check(idea.SubmittedBy != 0, "submitted_by", "must be provided")
	v.Check(idea.SubmittedBy > 0, "submitted_by", "must be a positive integer")
}
