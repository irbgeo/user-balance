package api

import (
	"context"
	"net/http"
	"strings"

	httpSwagger "github.com/swaggo/http-swagger"

	"github.com/irbgeo/user-balance/internal/service"
)

type svc interface {
	UpBalance(ctx context.Context, c service.BalanceChange) error
	DownBalance(ctx context.Context, c service.BalanceChange) error
	GetBalance(ctx context.Context, u service.User) (service.Balance, error)
}

func Routes(svc svc) {
	http.HandleFunc("/", handlerFunc(svc))
}

func handlerFunc(svc svc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "swagger") {
			httpSwagger.WrapHandler(w, r)
			return
		}

		switch r.Method {
		case http.MethodGet:
			getBalance(w, r, svc)
		case http.MethodPost:
			upBalance(w, r, svc)
		case http.MethodPut:
			downBalance(w, r, svc)
		}
	}
}
