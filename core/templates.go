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

func (t Templates) EcsCluster() Template {
	b, err := t.read("/templates/build/ecs-cluster.json")
	if err != nil {
		panic(err)
	}

	return Template(string(b))
}

type Template string
