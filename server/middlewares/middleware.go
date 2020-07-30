package middlewares

import (
	"fmt"
	"net/http"
	"net/http/httputil"
)

type Dumper struct {
}

func (am Dumper) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		Dumpster(next.ServeHTTP).ServeHTTP(w, r)
		return
	})
}

func Dumpster(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		dump, err := httputil.DumpRequest(r, true)
		if err != nil {
			// http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
			fmt.Println("##########: ", err)
			return
		}
		fmt.Println("------------------------------------------------------------------------------")
		fmt.Println(string(dump))
		next(w, r)

	})

}

type AuthMiddleware struct {
}

func (am AuthMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			SessionValidator(Authorization(next.ServeHTTP)).ServeHTTP(w, r)
			return
		})
}
