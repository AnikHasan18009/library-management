package main

import "library-service/app"

func main() {
	app := app.NewApplication()
	app.Init()
	app.Run()
	app.Wait()
	app.Cleanup()
}
