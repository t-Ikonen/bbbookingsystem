package main

import (
	"net/http"

	"github.com/justinas/nosurf"
)

// //WriteToConsole is middleware function
// func WriteToConsole(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		fmt.Println("Hit the page")
// 		next.ServeHTTP(w, r)
// 	})

// }
//NoSurf adds CSRT to all POST requests
func NoSurf(next http.Handler) http.Handler {
	csrtHandler := nosurf.New(next)

	csrtHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   appCnf.InProduction,
		SameSite: http.SameSiteLaxMode,
	})
	return csrtHandler
}

//SessionLoadloads and save the session on every request
func SessionLoad(next http.Handler) http.Handler {
	return session.LoadAndSave(next)
}
