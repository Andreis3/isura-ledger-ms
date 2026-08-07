package module

import (
	"net/http"

	"github.com/andreis3/isura-ledger-ms/internal/infra/dependency"
	"github.com/andreis3/isura-ledger-ms/internal/infra/factory"
	"github.com/andreis3/isura-ledger-ms/internal/transport/rest/middleware"
	"github.com/andreis3/isura-ledger-ms/internal/transport/rest/types"
	"github.com/andreis3/isura-ledger-ms/internal/util"
)

type AccountModule struct {
	baseDeps *dependency.BaseDeps
}

func NewAccountModule(
	baseDeps dependency.BaseDeps,
) *AccountModule {
	return &AccountModule{
		baseDeps: &baseDeps,
	}
}

func (m *AccountModule) Routes() types.RouteType {

	return types.RouteType{
		{
			Method: http.MethodPost,
			Path:   "/accounts",
			Handler: func(w http.ResponseWriter, r *http.Request) {
				factory.NewCreateAccountFactory(m.baseDeps).Handle(w, r)
			},
			Type: util.HandlerFunc,
			Middlewares: &types.Middlewares{
				middleware.Tracing(m.baseDeps.Tracer),
				middleware.Logging(m.baseDeps.Log.SlogJSON()),
				middleware.MetricsMiddleware(m.baseDeps.Prom),
			},
		},
	}
}
