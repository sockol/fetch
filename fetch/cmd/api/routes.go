package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	mux := httprouter.New()

	mux.NotFound = http.HandlerFunc(app.notFound)
	mux.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowed)

	mux.GET("/", app.health)
	mux.GET("/receipts/:id/points", app.receiptPoints)

	mux.POST("/receipts/process", app.receiptsProcess)

	return app.logAccess(app.recoverPanic(mux))
}
