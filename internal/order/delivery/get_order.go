package delivery

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"WB/internal/order/filter"
	"WB/internal/pkg/response"
)

func (d *OrderDelivery) GetOrder(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	filters, err := d.getFilters(params, w)
	if err != nil {
		return
	}

	order, err := d.service.GetUserBanner(r.Context(), strconv.FormatUint(filters.OrderID, 10))
	if err != nil {
		d.logger.Errorf("error getting order: %v", err)
		response.WriteResponse(w, response.Error{Err: "order not found"}, http.StatusNotFound, d.logger)
		return
	}

	response.WriteResponse(w, order, http.StatusOK, d.logger)
}

func (d *OrderDelivery) getFilters(params url.Values, w http.ResponseWriter) (*filter.Filter, error) {
	orderID := params.Get("order_id")
	var orderIDInt uint64
	var err error
	if orderID != "" {
		orderIDInt, err = d.parseParam("order_id", orderID, w)
		if err != nil {
			return nil, err
		}
	}
	return &filter.Filter{
		OrderID: orderIDInt,
	}, nil
}

func (d *OrderDelivery) parseParam(paramName, paramValue string, w http.ResponseWriter) (uint64, error) {
	value, err := strconv.ParseUint(paramValue, 10, 64)
	if err != nil || value < 1 {
		d.logger.Errorf("error in %s conversion: %s", paramName, err)
		response.WriteResponse(w, response.Error{Err: paramName + " must be a positive integer"}, http.StatusBadRequest, d.logger)
		return 0, fmt.Errorf("parse error")
	}
	return value, nil
}
