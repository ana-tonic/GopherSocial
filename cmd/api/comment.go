package main

import (
	"net/http"

	"github.com/ana-tonic/gopher-social/internal/store"
)

// CreateCommentPayload represents the payload for creating a new comment
// @Description Comment creation payload
type CreateCommentPayload struct {
	Content string `json:"content" validate:"required,max=300" example:"This is a comment"`
}

// @Summary		Create a comment
// @Description	Creates a new comment on a post
// @Tags			posts,comments
// @Accept			json
// @Produce		json
// @Security		ApiKeyAuth
// @Param			postID	path		int						true	"Post ID"
// @Param			comment	body		CreateCommentPayload	true	"Comment content"
// @Success		201		{object}	store.Comment
// @Failure		400		{object}	nil
// @Failure		404		{object}	nil
// @Failure		500		{object}	nil
// @Router			/posts/{postID}/comments [post]
func (app *application) createCommentHandler(w http.ResponseWriter, r *http.Request) {
	post := getPostFromContext(r)
	userId := 2
	ctx := r.Context()

	var payload CreateCommentPayload

	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	comment := &store.Comment{
		PostID:  post.ID,
		UserID:  int64(userId),
		Content: payload.Content,
	}

	if err := app.store.Comments.Create(ctx, comment); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusCreated, comment); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}
