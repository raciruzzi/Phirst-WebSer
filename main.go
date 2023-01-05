package main

import (
	"encoding/json"
	"log"
	"os"
	"strconv"
	"time"

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

var idToPut = len(prods) + 1

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
	productos.POST("", nuevoProd)
	router.Run()
}

func nuevoProd(c *gin.Context) {
	var elNuevo Producto
	if err := c.ShouldBindJSON(&elNuevo); err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	//Primero el id, como dice que no es necesario, por mas de que me pase uno, lo voy a sobreescribir
	elNuevo.Id = idToPut
	//Ahora todos los datos deben tener algo, asumo que no nos pueden pasar Quantity 0 ni Precio 0
	if elNuevo.Name == "" || elNuevo.Quantity == 0 || elNuevo.CodeValue == "" || elNuevo.Expiration == "" || elNuevo.Price == 0 {
		c.JSON(400, gin.H{
			"error": "Solo podes dejar en blanco los campos 'id' y 'is_published'",
		})
		return
	}
	//Nos fijamos que el code_value sea unico
	for _, prod := range prods {
		if prod.CodeValue == elNuevo.CodeValue {
			c.JSON(400, gin.H{
				"error": "Ese code_value ya existe",
			})
			return
		}
	}
	//Ahora chequeamos que la fecha este bien
	if _, err := time.Parse("02/01/2006", elNuevo.Expiration); err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	//Listo todos los chequeos lo agregamos
	idToPut++
	prods = append(prods, elNuevo)
	c.JSON(200, elNuevo)
}

func prodSearcher(c *gin.Context) {
	sPrice := c.Query("priceGt")
	if sPrice == "" {
		sPrice = "0"
	}
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
