package main

import (
	"context"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gymshark/go-hasher"
)

type Response struct {
	Word string `json:"word"`
	Hash string `json:"hash"`
	Type string `json:"type"`
}

var ginLambda *ginadapter.GinLambda

func main() {
	// create gin instance
	g := gin.Default()

	// set cors policy
	g.Use(cors.Default())

	// setup routes
	SetupRoutes(g)

	// set gin lambda instance
	ginLambda = ginadapter.New(g)

	// start lambda
	lambda.Start(Handler)
}

func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return ginLambda.ProxyWithContext(ctx, request)
}

func SetupRoutes(g *gin.Engine) {
	v1 := g.Group("api/v1")

	// ping
	v1.GET("ping", PingHandler())

	// hash
	v1.GET("hash/:word", HashHandler())
}

func PingHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Status(http.StatusNoContent)
	}
}

func HashHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// get word param
		word := c.Param("word")

		// calculate hash
		hash := GetHash(word)

		// build response
		res := Response{
			Word: word,
			Hash: hash,
			Type: "SHA_256",
		}

		// return response
		c.JSON(
			http.StatusOK,
			res,
		)
	}
}

func GetHash(word string) string {
	// call library
	hash := hasher.
		Sha256([]byte(word)).
		Hex()

	return hash
}
