package api

import (
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"os"
)

type responder func(http.ResponseWriter, *http.Request)

func requireAuth(routeHandler responder) responder {
	return func(w http.ResponseWriter, r *http.Request) {
		if len(r.Header["Authorization"]) != 1 {
			log.Printf("Failed auth (no header) on %s from %s", r.RequestURI, r.RemoteAddr)
			w.WriteHeader(401)
			w.Write([]byte("Unauthorized"))
			return
		}

		token := r.Header["Authorization"][0]
		err := bcrypt.CompareHashAndPassword([]byte(os.Getenv("API_TOKEN_HASH")), []byte(token))

		if err == nil {
			// They're authenticated, so let's go ahead and exec the route
			routeHandler(w, r)
		} else {
			log.Printf("Failed auth (bad auth: %s vs %s) on %s from %s", token, os.Getenv("API_TOKEN_HASH"), r.RequestURI, r.RemoteAddr)
			w.WriteHeader(401)
			w.Write([]byte("Unauthorized"))
		}
	}
}
