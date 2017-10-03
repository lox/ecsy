package taskdef

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"reflect"

	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/ghodss/yaml"
)

func ParseFile(file string, env []string) (*ecs.RegisterTaskDefinitionInput, error) {
	body, err := ioutil.ReadFile(file)
	if err != nil {
		return result, err
	}

	p := parser{
		Env:  env,
		Body: body,
	}

	return p.parse()
}

type parser struct {
	Env  []string
	Body []byte
}

func (p parser) parse() (*ecs.RegisterTaskDefinitionInput, error) {
	// Unmarshal the pipeline into an actual data structure
	unmarshaled, err := unmarshal(p.Body)
	if err != nil {
		return nil, err
	}

	// Recursively go through the entire pipeline and perform environment
	// variable interpolation on strings
	interpolated, err := p.interpolate(unmarshaled)
	if err != nil {
		return nil, err
	}

	// Return to json
	jsonBytes, err := json.Marshal(interpolated)
	if err != nil {
		return nil, err
	}

	var result ecs.RegisterTaskDefinitionInput

	// And then into the task definition 👌🏻 🤞🏻
	if err = json.Unmarshal(jsonBytes, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func unmarshal(body []byte) (interface{}, error) {
	var unmarshaled interface{}

	err := yaml.Unmarshal(body, &unmarshaled)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse: %v", err)
	}

	return unmarshaled, nil
}

func (p parser) interpolate(obj interface{}) (interface{}, error) {
	// Make sure there's something actually to interpolate
	if obj == nil {
		return nil, nil
	}

	// Wrap the original in a reflect.Value
	original := reflect.ValueOf(obj)

	// Make a copy that we'll add the new values to
	copy := reflect.New(original.Type()).Elem()

	err := p.interpolateRecursive(copy, original)
	if err != nil {
		return nil, err
	}

	// Remove the reflection wrapper
	return copy.Interface(), nil
}

func (p parser) interpolateRecursive(copy, original reflect.Value) error {
	switch original.Kind() {
	// If it is a pointer we need to unwrap and call once again
	case reflect.Ptr:
		// To get the actual value of the original we have to call Elem()
		// At the same time this unwraps the pointer so we don't end up in
		// an infinite recursion
		originalValue := original.Elem()

		// Check if the pointer is nil
		if !originalValue.IsValid() {
			return nil
		}

		// Allocate a new object and set the pointer to it
		copy.Set(reflect.New(originalValue.Type()))

		// Unwrap the newly created pointer
		err := p.interpolateRecursive(copy.Elem(), originalValue)
		if err != nil {
			return err
		}

	// If it is an interface (which is very similar to a pointer), do basically the
	// same as for the pointer. Though a pointer is not the same as an interface so
	// note that we have to call Elem() after creating a new object because otherwise
	// we would end up with an actual pointer
	case reflect.Interface:
		// Get rid of the wrapping interface
		originalValue := original.Elem()

		// Check to make sure the interface isn't nil
		if !originalValue.IsValid() {
			return nil
		}

		// Create a new object. Now new gives us a pointer, but we want the value it
		// points to, so we have to call Elem() to unwrap it
		copyValue := reflect.New(originalValue.Type()).Elem()

		err := p.interpolateRecursive(copyValue, originalValue)
		if err != nil {
			return err
		}

		copy.Set(copyValue)

	// If it is a struct we interpolate each field
	case reflect.Struct:
		for i := 0; i < original.NumField(); i += 1 {
			err := p.interpolateRecursive(copy.Field(i), original.Field(i))
			if err != nil {
				return err
			}
		}

	// If it is a slice we create a new slice and interpolate each element
	case reflect.Slice:
		copy.Set(reflect.MakeSlice(original.Type(), original.Len(), original.Cap()))

		for i := 0; i < original.Len(); i += 1 {
			err := p.interpolateRecursive(copy.Index(i), original.Index(i))
			if err != nil {
				return err
			}
		}

	// If it is a map we create a new map and interpolate each value
	case reflect.Map:
		copy.Set(reflect.MakeMap(original.Type()))

		for _, key := range original.MapKeys() {
			originalValue := original.MapIndex(key)

			// New gives us a pointer, but again we want the value
			copyValue := reflect.New(originalValue.Type()).Elem()
			err := p.interpolateRecursive(copyValue, originalValue)
			if err != nil {
				return err
			}

			copy.SetMapIndex(key, copyValue)
		}

	// If it is a string interpolate it (yay finally we're doing what we came for)
	case reflect.String:
		interpolated, err := Interpolate(p.Env, original.Interface().(string))
		if err != nil {
			return err
		}
		copy.SetString(interpolated)

	// And everything else will simly be taken from the original
	default:
		copy.Set(original)
	}

	return nil
}
