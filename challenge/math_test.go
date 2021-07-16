package challenge_test

import (
	"testing"

	"github.com/GPKyte/banter/challenge"
)

func TestExampleFactorial(t *testing.T) {
	t.Logf("Factorial of 5 is %d?", challenge.ExampleFactorial())
}
