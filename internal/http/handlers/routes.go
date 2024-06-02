package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/sgatu/ezmail/internal/domain/models/domain"
	internal_http "github.com/sgatu/ezmail/internal/http"
	"github.com/sgatu/ezmail/internal/http/common"
	"github.com/sgatu/ezmail/internal/http/handlers/auth/login"
	"github.com/sgatu/ezmail/internal/http/handlers/register"
)

func secureMiddleware() func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			currUser := r.Context().Value(internal_http.CurrentUserKey)
			if currUser == nil {
				fmt.Printf("%+v\n", currUser)
				common.ErrorResponse(common.UnauthorizedError(), w)
				return
			}
			h.ServeHTTP(w, r)
		})
	}
}

func SetupRoutes(r *chi.Mux, appContext *internal_http.AppContext) {
	login.LoginHandler(appContext, r)
	register.RegisterHandler(appContext, r)
	r.With(secureMiddleware()).Group(func(r chi.Router) {
		// after this authorized endpoints
		r.Get("/test", func(w http.ResponseWriter, r *http.Request) {
			errSave := appContext.DomainInfoRepository.Save(r.Context(), &domain.DomainInfo{
				Id:         appContext.SnowflakeNode.Generate().String(),
				DomainName: "google.es",
				UserId:     "xxxx",
				Created:    time.Now(),
				Validated:  false,
				DnsRecords: []domain.DnsRecord{{Type: "SPF", Value: "test", Status: domain.DNS_RECORD_STATUS_PENDING}, {Type: "SPF2", Value: "test2", Status: domain.DNS_RECORD_STATUS_PENDING}},
			})
			if errSave != nil {
				fmt.Printf("%+v\n", errSave)
			}
			// currUser, _ := r.Context().Value(internal_http.CurrentUserKey).(*user.User)
			// fmt.Fprintf(w, "Current user - Name: %s, Id: %s", currUser.Name, currUser.Id)
		})
	})
}
