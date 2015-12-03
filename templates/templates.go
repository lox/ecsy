package templates

import (
	"io/ioutil"
	"net/http"
)

var FileSystem http.FileSystem

func init() {
	FileSystem = FS(false)
}

func readTemplateBytes(path string) ([]byte, error) {
	tpl, err := FileSystem.Open(path)
	if err != nil {
		return nil, err
	}

	defer tpl.Close()
	tplB, err := ioutil.ReadAll(tpl)
	if err != nil {
		return nil, err
	}

	return tplB, nil
}

type Template string

func EcsStack() string {
	b, err := readTemplateBytes("/templates/build/ecs-stack.json")
	if err != nil {
		panic(err)
	}
	return string(b)
}

func EcsService() string {
	b, err := readTemplateBytes("/templates/build/ecs-service.json")
	if err != nil {
		panic(err)
	}
	return string(b)
}

func NetworkStack() string {
	b, err := readTemplateBytes("/templates/build/network-stack.json")
	if err != nil {
		panic(err)
	}
	return string(b)
}
