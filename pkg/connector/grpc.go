package connector

import (
	"fmt"
	"log"

	"cuelang.org/go/cue/format"
	"cuelang.org/go/encoding/protobuf"
)

func Load() {
	file, err := protobuf.Extract("testdata/basic.proto", nil, &protobuf.Config{
		Paths: []string{ /* paths to proto includes */ },
	})

	if err != nil {
		log.Fatal(err, "")
	}

	b, _ := format.Node(file)
	fmt.Println(string(b))
	// ioutil.WriteFile("out.cue", b, 0644)
}
