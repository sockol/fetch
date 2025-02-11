package main

import (
	"api/internal/request"
	"api/internal/response"
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
)

func (app *application) receiptsProcess(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var receipt Receipt
	err := request.DecodeJSONStrict(w, r, &receipt)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	// Validate the input is not faulty (I assume it will be?).
	// Ideally we do validation in middleware but this doesn't need to be prod ready.
	_, err = receipt.GetPoints()
	if err != nil {
		app.badRequest(w, r, errors.New("failed validation"))
		return
	}

	id := app.db.Add(&receipt)
	data := map[string]string{
		"id": id.String(),
	}

	err = response.JSON(w, http.StatusOK, data)
	if err != nil {
		app.serverError(w, r, err)
	}
}

func (app *application) receiptPoints(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id, err := uuid.Parse(p.ByName("id"))
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	receipt := app.db.Get(id)
	if receipt == nil {
		app.notFound(w, r)
		return
	}
	points, err := receipt.GetPoints()
	if err != nil {
		app.badRequest(w, r, err)
		return
	}
	data := map[string]string{
		"points": fmt.Sprintf("%d", points),
	}

	err = response.JSON(w, http.StatusOK, data)
	if err != nil {
		app.serverError(w, r, err)
	}
}

func (app *application) health(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	data := map[string]string{
		"Status": "OK",
	}

	err := response.JSON(w, http.StatusOK, data)
	if err != nil {
		app.serverError(w, r, err)
	}
}
