package function

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/serverlessworkflow/sdk-go/v2/model"
	"github.com/stretchr/testify/assert"
)

var specFhub = `
name: "test"
specVersion: "1.0"
version: "v1"
serving: {
	http: {
		url: "%s/path-test"
	}
}
functions: {
  test: {
    input: {
      arg0: string
      arg1: string
    }
    output: {
      ok: bool
    }
  }
}
`

func TestFhub(t *testing.T) {
	var spec []byte
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		switch req.URL.Path {
		case "/function.json":
			if req.Method != http.MethodGet {
				rw.Write([]byte{})
				t.Failed()
				return
			}
			rw.Write(spec)
		case "/path-test/v1/test/test":
			if req.Method != http.MethodPost {
				rw.Write([]byte{})
				t.Failed()
				return
			}
			rw.Write([]byte(`{"ok": true}`))
		default:
			t.Failed()
		}
	}))
	defer server.Close()

	spec = []byte(fmt.Sprintf(specFhub, server.URL))

	t.Run("init operation not found", func(t *testing.T) {
		functionRest := FunctionFhub{
			Http:      server.Client(),
			Operation: server.URL + "/function.json#test",
		}

		err := functionRest.Init()
		assert.NoError(t, err)

		dataIn := model.FromInterface(map[string]any{
			"arg0": "test",
			"arg1": "test2",
		})
		dataOut, err := functionRest.Run(dataIn)
		assert.NoError(t, err)

		dataOut2 := model.FromInterface(map[string]any{"ok": true})
		assert.Equal(t, dataOut2, dataOut)

		dataIn = model.FromInterface(map[string]any{
			"arg0": "test",
			"arg2": "test2",
		})
		_, err = functionRest.Run(dataIn)
		assert.Error(t, err)
	})
}
