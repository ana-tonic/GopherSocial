package main

import (
	"net/http"

	"github.com/ana-tonic/gopher-social/internal/store"
)

// @Summary		Fetches the user feed
// @Description	Fetches the user feed
// @Tags			feed
// @Accept			json
// @Produce		json
// @Param			limit	query		int			false	"Number of posts to return (default 20)"	minimum(1)	maximum(20)
// @Param			offset	query		int			false	"Offset for pagination (default 0)"			minimum(0)
// @Param			sort	query		string		false	"Sort order (asc or desc, default desc)"	Enums(asc,desc)
// @Param			search	query		string		false	"Search in title and content"
// @Param			tags	query		[]string	false	"Filter by tags (comma separated)"
// @Param			since	query		string		false	"Filter posts since date (format: 2006-01-02 15:04:05)"
// @Param			until	query		string		false	"Filter posts until date (format: 2006-01-02 15:04:05)"
// @Success		200		{object}	[]store.PostWithMetadata
// @Failure		400		{object}	error
// @Failure		500		{object}	error
// @Security		ApiKeyAuth
// @Router			/users/feed [get]
func (app *application) getUserFeedHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Get authenticated user ID from context
	// For now, we'll use a hardcoded ID for testing
	userID := int64(1) // This should come from auth context

	fq := store.PaginatedFeedQuery{
		Limit:  20,
		Offset: 0,
		Sort:   "desc",
	}

	fq, err := fq.Parse(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(fq); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	ctx := r.Context()

	feed, err := app.store.Posts.GetUserFeed(ctx, userID, fq)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, feed); err != nil {
		app.internalServerError(w, r, err)
	}
}
