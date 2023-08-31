package function

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	fhubModel "github.com/galgotech/fhub-go/model"
	"github.com/serverlessworkflow/sdk-go/v2/model"

	"github.com/galgotech/fermions-workflow/pkg/worker/data"
)

func newFhub(f model.Function) *FunctionFhub {
	return &FunctionFhub{
		Http:      &http.Client{Timeout: time.Duration(1) * time.Second},
		Operation: f.Operation,
	}
}

type FunctionFhub struct {
	Http      *http.Client
	Operation string
	url       string
	function  fhubModel.Function
}

func (w *FunctionFhub) Init() error {
	_, functionName, ok := strings.Cut(w.Operation, "#")
	if !ok {
		return errors.New("operation not specified")
	}

	resp, err := http.Get(w.Operation)
	if err != nil {
		return err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	fhub, err := fhubModel.UnmarshalBytes(body)
	if err != nil {
		return err
	}

	exists := false
	for name, function := range fhub.Functions {
		if name == functionName {
			exists = true
			w.function = function
			break
		}
	}

	if !exists {
		return errors.New("operation not found")
	}

	w.url = fmt.Sprintf("%s/%s/%s/%s", fhub.Serving.Http.Url, fhub.Version, fhub.Name, functionName)
	return nil
}

func (w *FunctionFhub) Run(dataIn data.Data[any]) (data.Data[any], error) {
	jsonData, err := dataIn.Marshal()
	if err != nil {
		return nil, err
	}

	if err := w.function.ValidateInput(jsonData); err != nil {
		return nil, fmt.Errorf("invalid input: %q", err)
	}

	req, err := http.NewRequest(http.MethodPost, w.url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	resp, err := w.Http.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if err := w.function.ValidateOutput(body); err != nil {
		return nil, fmt.Errorf("invalid output: %q", err)
	}

	dataOut := data.Data[any]{}
	err = dataOut.Unmarshal(body)
	if err != nil {
		return nil, err
	}

	return dataOut, nil
}
