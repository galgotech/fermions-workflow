package data

import (
	"encoding/json"
	"fmt"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/imdario/mergo"
	"github.com/serverlessworkflow/sdk-go/v2/model"
)

var ObjectNil = model.Object{}

func FromEvent(event cloudevents.Event) (model.Object, error) {
	var dataOut any
	err := json.Unmarshal(event.Data(), &dataOut)
	if err != nil {
		return ObjectNil, err
	}
	return model.FromInterface(dataOut), nil
}

func Merge(dst, src model.Object) (model.Object, error) {
	if src.Type == model.Null {
		return dst, nil
	}

	if dst.Type != model.Map {
		return ObjectNil, fmt.Errorf("dst data.Merge need be a map. Current %d", dst.Type)
	}
	if src.Type != model.Map {
		return ObjectNil, fmt.Errorf("src data.Merge need be a map. Current %d", src.Type)
	}

	err := mergo.Merge(&dst.MapValue, src.MapValue, mergo.WithOverride)
	if err != nil {
		return ObjectNil, err
	}
	return dst, nil
}

func Unmarshal(data []byte) (model.Object, error) {
	object := model.Object{}
	err := object.UnmarshalJSON(data)
	if err != nil {
		return ObjectNil, err
	}
	return object, nil
}

func Marshal(object model.Object) ([]byte, error) {
	return object.MarshalJSON()
}
