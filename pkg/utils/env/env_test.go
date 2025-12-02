package env

import (
	"os"
	"testing"
)

func TestExpandEnvWithDefault_EnvPresent(t *testing.T) {
	os.Setenv("FOO", "bar")
	defer os.Unsetenv("FOO")

	input := "Value: ${FOO:default}"
	expected := "Value: bar"

	result := ExpandEnvWithDefault(input)
	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}

func TestExpandEnvWithDefault_EnvAbsent_UsesDefault(t *testing.T) {
	os.Unsetenv("FOO")

	input := "Value: ${FOO:default}"
	expected := "Value: default"

	result := ExpandEnvWithDefault(input)
	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}

func TestExpandEnvWithDefault_EnvAbsent_NoDefault(t *testing.T) {
	os.Unsetenv("FOO")

	input := "Value: ${FOO}"
	expected := "Value: " // default is empty string

	result := ExpandEnvWithDefault(input)
	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}

func TestExpandEnvWithDefault_MultipleVariables(t *testing.T) {
	os.Setenv("A", "1")
	os.Unsetenv("B")
	os.Setenv("C", "3")
	defer func() {
		os.Unsetenv("A")
		os.Unsetenv("C")
	}()

	input := "${A} + ${B:2} + ${C} = 1 + 2 + 3"
	expected := "1 + 2 + 3 = 1 + 2 + 3"

	result := ExpandEnvWithDefault(input)
	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}

func TestExpandEnvWithDefault_NoVariables(t *testing.T) {
	input := "Plain string"
	expected := "Plain string"

	result := ExpandEnvWithDefault(input)
	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}
