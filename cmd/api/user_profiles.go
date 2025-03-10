package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/OpenConnectOUSL/backend-api-v1/internal/data"
	"github.com/OpenConnectOUSL/backend-api-v1/internal/validator"
	"github.com/google/uuid"
)

func (app *application) createProfileHandler(w http.ResponseWriter, r *http.Request) {

	user := app.contextGetUser(r)

	_, err := app.models.UserProfile.GetByUserID(user.ID)
	if err == nil {
		app.failedValidationResponse(w, r, map[string]string{"profile": "profile already exists for this user"})
		return
	} else if !errors.Is(err, data.ErrRecordNotFound) {
		app.serverErrorResponse(w, r, err)
		return
	}

	var input struct {
		Firstname string   `json:"firstname"`
		Lastname  string   `json:"lastname"`
		Avatar    string   `json:"avatar"`
		Title     string   `json:"title"`
		Bio       string   `json:"bio"`
		Faculty   string   `json:"faculty"`
		Program   string   `json:"program"`
		Degree    string   `json:"degree"`
		Year      string   `json:"year"`
		Uni       string   `json:"uni"`
		Mobile    string   `json:"mobile"`
		LinkedIn  string   `json:"linkedin"`
		GitHub    string   `json:"github"`
		FB        string   `json:"fb"`
		Skills    []string `json:"skills"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	avatarID, err := app.processAndSaveAvatar(input.Avatar, w, r)
	if err != nil {
		return
	}

	profile := &data.Profile{
		UserID:    user.ID,
		Firstname: input.Firstname,
		Lastname:  input.Lastname,
		Avatar:    avatarID,
		Title:     input.Title,
		Bio:       input.Bio,
		Faculty:   input.Faculty,
		Program:   input.Program,
		Degree:    input.Degree,
		Year:      input.Year,
		Uni:       input.Uni,
		Mobile:    input.Mobile,
		LinkedIn:  input.LinkedIn,
		GitHub:    input.GitHub,
		FB:        input.FB,
		Skills:    input.Skills,
	}

	// Validate the profile data
	v := validator.New()
	validateProfile(v, profile)
	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	// Insert the profile
	err = app.models.UserProfile.Insert(profile)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	user.HasProfileCreated = true
	err = app.models.Users.Update(user)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// Return the created profile
	err = app.writeJSON(w, http.StatusCreated, envelope{"profile": profile}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) updateProfileHandler(w http.ResponseWriter, r *http.Request) {
	// Get current user from context
	user := app.contextGetUser(r)

	// Get existing profile
	profile, err := app.models.UserProfile.GetByUserID(user.ID)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// Parse the request body
	var input struct {
		Firstname *string  `json:"firstname"`
		Lastname  *string  `json:"lastname"`
		Avatar    *string  `json:"avatar"`
		Title     *string  `json:"title"`
		Bio       *string  `json:"bio"`
		Faculty   *string  `json:"faculty"`
		Program   *string  `json:"program"`
		Degree    *string  `json:"degree"`
		Year      *string  `json:"year"`
		Uni       *string  `json:"uni"`
		Mobile    *string  `json:"mobile"`
		LinkedIn  *string  `json:"linkedin"`
		GitHub    *string  `json:"github"`
		FB        *string  `json:"fb"`
		Skills    []string `json:"skills"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// Update profile fields if provided
	if input.Firstname != nil {
		profile.Firstname = *input.Firstname
	}
	if input.Lastname != nil {
		profile.Lastname = *input.Lastname
	}
	if input.Avatar != nil && *input.Avatar != "" {
		// Process and save new avatar image
		avatarID, err := app.processAndSaveAvatar(*input.Avatar, w, r)
		if err != nil {
			return
		}
		profile.Avatar = avatarID
	}
	if input.Title != nil {
		profile.Title = *input.Title
	}
	if input.Bio != nil {
		profile.Bio = *input.Bio
	}
	if input.Faculty != nil {
		profile.Faculty = *input.Faculty
	}
	if input.Program != nil {
		profile.Program = *input.Program
	}
	if input.Degree != nil {
		profile.Degree = *input.Degree
	}
	if input.Year != nil {
		profile.Year = *input.Year

	}
	if input.Uni != nil {
		profile.Uni = *input.Uni
	}
	if input.Mobile != nil {
		profile.Mobile = *input.Mobile
	}
	if input.LinkedIn != nil {
		profile.LinkedIn = *input.LinkedIn
	}
	if input.GitHub != nil {
		profile.GitHub = *input.GitHub
	}
	if input.FB != nil {
		profile.FB = *input.FB
	}
	if input.Skills != nil {
		profile.Skills = input.Skills
	}

	// Validate the updated profile
	v := validator.New()
	validateProfile(v, profile)
	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	// Update the profile
	err = app.models.UserProfile.Update(profile)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// Return the updated profile
	err = app.writeJSON(w, http.StatusOK, envelope{"profile": profile}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deleteProfileHandler(w http.ResponseWriter, r *http.Request) {
	// Get current user from context
	user := app.contextGetUser(r)

	// Delete profile for user
	err := app.models.UserProfile.Delete(user.ID)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "profile successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) getProfileHandler(w http.ResponseWriter, r *http.Request) {
	// Get current user from context
	user := app.contextGetUser(r)

	fmt.Println(user)

	// Get full profile
	profile, err := app.models.UserProfile.GetFullProfile(user.ID)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			profile = &data.UserProfile{
				ID:        user.ID,
				Username:  user.UserName,
				Email:     user.Email,
				UserType:  user.UserType,
				CreatedAt: user.CreatedAt,
			}

			// Return an empty profile with basic user information
			err = app.writeJSON(w, http.StatusOK, envelope{"profile": profile}, nil)
			if err != nil {
				app.serverErrorResponse(w, r, err)
			}
			return // No profile found, return empty profile
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	response := map[string]interface{}{
		"profile":           profile,
		"hasProfileCreated": user.HasProfileCreated,
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"profile": response}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) getProfileByUsernameHandler(w http.ResponseWriter, r *http.Request) {
	// Get username from URL
	username := app.readStringParam(r, "username")

	// Get profile by username
	profile, err := app.models.UserProfile.GetByUsername(username)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"profile": profile}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) listProfilesWithIdeasHandler(w http.ResponseWriter, r *http.Request) {
	// Get query parameters
	qs := r.URL.Query()

	v := validator.New()
	limit := app.readInt(qs, "limit", 20, v)
	offset := app.readInt(qs, "offset", 0, v)

	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	// Get all profiles with ideas
	profilesWithIdeas, err := app.models.UserProfile.GetAllProfilesWithIdeas(limit, offset)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// Return the profiles with ideas
	err = app.writeJSON(w, http.StatusOK, envelope{"profiles": profilesWithIdeas, "count": len(profilesWithIdeas)}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) searchProfilesHandler(w http.ResponseWriter, r *http.Request) {
	// Get query parameters
	qs := r.URL.Query()
	query := app.readString(qs, "q", "")

	v := validator.New()
	limit := app.readInt(qs, "limit", 20, v)
	offset := app.readInt(qs, "offset", 0, v)

	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	// Search profiles
	profiles, err := app.models.UserProfile.Search(query, limit, offset)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// Return search results
	err = app.writeJSON(w, http.StatusOK, envelope{"profiles": profiles, "count": len(profiles)}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// Helper function to validate profile data
func validateProfile(v *validator.Validator, profile *data.Profile) {
	v.Check(profile.UserID != uuid.Nil, "user_id", "must be provided")

	if profile.Firstname != "" {
		v.Check(len(profile.Firstname) <= 100, "firstname", "must not exceed 100 characters")
	}

	if profile.Lastname != "" {
		v.Check(len(profile.Lastname) <= 100, "lastname", "must not exceed 100 characters")
	}

	if profile.Title != "" {
		v.Check(len(profile.Title) <= 100, "title", "must not exceed 100 characters")
	}

	if profile.Bio != "" {
		v.Check(len(profile.Bio) <= 1000, "bio", "must not exceed 1000 characters")
	}

	// Validate skills
	if len(profile.Skills) > 0 {
		v.Check(len(profile.Skills) <= 20, "skills", "cannot have more than 20 skills")

		for i, skill := range profile.Skills {
			v.Check(len(skill) <= 50, fmt.Sprintf("skills.%d", i), "must not exceed 50 characters")
		}
	}
}

func (app *application) getCurrentUserProfileHandler(w http.ResponseWriter, r *http.Request) {
	// Get the current authenticated user
	user := app.contextGetUser(r)

	// Retrieve the full profile for the user
	profile, err := app.models.UserProfile.GetFullProfile(user.ID)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			// Return an empty profile with basic user information
			profile = &data.UserProfile{
				ID:        user.ID,
				Username:  user.UserName,
				Email:     user.Email,
				UserType:  user.UserType,
				CreatedAt: user.CreatedAt,
			}
			err = app.writeJSON(w, http.StatusOK, envelope{"profile": profile}, nil)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	response := map[string]interface{}{
		"profile": profile,
	}

	if profile.Avatar != "" && profile.Avatar != "no key" {
		avatarBase64, err := app.getAvatarBase64(profile.Avatar)
		if err != nil {
			// Just log the error - don't add avatar to response
			app.logger.PrintError(err, nil)
		} else if avatarBase64 != "" {
			// Only add avatar to response if we actually got one
			response["avatarBase64"] = avatarBase64
		}
	}

	// Return the response
	err = app.writeJSON(w, http.StatusOK, envelope{"response": response}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) getProfileWithIdeasByUserIDHandler(w http.ResponseWriter, r *http.Request) {

	userIdParam := app.readStringParam(r, "userId")

	userID, err := uuid.Parse(userIdParam)
	if err != nil {
        app.badRequestResponse(w, r, fmt.Errorf("invalid user ID: %s", err.Error()))
		return
	}

	profile, err := app.models.UserProfile.GetFullProfile(userID)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	qs := r.URL.Query()
	v := validator.New()
	limit := app.readInt(qs, "limit", 10, v)
	offset := app.readInt(qs, "offset", 0, v)

	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	ideas, totalCount, err := app.models.Ideas.GetAllByUserID(userID, limit, offset)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	response := map[string]interface{}{
		"profile": profile,
		"ideas": ideas,
		"ideas_count": totalCount,
		"limit": limit,
		"offset": offset,
	}	

	if profile.Avatar != "" && profile.Avatar != "no key" {
        response["avatarURL"] = fmt.Sprintf("/v1/avatars/%s", profile.Avatar)
    }

	err = app.writeJSON(w, http.StatusOK, envelope{"response": response}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}
