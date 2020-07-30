package middlewares

import (
	"net/http"

	"github.com/seknox/trasa/core/logs"
	"github.com/seknox/trasa/models"
)

//This middleware should be placed on outermost/topmost place
type InAppTrail struct {
}

func (iat InAppTrail) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			userContext := r.Context().Value("user").(models.UserContext)
			next.ServeHTTP(w, r)
			description := w.Header().Get("trailDescription")
			status := w.Header().Get("trailStatus")
			w.Header().Del("trailIntent")
			w.Header().Del("trailStatus")
			//logrus.Tracef(`intent is %s `, description)

			//TODO handle file download
			if description != "" && r.Method != "GET" {
				go logs.Store.LogInAppTrail(r.RemoteAddr, r.UserAgent(), description, userContext.User, status == "success")
			}

		})
}

//func (iat InAppTrail) Handler(next http.Handler) http.Handler {
//	return http.HandlerFunc(
//		func(w http.ResponseWriter, r *http.Request) {
//
//
//			userContext := r.Context().Value("user").(models.UserContext)
//
//			rec:=httptest.NewRecorder()
//			next.ServeHTTP(rec, r)
//
//
//
//			description := rec.Header().Get("trailDescription")
//			status := rec.Header().Get("trailStatus")
//			rec.Header().Del("trailDescription")
//			rec.Header().Del("trailStatus")
//
//			//logrus.Tracef(`intent is %s `, description)
//
//			//TODO handle file download
//			if description != "" && r.Method != "GET" {
//				go logs.Store.LogInAppTrail(r.RemoteAddr, r.UserAgent(), description, userContext.User, status == "success")
//			}
//			for k, v := range rec.Header() {
//				w.Header()[k] = v
//			}
//			w.WriteHeader(rec.Code)
//			rec.Body.WriteTo(w)
//
//		})
//}
