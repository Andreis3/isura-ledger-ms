package module

import (
	"net/http"

	"github.com/andreis3/isura-ledger-ms/internal/transport/rest/types"
	"github.com/andreis3/isura-ledger-ms/internal/util"
	"github.com/go-chi/chi/v5/middleware"
)

type PPROF struct{}

func NewPPROF() *PPROF {
	return &PPROF{}
}

func (m *PPROF) Routes() types.RouteType {
	return types.RouteType{
		{
			Method:      http.MethodGet,
			Path:        "/debug",
			Handler:     middleware.Profiler(),
			Type:        util.Mount,
			Middlewares: nil,
		},
	}
}
