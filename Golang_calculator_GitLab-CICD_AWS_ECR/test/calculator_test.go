package calculator_test

import (
	"calculator/internal/calculator"
	"math"
	"testing"
)

func TestAdd(t *testing.T) {
	t.Parallel()

	type testCase struct {
		a, b float64
		want float64
	}

	testCases := []testCase{
		{a: 2, b: 3, want: 5},
		{a: 3, b: 4.4, want: 7.4},
		{a: 0, b: 2, want: 2},
	}

	for _, tc := range testCases {
		got := calculator.Add(tc.a, tc.b)
		if tc.want != got {
			t.Errorf("Add(%f, %f) want %f, got %f", tc.a, tc.b, tc.want, got)
		}
	}
}

func TestSubtract(t *testing.T) {
	t.Parallel()
	type testCase struct {
		a, b float64
		want float64
	}

	testCases := []testCase{
		{a: 5.5, b: 3, want: 2.5},
		{a: 4.4, b: 2.2, want: 2.2},
		{a: 0, b: 5, want: -5},
	}

	for _, tc := range testCases {
		got := calculator.Subtract(tc.a, tc.b)
		if tc.want != got {
			t.Errorf("Subtract(%f, %f) want %f, got %f", tc.a, tc.b, tc.want, got)
		}
	}
}

func TestMultiply(t *testing.T) {
	t.Parallel()

	type testCase struct {
		a, b float64
		want float64
	}

	testCases := []testCase{
		{a: 4, b: 3, want: 12},
		{a: 3.3, b: 0, want: 0},
		{a: 1, b: 1.1, want: 1.1},
	}

	for _, tc := range testCases {
		got := calculator.Multiply(tc.a, tc.b)
		if tc.want != got {
			t.Errorf("Mulitply(%f, %f) want %f, got %f", tc.a, tc.b, tc.want, got)
		}
	}
}

func closeEnough(a, b, tolerance float64) bool {
	return math.Abs(a-b) <= tolerance
}

func TestDivide(t *testing.T) {
	t.Parallel()

	type testCase struct {
		a, b float64
		want float64
	}

	testCases := []testCase{
		{a: 1, b: 2, want: 0.5},
		{a: 2, b: 1, want: 2},
		{a: 0, b: 22, want: 0},
		{a: 1, b: 3, want: 0.3333},
	}

	for _, tc := range testCases {
		got, err := calculator.Divide(tc.a, tc.b)
		if err != nil {
			t.Fatalf("want no error, invalid input %f", err)
		}
		if !closeEnough(tc.want, got, 0.001) {
			t.Errorf("Divide (%f, %f) want %f, got %f", tc.a, tc.b, tc.want, got)
		}
	}
}

func TestSqrt(t *testing.T) {
	t.Parallel()

	type testCase struct {
		a    float64
		want float64
	}

	testCases := []testCase{
		{a: 25, want: 5},
		{a: 1, want: 1},
		{a: 0, want: 0},
	}

	for _, tc := range testCases {
		got, err := calculator.Sqrt(tc.a)
		if err != nil {
			t.Fatalf("Fatal %q", err)
		}
		if !closeEnough(tc.want, got, 0.000000000001) {
			t.Errorf("Sqrt(%f) want %f got %f", tc.a, tc.want, got)
		}
	}
}
