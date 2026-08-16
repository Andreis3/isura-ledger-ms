package module

import (
	"net/http"

	"github.com/andreis3/isura-ledger-ms/internal/infra/dependency"
	"github.com/andreis3/isura-ledger-ms/internal/infra/factory"
	"github.com/andreis3/isura-ledger-ms/internal/transport/rest/middleware"
	"github.com/andreis3/isura-ledger-ms/internal/transport/rest/types"
	"github.com/andreis3/isura-ledger-ms/internal/util"
)

type TransactionModule struct {
	baseDeps *dependency.BaseDeps
}

func NewTransactionModule(
	baseDeps dependency.BaseDeps,
) *TransactionModule {
	return &TransactionModule{
		baseDeps: &baseDeps,
	}
}

func (m *TransactionModule) Routes() types.RouteType {

	return types.RouteType{
		{
			Method: http.MethodPost,
			Path:   "/transactions",
			Handler: func(w http.ResponseWriter, r *http.Request) {
				factory.NewCreateTransactionFactory(m.baseDeps).Handle(w, r)
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
