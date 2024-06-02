package handlers

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
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
			/*			domainInfo := &domain.DomainInfo{
							Id:         appContext.SnowflakeNode.Generate().String(),
							DomainName: "google.es",
							UserId:     "xxxx",
							Created:    time.Now(),
							Validated:  false,
						}
						domainInfo.SetDnsRecords([]domain.DnsRecord{{Type: "SPF", Value: "test", Status: domain.DNS_RECORD_STATUS_PENDING}, {Type: "SPF2", Value: "test2", Status: domain.DNS_RECORD_STATUS_PENDING}})*/
			//			errSave := appContext.DomainInfoRepository.Save(r.Context(), domainInfo)
			getDomains, err := appContext.DomainInfoRepository.GetAllByUserId(r.Context(), "xxxx")
			if err == nil && len(getDomains) > 0 {
				getDomains[0].GetDnsRecords()
				fmt.Printf("DomainInfo: %+v\n", getDomains)
			}
			// currUser, _ := r.Context().Value(internal_http.CurrentUserKey).(*user.User)
			// fmt.Fprintf(w, "Current user - Name: %s, Id: %s", currUser.Name, currUser.Id)
		})
	})
}
