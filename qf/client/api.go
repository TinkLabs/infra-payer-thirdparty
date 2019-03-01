// Package client provides a Qf client for invoking APIs across all resources
package client

import (
	qf "github.com/TinkLabs/payer-thirdparty/qf"

	"github.com/TinkLabs/payer-thirdparty/qf/charge"
	"github.com/TinkLabs/payer-thirdparty/qf/query"
	"github.com/TinkLabs/payer-thirdparty/qf/refund"
)

// API is the Qf client. It contains all the different resources available.
type API struct {
	// Charges is the client used to invoke /charge APIs.
	Charges *charge.Client
	// Refunds is the client used to invoke /refund APIs.
	Refunds *refund.Client
	// Querys is the client used to invoke /query APIs.
	Querys *query.Client
}

// Init initializes the Qf client with the appropriate secret key
// as well as providing the ability to override the backend as needed.
func (a *API) Init(key string, backends *qf.Backends) {
	if backends == nil {
		backends = &qf.Backends{
			API: qf.GetBackend(qf.APIBackend),
		}
	}

	a.Charges = &charge.Client{B: backends.API, AppCode: key}
	a.Refunds = &refund.Client{B: backends.API, AppCode: key}
	a.Querys = &query.Client{B: backends.API, AppCode: key}
}
