package function

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	openapiloads "github.com/go-openapi/loads"
	"github.com/go-openapi/spec"

	"github.com/galgotech/fermions-workflow/pkg/worker/data"
)

var ErrorOperationNotFound = errors.New("operation not found")
var ErrorOperationNotInitialize = errors.New("operation not initialize")

func newRest(operation string) *FunctionRest {
	return &FunctionRest{
		Http:      &http.Client{Timeout: time.Duration(1) * time.Second},
		Operation: operation,
	}
}

type FunctionRest struct {
	Http      *http.Client
	Operation string
	method    string
	url       string
	op        *spec.Operation
}

func (w *FunctionRest) Init() error {
	document, err := openapiloads.Spec(w.Operation)
	if err != nil {
		return err
	}

	operationParse, err := url.Parse(w.Operation)
	if err != nil {
		return err
	}

	var ok bool
	w.method, w.url, w.op, ok = document.Analyzer.OperationForName(operationParse.Fragment)
	if !ok {
		return ErrorOperationNotFound
	}

	w.url = fmt.Sprintf("%s%s%s", document.Host(), document.BasePath(), w.url)

	return nil
}

func (w *FunctionRest) Run(dataInput data.Data[any]) (data.Data[any], error) {
	if w.method == "" || w.url == "" || w.op == nil {
		return nil, ErrorOperationNotInitialize
	}

	req, err := http.NewRequest(w.method, w.url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := w.Http.Do(req)
	if err != nil {
		return nil, err
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	dataOutput := data.Data[any]{}
	err = json.Unmarshal(bodyBytes, &dataOutput)
	if err != nil {
		return nil, err
	}

	return dataOutput, nil
}
