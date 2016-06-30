package slog

type HandlerConfig struct {
	Levels  []string
	Handler Handler
}

type Config struct {
	Handlers []HandlerConfig
}

func LoadConfig(config Config) {
	newHandlers := make(map[string][]Handler)
	for _, handler := range config.Handlers {
		for _, level := range handler.Levels {
			if levelHandlers := handlers[level]; levelHandlers != nil {
				newHandlers[level] = append(levelHandlers, handler.Handler)
			} else {
				newHandlers[level] = []Handler{handler.Handler}
			}
		}
	}
	handlers = newHandlers
}
