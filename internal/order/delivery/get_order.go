package delivery

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"WB/internal/order/filter"
	"WB/internal/pkg/response"
)

func (d *OrderDelivery) GetBanners(w http.ResponseWriter, r *http.Request) {

	params := r.URL.Query()
	filters, err := d.getFilters(params, w)
	if err != nil {
		return
	}
	banners, err := d.service.GetOrders(r.Context(), *filters)
	if err != nil {
		d.logger.Errorf("internal server error in getting banners: %v", err)
		response.WriteResponse(w, response.Error{Err: response.ErrInternal.Error()}, http.StatusInternalServerError, d.logger)
		return
	}

	response.WriteResponse(w, banners, http.StatusOK, d.logger)
}

func (d *OrderDelivery) getFilters(params url.Values, w http.ResponseWriter) (*filter.Filter, error) {
	featureID := params.Get("feature_id")
	var orderIDInt uint64
	var err error
	if featureID != "" {
		orderIDInt, err = d.parseParam("feature id", featureID, w)
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
		errText := fmt.Sprintf("%s must positive integer", paramName)
		response.WriteResponse(w, response.Error{Err: errText},
			http.StatusBadRequest, d.logger)
		return 0, fmt.Errorf("parse error")
	}
	return value, nil
}
