package main

import (
	"testing"
)

func TestVerifyData(t *testing.T) {

	t.Run("data is equal, medical plan", func(t *testing.T) {
		got, err := VerifyData("Ipe", "Ipe", MedicalPlan)
		want := Success

		if err != nil {
			t.Fatalf("Got error and error should have been nil")
		}

		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	t.Run("data is not equal, medical_pan", func(t *testing.T) {
		_, err := VerifyData("Ipe", "Unimed", MedicalPlan)
		want := ErrMedicalPlanNotMatch

		if err == nil {
			t.Fatalf("Expected an error, got nil")
		}
		if err != want {
			t.Errorf("got %q want %q", err, want)
		}
	})

	t.Run("data type does not exist", func(t *testing.T) {
		_, err := VerifyData("Ipe", "Unimed", "Something")
		want := ErrNotFound

		if err == nil {
			t.Fatalf("Expected an error, got nil")
		}
		if err != want {
			t.Errorf("got %q want %q", err, want)
		}
	})
}

// func assertVerifyData(t *testing.T, got string, want string) {
// 	t.Helper()
// 	if got != want {
// 		t.Errorf("got %q, want %q", got, want)
// 	}
// }
