package utils

import (
	"library-service/types"
	"net/http"
	"strconv"
	"strings"
)

func GetQueryParametersForOrders(r *http.Request, qp *types.OrdersQueryParams) {

	parameters := r.URL.Query()

	if x := parameters.Get("order_type"); x != "" {
		qp.OrderType = strings.ToLower(x)
	}

	if x := parameters.Get("order_key"); x != "" {
		qp.OrderKey = strings.ToLower(x)
	}

	if x := parameters.Get("page"); x != "" {
		qp.Page, _ = strconv.Atoi(x)
	}

	if x := parameters.Get("limit"); x != "" {
		qp.Limit, _ = strconv.Atoi(x)
	}

}
