package app

import (
	"library-service/config"
	"library-service/db"
	"library-service/logger"
	"library-service/web"
	Redis "library-service/web/redis"
	"sync"
)

type Application struct {
	wg sync.WaitGroup
}

func NewApplication() *Application {
	return &Application{}
}

func (app *Application) Init() {
	config.LoadConfig()
	logger.SetupLogger(config.GetConfig().ServiceName)
	db.InitDB()
	db.Init()
	Redis.InitRedis()
	//utils.InitValidator()
}

func (app *Application) Run() {

	web.StartServer(&app.wg)
}

func (app *Application) Wait() {
	app.wg.Wait()
}

func (app *Application) Cleanup() {
	db.CloseDB()

}
