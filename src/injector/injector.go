package injector

import (
	"github.com/tamurakeito/avocado_map_backend/src/presentation"
)

func InjectHttpHandler() presentation.HttpHandler {
	return presentation.NewHttpHandler()
}

func InjectWsHandler() presentation.WsHandler {
	return presentation.NewWsHandler()
}
