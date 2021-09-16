package middleware

import (
	"UserAPI/UserDS"
	"context"
	"fmt"
	"net/http"
	"time"
)

type KeyUser struct {
}

func ValidateUser(ph http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		user := &UserDS.User{}
		if r.Method == http.MethodPost {
			user.CreatedOn = time.Now()
		}
		user.UpdatedOn = time.Now()
		err := user.FromJson(r.Body)
		if err != nil {
			http.Error(rw, fmt.Sprintln("Unable to parse request. ", err), http.StatusBadRequest)
			return
		}

		if len(user.FirstName) == 0 || len(user.LastName) == 0 {
			http.Error(rw, "Invalid Format.", http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), KeyUser{}, user)

		r = r.WithContext(ctx)

		ph.ServeHTTP(rw, r)
	})
}
