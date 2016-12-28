package testutils

import (
	"reflect"
	"testing"
)

var errorTemplate = "Got %q, expected %q"

func AssertDeepEqual(t *testing.T, got interface{}, expected interface{}) {
	if !reflect.DeepEqual(got, expected) {
		TError(t, got, expected)
	}
}

func TError(t *testing.T, got interface{}, expected interface{}) {
	t.Error("Got:", got, "Expected:", expected)
}
