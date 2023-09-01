package data

import (
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/imdario/mergo"
	"github.com/serverlessworkflow/sdk-go/v2/model"
)

type Data[T any] map[string]T

var ObjectNil = model.Object{}

func FromEvent(event cloudevents.Event) (model.Object, error) {
	var dataOut any
	err := event.DataAs(&dataOut)
	if err != nil {
		return ObjectNil, err
	}
	object := model.FromInterface(dataOut)
	return object, nil
}

func FromInterface(data any) model.Object {
	object := model.FromInterface(data)
	return object
}

func ToInterface(data model.Object) any {
	return model.ToInterface(data)
}

func Merge(dst, src model.Object) error {
	srcAny := model.ToInterface(src)
	dstAny := model.ToInterface(dst)
	if err := mergo.Merge(dstAny, srcAny, mergo.WithOverride); err != nil {
		return err
	}
	return nil
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
