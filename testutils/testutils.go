package testutils

import (
	"encoding/json"
	"fmt"
	"reflect"
)

func JsonEq(foundBody, expectedBody string) (bool, error) {
	var o1 interface{}
	var o2 interface{}

	var err error
	err = json.Unmarshal([]byte(foundBody), &o1)
	if err != nil {
		return false, fmt.Errorf("json syntax error in found body: %s", err)
	}
	err = json.Unmarshal([]byte(expectedBody), &o2)
	if err != nil {
		return false, fmt.Errorf("json syntax error in expected body: %s", err)
	}

	return reflect.DeepEqual(o1, o2), nil
}
