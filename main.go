package main

import (
	"encoding/json"
	"log"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Producto struct {
	Id          int     `json:"id"`
	Name        string  `json:"name"`
	Quantity    int     `json:"quantity"`
	CodeValue   string  `json:"code_value"`
	IsPublished bool    `json:"is_published"`
	Expiration  string  `json:"expiration"`
	Price       float64 `json:"price"`
}

var prods []Producto

func main() {
	if err := levantarJson("products.json"); err != nil {
		log.Fatal(err)
	}
	router := gin.Default()
	router.GET("ping", func(ctx *gin.Context) {
		ctx.String(200, "pong")
	})
	productos := router.Group("/products")
	productos.GET("", getProds)
	productos.GET("/:id", getProd)
	productos.GET("/search", prodSearcher)
	router.Run()
}

func prodSearcher(c *gin.Context) {
	sPrice := c.Query("priceGt")
	price, err := strconv.ParseFloat(sPrice, 8)
	if err != nil {
		c.String(500, "No se pudo convertir el price a float")
		return
	}
	var mayores []Producto
	for _, prod := range prods {
		if prod.Price > price {
			mayores = append(mayores, prod)
		}
	}
	c.IndentedJSON(200, mayores)
}

func getProd(c *gin.Context) {
	defer func() {
		err := recover()
		if err != nil {
			c.String(404, "No existe ese producto")
		}
	}()
	id, err2 := strconv.Atoi(c.Param("id"))
	if err2 != nil {
		c.String(500, "Error en la conversion a int")
		return
	}
	prod := prods[id-1]
	c.IndentedJSON(200, prod)
}

func getProds(c *gin.Context) {
	c.IndentedJSON(200, prods)
}

func levantarJson(sheison string) error {
	bProducts, err := os.ReadFile(sheison)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(bProducts, &prods); err != nil {
		return err
	}
	return nil
}

//EJERCICIO 1
// import (
// 	"encoding/json"
// 	"io"
// 	"log"

// 	"github.com/gin-gonic/gin"
// )

// type Persona struct {
// 	Nombre   string `json:"nombre"`
// 	Apellido string `json:"apellido"`
// }

// func main() {
// 	router := gin.Default()

// 	router.GET("/ping", func(ctx *gin.Context) {
// 		ctx.String(200, "pong")
// 	})

// 	router.POST("/saludo", func(ctx *gin.Context) {
// 		myDecoder := json.NewDecoder(ctx.Request.Body)
// 		var alguien Persona

// 		for {
// 			if err := myDecoder.Decode(&alguien); err == io.EOF {
// 				break
// 			} else if err != nil {
// 				log.Fatal(err)
// 			}
// 		}

// 		ctx.String(200, "Hola "+alguien.Nombre+" "+alguien.Apellido)
// 	})

// 	router.Run()
// }
