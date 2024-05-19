package delivery

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"WB/internal/order/delivery/dto"
	"WB/internal/order/storage"
	"WB/internal/pkg/response"
)

const adminRole = "admin"

func (d *OrderDelivery) AddOrder(w http.ResponseWriter, r *http.Request) {
	user, err := getUserFromContext(r.Context())
	if err != nil {
		d.logger.Errorf("error in getting user from context: %v", err)
		response.WriteResponse(w, response.Error{Err: response.ErrInternal.Error()}, http.StatusInternalServerError, d.logger)
		return
	}
	if user.Role != adminRole {
		d.logger.Errorf("user %d has got no access dor adding order", user.TagID)
		w.WriteHeader(http.StatusForbidden)
		return
	}
	rBody, err := io.ReadAll(r.Body)
	if err != nil {
		d.logger.Errorf("error in reading request body: %v", err)
		response.WriteResponse(w, response.Error{Err: response.ErrInternal.Error()}, http.StatusInternalServerError, d.logger)
		return
	}

	defer func() {
		err = r.Body.Close()
		if err != nil {
			d.logger.Errorf("error in closing request body")
		}
	}()

	var orderDTO dto.AddOrderDTO
	err = json.Unmarshal(rBody, &orderDTO)
	if err != nil {
		var jsonErr *json.SyntaxError
		if errors.As(err, &jsonErr) {
			d.logger.Errorf("invalid json: %s", string(rBody))
			response.WriteResponse(w, response.Error{Err: response.ErrInvalidJSON.Error()}, http.StatusBadRequest, d.logger)
			return
		}
		d.logger.Errorf("error in response body unmarshalling: %v", err)
		response.WriteResponse(w, response.Error{Err: response.ErrInternal.Error()}, http.StatusInternalServerError, d.logger)
		return
	}

	err = orderDTO.Validate()
	if err != nil {
		d.logger.Errorf("validation errors in adding order: %v", err)
		response.WriteResponse(w, response.Error{Err: err.Error()}, http.StatusBadRequest, d.logger)
		return
	}

	orderToAdd := dto.ConvertToOrder(orderDTO)
	addedOrder, err := d.service.AddOrder(r.Context(), orderToAdd)
	if err != nil {
		if errors.Is(err, storage.ErrDuplicateFeatureTag) {
			d.logger.Errorf("order with one of combibnations of featiure + tag already exists: %d, %v", orderDTO.OrderUID, orderDTO.Items)
			response.WriteResponse(w, response.Error{Err: err.Error()}, http.StatusBadRequest, d.logger)
			return
		}
		d.logger.Errorf("internal server error in adding order: %v", err)
		response.WriteResponse(w, response.Error{Err: response.ErrInternal.Error()}, http.StatusInternalServerError, d.logger)
		return
	}

	response.WriteResponse(w, dto.OrderResponse{UID: addedOrder.OrderUID}, http.StatusCreated, d.logger)
}
