package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type httpServerObj struct {
	port     string
	username string
	password string
}

func httpServer(port string, username string, password string) httpServerObj {
	return httpServerObj{port, username, password}
}

func (obj httpServerObj) run() {
	r := gin.Default()
	r.GET("/conf", func(c *gin.Context) {
		// c.JSON(200, gin.H{
		// 	"message": "pong",
		// })
		// c.String(http.StatusOK, commandDealConfPath)
		confMap := utilsConf().getConf(commandDealConfPath)
		if utilsConf().isExistKeyOfMap("system", confMap) == true {
			delete(confMap, "system")
		}
		c.IndentedJSON(http.StatusOK, confMap)
	})

	// r.POST("/conf", func(c *gin.Context) {
	// 	body, _ := ioutil.ReadAll(c.Request.Body)
	// 	var retStr = ""
	// 	if body != nil {
	// 		log.Println(string(body))
	// 		retStr = string(body)
	// 	}
	// 	c.String(http.StatusOK, retStr)
	// })

	r.GET("/reload", func(c *gin.Context) {
		// c.JSON(200, gin.H{
		// 	"message": "pong",
		// })
		// c.String(http.StatusOK, commandDealConfPath)
		shutdownAll()
		c.String(http.StatusOK, `{"status":"ok"}`)
	})

	r.Run("localhost:" + obj.port) // listen and serve on 0.0.0.0:8080
}
