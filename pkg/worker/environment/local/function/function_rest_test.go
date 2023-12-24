package function

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/serverlessworkflow/sdk-go/v2/model"
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
      parameters:
      - name: test
        in: query
        required: false
        type: string
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
			fmt.Println("req", req.URL.Query().Get("name"))
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
		assert.Error(t, err)
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

		_, err := functionRest.Run(model.Object{})
		assert.Error(t, err)
	})

	t.Run("run", func(t *testing.T) {
		functionRest := FunctionRest{
			Http:      server.Client(),
			Operation: server.URL + "/function.json#test",
		}

		err := functionRest.Init()
		assert.Nil(t, err)

		dataOut, err := functionRest.Run(model.Object{})
		assert.Nil(t, err)
		assert.Equal(t, model.FromMap(map[string]any{"test": "test"}), dataOut)
	})
}
