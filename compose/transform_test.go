package compose

import "testing"

func TestTransformHelloWorld(t *testing.T) {
	trf := Transformer{
		ComposeFiles: []string{"../examples/helloworld/docker-compose.yml"},
		ProjectName:  "helloworld",
	}

	_, err := trf.Transform()
	if err != nil {
		t.Fatal(err)
	}
}

func TestTransformComplex(t *testing.T) {
	trf := Transformer{
		ComposeFiles: []string{"../examples/complex/docker-compose.yml"},
		ProjectName:  "complex",
		EnvironmentLookup: envMap{
			"ECR_REPOSITORY": []string{"http://example.org"},
			"DATABASE_URL":   []string{"http://database.org"},
		},
	}

	_, err := trf.Transform()
	if err != nil {
		t.Fatal(err)
	}
}
