package controllers

import (
	"fmt"
	"minitwit/config"

	"github.com/gin-gonic/gin"
)

var HttpHandlers = []interface{}{
	loginHandlers,
	logoutHandlers,
	userHandlers,
	registerHandlers,
	timelineHandlers,
	staticHandlers,
	addMessageHandlers,
	simulationHandlers,
}

// HandleRESTRequests - handles the rest requests
func HandleRESTRequests() {

	router := gin.Default()
	router.SetTrustedProxies(nil)

	for _, handler := range HttpHandlers {
		handler.(func(engine *gin.Engine))(router)
	}

	router.Run(fmt.Sprintf(":%s", config.GetConfig().Server.Port))

}