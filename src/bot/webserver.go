package main

import (
	"log"
	"net/http"
	"time"

	slackbot "github.com/BeepBoopHQ/go-slackbot"
	"github.com/gorilla/mux"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

func NewRouter(bot *slackbot.Bot) *mux.Router {
	router := mux.NewRouter().StrictSlash(false)
	for _, route := range routes {

		var handler http.Handler
		handler = route.HandlerFunc
		handler = Logger(handler, route.Name)
		handler = Botter(handler, bot)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}
	return router
}

var routes = Routes{
	Route{
		"get",
		"GET",
		"/get",
		emptyHandler_func,
	},
	Route{
		"get2",
		"GET",
		"/ge",
		emptyHandler_func,
	},
}

func emptyHandler_func(rw http.ResponseWriter, req *http.Request) {
	log.Printf("emptyHandler_func")
}

func Botter(inner http.Handler, bot *slackbot.Bot) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		inner.ServeHTTP(w, r)
		if r.RequestURI == "/get" {
			bot.RTM.NewOutgoingMessage("Hello", "#ascii-art-channel")
		} else {
			log.Printf("noooo")
		}

	})
}

func Logger(inner http.Handler, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		inner.ServeHTTP(w, r)

		log.Printf(
			"%s\t%s\t%s\t%s",
			r.Method,
			r.RequestURI,
			name,
			time.Since(start),
		)
	})
}
