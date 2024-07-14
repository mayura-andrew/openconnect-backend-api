package data

import (
	"database/sql"
	"errors"
	"time"

	"github.com/OpenConnectOUSL/backend-api-v1/internal/validator"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type Idea struct {
	ID              uuid.UUID     `json:"id"` // Unique identifier for the idea
	CreatedAt       time.Time `json:"created_at"`
	Title           string    `json:"title"`            // Title of the idea
	Description     string    `json:"description"`      // Detailed description of the idea
	Category        string    `json:"category"`         // Category of the idea
	Pdf             string    `json:"pdf"`              // PDF file of the idea
	Tags            []string  `json:"tags"`             // List of tags associated with the idea
	SubmittedBy     uuid.UUID       `json:"submitted_by"`     // User ID of the person who submitted the idea
	UpdatedAt     time.Time `json:"updated_at"`     // Timestamp of when the idea was submitted
	Upvotes         int       `json:"upvotes"`          // Number of upvotes received
	Downvotes       int       `json:"downvotes"`        // Number of downvotes received
	Status          string    `json:"status"`           // Current status of the idea (e.g., pending, approved, rejected)
	Comments        []uuid.UUID `json:"comments"`         // List of comments on the idea
	InterestedUsers []uuid.UUID    `json:"interested_users"` // List of user IDs who are interested in the idea
	Version         int32     `json:"version"`
}

type Comment struct {
	ID          uuid.UUID       `json:"id"`           // Unique identifier for the comment
	IdeaID      uuid.UUID       `json:"idea_id"`      // ID of the idea the comment is related to
	CommentedBy uuid.UUID       `json:"commented_by"` // User ID of the person who made the comment
	Content     string    `json:"content"`      // Content of the comment
	CreatedAt   time.Time `json:"created_at"`   // Timestamp of when the comment was created
}

func ValidateIdea(v *validator.Validator, idea *Idea) {
	v.Check(idea.Title != "", "title", "must be provided")
	v.Check(len(idea.Title) <= 100, "title", "must not be more than 100 bytes long")

	v.Check(idea.Description != "", "description", "must be provided")
	v.Check(len(idea.Description) <= 1000, "description", "must not be more than 1000 bytes long")

	v.Check(idea.Pdf != "", "pdf", "must be provided")
	v.Check(len(idea.Pdf) > 0, "pdf", "must not be empty")

	// isPDF := strings.EqualFold(filepath.Ext(idea.Pdf), ".pdf")

	// v.Check(isPDF, "pdf", "must be a PDF file")

	// fileInfo, err := os.Stat(idea.Pdf)
	// if err != nil {
	// 	v.AddError("pdf", "unable to get file into")
	// } else {
	// 	const maxPDFSize = 5 * 1024 * 1024 // 5MB
	// 	v.Check(fileInfo.Size() <= maxPDFSize, "pdf", "file size must be less than 5MB")
	// }

	v.Check(idea.Category != "", "category", "must be provided")
	v.Check(len(idea.Category) <= 50, "category", "must not be more than 50 bytes long")

	v.Check(idea.Tags != nil, "tags", "must be provided")
	v.Check(len(idea.Tags) >= 1, "tags", "must contain at least one tag")
	v.Check(validator.Unique(idea.Tags), "tags", "must not contain duplicate values")

	// v.Check(idea.SubmittedBy != 0, "submitted_by", "must be provided")
	// v.Check(idea.SubmittedBy > 0, "submitted_by", "must be a positive integer")
	if !ValidateUUID(idea.SubmittedBy.String()) {
		v.AddError("submitted_by", "must be a valid UUID")
	}
}

func ValidateUUID(uuidStr string) bool {
	_, err := uuid.Parse(uuidStr)
	return err == nil
}
type IdeaModel struct {
	DB *sql.DB
}

func (i IdeaModel) Insert(idea *Idea) error {
	query := `INSERT INTO ideas (title, description, submitted_by, idea_source_id, category, tags)
    VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, created_at, version`

	args := []any{idea.Title, idea.Description, idea.SubmittedBy, idea.Pdf, idea.Category, pq.Array(idea.Tags)}

	return i.DB.QueryRow(query, args...).Scan(&idea.ID, &idea.CreatedAt, &idea.Version)

}

func (i IdeaModel) Update(idea *Idea) error {
	query := `UPDATE ideas 
	SET title = $1, description = $2, submitted_by = $3, idea_source_id = $4, category = $5, tags = $6, version = version + 1 WHERE id = $7 RETURNING version`

	args := []any{
		idea.Title,
		idea.Description,
		idea.SubmittedBy,
		idea.Pdf,
		idea.Category,
		pq.Array(idea.Tags),
		idea.ID,
	}

	return i.DB.QueryRow(query, args...).Scan(&idea.Version)
}

func (i IdeaModel) Get(id uuid.UUID) (*Idea, error) {
	// id should be a valid UUID

	query := `SELECT id, created_at, updated_at, title, description, submitted_by, idea_source_id, category, tags, upvotes, downvotes, status, comments,  interested_users, version FROM ideas WHERE id = $1`

	var idea Idea

	err := i.DB.QueryRow(query, id).Scan(
		&idea.ID,
		&idea.CreatedAt,
		&idea.UpdatedAt,
		&idea.Title,
		&idea.Description,
		&idea.SubmittedBy,
		&idea.Pdf,
		&idea.Category,
		pq.Array(&idea.Tags),
		&idea.Upvotes,
		&idea.Downvotes,
		&idea.Status,
		pq.Array(&idea.Comments),
		pq.Array(&idea.InterestedUsers),
		&idea.Version,
	)

	if err != nil {
		switch {
			case errors.Is(err, sql.ErrNoRows):
				return nil, ErrRecordNotFound
			default:
				return nil, err
		}
	}

	return &idea, nil
	
}

func (i IdeaModel) Delete(id uuid.UUID) error {
	query := `DELETE FROM ideas WHERE id = $1`

	result, err := i.DB.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}
