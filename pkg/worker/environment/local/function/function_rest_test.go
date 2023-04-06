package function

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/galgotech/fermions-workflow/pkg/worker/data"
	"github.com/stretchr/testify/assert"
)

var specOpenAPI = `
openapi: 2.0
info:
  version: 1.0.0
host: %s
basePath: /v1
schemes:
  - "https"
  - "http"
paths:
  /path-test:
    get:
      operationId: test
      responses:
        "200":
          description: "OK"
          content:
            application/json:
              schema:
                type: object
                properties:
                  name:
                    type: string
`

func TestRest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		switch req.URL.Path {
		case "/function.json":
			rw.Write([]byte(specOpenAPI))
		case "/v1/path-test":
			rw.Write([]byte("{\"test\":\"test\"}"))
		default:
			t.Failed()
		}
	}))
	defer server.Close()

	specOpenAPI = fmt.Sprintf(specOpenAPI, strings.TrimRight(server.URL, "/"))

	t.Run("init operation not found", func(t *testing.T) {
		functionRest := FunctionRest{
			Http:      server.Client(),
			Operation: server.URL + "/function.json#testNotFound",
		}

		err := functionRest.Init()
		assert.ErrorIs(t, ErrorOperationNotFound, err)
	})

	t.Run("init operation", func(t *testing.T) {
		functionRest := FunctionRest{
			Http:      server.Client(),
			Operation: server.URL + "/function.json#test",
		}

		err := functionRest.Init()
		assert.Nil(t, err)
	})

	t.Run("run not init", func(t *testing.T) {
		functionRest := FunctionRest{
			Http:      server.Client(),
			Operation: server.URL + "/function.json#test",
		}

		_, err := functionRest.Run(data.Data[any]{})
		assert.ErrorIs(t, ErrorOperationNotInitialize, err)
	})

	t.Run("run", func(t *testing.T) {
		functionRest := FunctionRest{
			Http:      server.Client(),
			Operation: server.URL + "/function.json#test",
		}

		err := functionRest.Init()
		assert.Nil(t, err)

		dataOut, err := functionRest.Run(data.Data[any]{})
		assert.Nil(t, err)
		assert.Equal(t, data.Data[any]{"test": "test"}, dataOut)
	})
}
