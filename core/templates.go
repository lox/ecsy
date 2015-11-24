package core

import (
	"io/ioutil"
	"net/http"
)

type Templates struct {
	http.FileSystem
}

func (t Templates) read(path string) ([]byte, error) {
	tpl, err := t.FileSystem.Open(path)
	if err != nil {
		return nil, err
	}

	defer tpl.Close()
	tplB, err := ioutil.ReadAll(tpl)
	if err != nil {
		return nil, err
	}

	return tplB, err
}

func (t Templates) EcsStack() Template {
	b, err := t.read("/templates/build/ecs-stack.json")
	if err != nil {
		panic(err)
	}

	return Template(string(b))
}

func (t Templates) EcsService() Template {
	b, err := t.read("/templates/build/ecs-service.json")
	if err != nil {
		panic(err)
	}

	return Template(string(b))
}

type Template string
