package main

// go build -buildmode=plugin -o protobuf_response.so protobuf_response.go

import (
	"net/http"

	v1 "protobuf_response/generated/server-faker/v1"

	"github.com/gin-gonic/gin"
)

type plgHandler struct{}

func (plg *plgHandler) Respond(c *gin.Context) {
	person := v1.Person{
		Name:  "John Doe",
		Id:    12345,
		Email: "john.doe@doe.te",
	}
	c.ProtoBuf(http.StatusOK, &person)
}

func main() {}

// exported
var PlgHandler plgHandler
