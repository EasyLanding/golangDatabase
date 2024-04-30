package main

import (
	"bytes"
	"fmt"
	"testing"
)

func TestHelloWorldFunc(t *testing.T) {
	var buf bytes.Buffer

	fmt.Fprint(&buf, HelloWorld())

	output := buf.String()
	expected := "Hello world!"

	if output != expected {
		t.Errorf("Unexpected output: %s", output)
	}
}
