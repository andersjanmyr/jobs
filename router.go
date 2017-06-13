package main

import (
	"github.com/andersjanmyr/jobs/controllers"
	"github.com/gorilla/mux"
)

func setupRouter(router *mux.Route, controller controllers.RestController) *mux.Router {
	var subRouter = router.Subrouter()
	subRouter.HandleFunc("/", controller.Index).Methods("GET")
	subRouter.HandleFunc("/", controller.Create).Methods("POST")
	subRouter.HandleFunc("/{slug}", controller.Show).Methods("GET")
	subRouter.HandleFunc("/{slug}", controller.Update).Methods("PUT")
	subRouter.HandleFunc("/{slug}", controller.Destroy).Methods("DELETE")
	subRouter.HandleFunc("/{slug}/new", controller.New).Methods("GET")
	subRouter.HandleFunc("/{slug}/edit", controller.Edit).Methods("GET")
	return subRouter
}
