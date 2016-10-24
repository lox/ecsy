package compose

import (
	"log"
	"testing"
)

func TestTransformV1(t *testing.T) {
	input, err := TransformComposeFile("../examples/helloworld/docker-compose.yml", "helloworld")
	if err != nil {
		t.Fatal(err)
	}

	log.Printf("%#v", input)
}
