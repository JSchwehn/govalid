package govalid_test

import (
	"encoding/json"
	"fmt"
	v "github.com/gima/govalid/v1"
	"reflect"
	"testing"
)

func TestExample(t *testing.T) {

	// set up raw json data
	rawJson := []byte(`
		{
    	"status": true,
      "data": {
      	"token": "CAFED00D",
	      "debug": 69306,
      	"items": [
      	  { "url": "https://news.ycombinator.com/", "comment": "clickbaits" },
          { "url": "http://golang.org/", "comment": "some darn gophers" },
          { "url": "http://www.kickstarter.com/", "comment": "opensource projects. yeah.." }
       	],
       	"ghost2": null,
       	"meta": {
       	 "g": 1,
         "xyzzy": 0.25,
         "blöö": 0.5
       	}
      }
		}`)

	// decode json
	var data interface{}
	if err := json.Unmarshal(rawJson, &data); err != nil {
		t.Fatal("JSON parsing failed. Err =", err)
	}

	// set up a custom validator function
	myValidatorFunc := func(data interface{}) (path string, err error) {
		path = "myValidatorFunc"

		validate, ok := data.(string)
		if !ok {
			return path, fmt.Errorf("expected string, got %v", reflect.TypeOf(data))
		}

		if validate != "CAFED00D" {
			return path, fmt.Errorf("expected CAFED00D, got %s", validate)
		}

		return "", nil
	}

	// construct the schema which is used to validate data
	schema := v.Object(
		v.ObjKV("status", v.Boolean()),
		v.ObjKV("data", v.Object(
			v.ObjKV("token", v.Function(myValidatorFunc)),
			v.ObjKV("debug", v.Number(v.NumMin(1), v.NumMax(99999))),
			v.ObjKV("items", v.Array(v.ArrEach(v.Object(
				v.ObjKV("url", v.String(v.StrMin(1))),
				v.ObjKV("comment", v.Optional(v.String())),
			)))),
			v.ObjKV("ghost", v.Optional(v.String())),
			v.ObjKV("ghost2", v.Optional(v.String())),
			v.ObjKV("meta", v.Object(
				v.ObjKeys(v.String()),
				v.ObjValues(v.Or(v.Number(v.NumMin(.01), v.NumMax(1.1)), v.String())),
			)),
		)),
	)

	// do the actual validation
	if path, err := schema.Validate(data); err == nil {
		t.Log("Validation passed.")
	} else {
		t.Fatalf("Validation failed at %s. Error (%s)", path, err)
	}
}
