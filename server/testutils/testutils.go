package testutils

import (
	"context"
	"github.com/seknox/trasa/server/models"
	"net/http"
)

// SddTestUserContext is a middleware that adds  mock userContext
func AddTestUserContext(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		userContext := models.UserContext{
			User: &models.User{
				ID:         "13c45cfb-72ca-4177-b968-03604cab6a27",
				OrgID:      "153f7582-5ae2-46ba-8c1c-79ef73fe296e",
				UserName:   "root",
				FirstName:  "Bhargab",
				MiddleName: "",
				LastName:   "Acharya",
				Email:      "bhargab@seknox.com",
				Groups:     nil,
				UserRole:   "orgAdmin",
				Status:     true,
				IdpName:    "trasa",
			},
			Org: models.Org{
				ID:             "153f7582-5ae2-46ba-8c1c-79ef73fe296e",
				OrgName:        "Trasa",
				Domain:         "trasa.io",
				PrimaryContact: "",
				Timezone:       "Asia/Kathmandu",
				PhoneNumber:    "",
				CreatedAt:      0,
				PlatformBase:   "",
				License:        models.License{},
			},
			DeviceID:  "",
			BrowserID: "",
		}
		ctx := context.WithValue(r.Context(), "user", userContext)
		next(w, r.WithContext(ctx))

	})

}
