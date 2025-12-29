package clictx

import "github.com/mxV03/warhousemanagementsystem/ent"

type App struct {
	client *ent.Client
}

var app *App

func Init(client *ent.Client) {
	if client == nil {
		panic("cli.Init: client is nil")
	}
	app = &App{client: client}
}

func AppCtx() *App {
	if app == nil {
		panic("cli not initialized")
	}
	return app
}

func (a *App) Client() *ent.Client {
	return a.client
}
