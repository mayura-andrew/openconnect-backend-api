package data

import "time"

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
