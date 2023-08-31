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
	url       *url.URL
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
	method, uri, op, ok := document.Analyzer.OperationForName(operationParse.Fragment)
	if !ok {
		return errors.New("operation not found")
	}

	w.method = method
	w.op = op
	w.url, err = url.Parse(fmt.Sprintf("%s%s%s", document.Host(), document.BasePath(), uri))
	if err != nil {
		return err
	}

	return nil
}

func (w *FunctionRest) Run(dataIn data.Data[any]) (data.Data[any], error) {
	if w.method == "" || w.url == nil && w.op == nil {
		return nil, errors.New("operation not initialized")
	}

	url := *w.url
	for _, parameter := range w.op.Parameters {
		value, ok := dataIn[parameter.Name].(string)
		if !ok && parameter.Required {
			return nil, fmt.Errorf("not found parameter %q", value)
		}

		if parameter.In == "query" {
			url.Query().Add(parameter.Name, value)
		} else if parameter.In == "path" {
		} else if parameter.In == "body" && (w.method == "POST" || w.method == "PUT" || w.method == "PATCH") {
		}
	}

	req, err := http.NewRequest(w.method, w.url.String(), nil)
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

	dataOut := data.Data[any]{}
	err = json.Unmarshal(bodyBytes, &dataOut)
	if err != nil {
		return nil, err
	}

	return dataOut, nil
}
