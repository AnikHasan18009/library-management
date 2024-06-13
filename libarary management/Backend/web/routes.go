package web

import (
	handler "library-service/web/handlers"
	"library-service/web/middlewares"

	"net/http"
)

func InitRoutes(mux *http.ServeMux, manager *middlewares.Manager) {

	mux.Handle(
		"POST /readers",
		manager.With(
			http.HandlerFunc(handler.CreateReader), middlewares.LoggingMiddleware,
		),
	)

	mux.Handle(
		"POST /admins",
		manager.With(
			http.HandlerFunc(handler.CreateAdmin), middlewares.LoggingMiddleware,
		),
	)
	mux.Handle(
		"POST /admins/approve",
		manager.With(
			http.HandlerFunc(handler.ApproveReader), middlewares.LoggingMiddleware,
		),
	)

	mux.Handle(
		"POST /admins/approve",
		manager.With(
			http.HandlerFunc(handler.ApproveReader), middlewares.LoggingMiddleware,
		),
	)

	mux.Handle(
		"POST /hash",
		manager.With(
			http.HandlerFunc(handler.HashVal), middlewares.LoggingMiddleware,
		),
	)
	mux.Handle(
		"GET /users/logout",
		manager.With(
			http.HandlerFunc(handler.LogoutUser), middlewares.Authenticate,
		),
	)
	mux.Handle(
		"GET /users/refresh-token", manager.With(http.HandlerFunc(handler.RefreshToken), middlewares.RefreshTokenVerification),
	)

	mux.Handle(
		"POST /users/login",
		manager.With(
			http.HandlerFunc(handler.LoginUser), middlewares.LoggingMiddleware,
		),
	)

}
