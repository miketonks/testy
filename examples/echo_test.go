package examples

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"testy"
)

func TestEchoClient(t *testing.T) {
	handler := EchoHandler()
	api := testy.New(handler)

	response := api.Get("/hello")
	assert.Equal(t, 200, response.StatusCode, "OK response is expected")
	assert.Equal(t, "hello, world!", response.String())

	response = api.SetHeader("X-UserName", "bob").Get("/hello")
	assert.Equal(t, 200, response.StatusCode, "OK response is expected")
	assert.Equal(t, "hello, bob!", response.String())
	fmt.Println(response.String())
}
