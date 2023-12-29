package main

import (
	"errors"
	"net/http"
	"strconv"

	"tsundokukeeper/internal/models"

	"github.com/julienschmidt/httprouter"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	userID := extractUserIDFromRequest(r)

	books, err := app.books.Finished(userID)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.render(w, r, http.StatusOK, "home.html", templateData{
		Books: books,
	})
}

func (app *application) about(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Display the origin of the meaning tsundoku and the app..."))
}

func (app *application) tsunbookView(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	userID := extractUserIDFromRequest(r)

	book, err := app.books.Get(userID, id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, r, err)
		}
		return
	}

	app.render(w, r, http.StatusOK, "view.html", templateData{
		Book: book,
	})
}

func (app *application) addBook(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Display form for adding a new book..."))
}

func (app *application) addBookPost(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Display added book..."))
}

func (app *application) stats(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("DIsplay different stats of finished/unfinished books..."))
}

func extractUserIDFromRequest(r *http.Request) int {
	return 0
}
