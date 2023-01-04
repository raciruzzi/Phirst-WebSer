package main

import (
	"encoding/json"
	"io"
	"log"

	"github.com/gin-gonic/gin"
)

type Persona struct {
	Nombre   string `json:"nombre"`
	Apellido string `json:"apellido"`
}

func main() {
	router := gin.Default()

	router.GET("/ping", func(ctx *gin.Context) {
		ctx.String(200, "pong")
	})

	router.POST("/saludo", func(ctx *gin.Context) {
		myDecoder := json.NewDecoder(ctx.Request.Body)
		var alguien Persona

		for {
			if err := myDecoder.Decode(&alguien); err == io.EOF {
				break
			} else if err != nil {
				log.Fatal(err)
			}
		}

		ctx.String(200, "Hola "+alguien.Nombre+" "+alguien.Apellido)
	})

	router.Run()
}
