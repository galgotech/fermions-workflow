package function

import (
	"errors"

	"github.com/serverlessworkflow/sdk-go/v2/model"

	"github.com/galgotech/fermions-workflow/pkg/worker/environment"
)

var ErrFunctionNotImplemented = errors.New("function not implemented")

func New(f model.Function) (fPrepared environment.Function, err error) {
	switch f.Type {
	case model.FunctionTypeREST:
		fPrepared = newRest(f.Operation)
	case model.FunctionTypeExpression:
		fPrepared = newExpression(f.Operation)
	default:
		return nil, ErrFunctionNotImplemented
	}

	err = fPrepared.Init()
	if err != nil {
		return nil, err
	}

	return fPrepared, err
}
