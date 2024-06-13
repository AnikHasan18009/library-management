package web

import (
	"fmt"
	"library-service/config"
	"library-service/web/middlewares"
	"library-service/web/swagger"
	"log/slog"
	"net/http"
	"sync"
)

func StartServer(wg *sync.WaitGroup) {
	manager := middlewares.NewManager()

	// manager.Use(
	// 	middlewares.Recover,
	// 	middlewares.Logger,
	// )

	mux := http.NewServeMux()

	InitRoutes(mux, manager)

	handler := middlewares.EnableCors(mux)

	swagger.SetupSwagger(mux, manager)

	wg.Add(1)

	go func() {
		defer wg.Done()

		conf := config.GetConfig()

		addr := fmt.Sprintf(":%d", conf.HttpPort)
		slog.Info(fmt.Sprintf("Listening at %s", addr))

		if err := http.ListenAndServe(addr, handler); err != nil {
			slog.Error(err.Error())
		}
	}()
}
