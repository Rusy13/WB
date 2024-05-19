package routes

import (
	"net/http"

	"WB/internal/middleware"
	"WB/internal/order/delivery"
	"github.com/gorilla/mux"
)

func GetRouter(handlers *delivery.OrderDelivery, mw *middleware.Middleware) *mux.Router {
	router := mux.NewRouter()
	assignRoutes(router, handlers)
	assignMiddleware(router, mw)
	return router
}

func assignRoutes(router *mux.Router, handlers *delivery.OrderDelivery) {
	router.HandleFunc("/order", handlers.AddOrder).Methods(http.MethodPost)
	router.HandleFunc("/order/{id}", handlers.GetOrder).Methods(http.MethodGet)

}

func assignMiddleware(router *mux.Router, mw *middleware.Middleware) {
	router.Use(mw.AccessLog)
	//router.Use(mw.Auth)
}
