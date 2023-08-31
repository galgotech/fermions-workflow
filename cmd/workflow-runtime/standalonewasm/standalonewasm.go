package main

import (
	_ "crypto/sha512"
	_ "os"
	"syscall/js"

	"github.com/galgotech/fermions-workflow/pkg/standalonewasm"
)

var fermionsStandalone *standalonewasm.StandaloneWasm

func main() {
	var err error
	setting := standalonewasm.NewSetting()
	fermionsStandalone, err = standalonewasm.Initialize(setting)
	if err != nil {
		panic(err)
	}

	done := make(chan struct{}, 0)
	js.Global().Set("fermionsExecWorkflow", js.FuncOf(fermionsExecWorkflow))
	js.Global().Set("fermionsKillWorkflow", js.FuncOf(fermionsKillWorkflow))
	<-done
}

func fermionsExecWorkflow(this js.Value, args []js.Value) interface{} {
	err := fermionsStandalone.Execute([]byte(args[0].String()))
	if err != nil {
		return err.Error()
	}
	return ""
}

func fermionsKillWorkflow(this js.Value, args []js.Value) interface{} {
	err := fermionsStandalone.Kill(args[0].String())
	if err != nil {
		return err.Error()
	}
	return ""
}
