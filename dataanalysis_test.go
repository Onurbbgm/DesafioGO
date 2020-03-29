package main

import (
	"testing"
)

func TestVerifyData(t *testing.T) {
	verifyDataErrorsTests := []struct {
		name     string
		valueOne string
		valueTwo string
		dataType string
	}{
		{"data type does not exist", "Ipe", "Ipe", "Something"},
		{"error in medical_plan", "Ipe", "Unimed", MedicalPlan},
	}
	for _, tt := range verifyDataErrorsTests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := VerifyData(tt.valueOne, tt.valueTwo, tt.dataType)
			want, found := Errors[tt.dataType]

			if !found {
				if err == ErrNotFound {
					return
				}
				if err != nil && err != ErrNotFound {
					t.Fatalf("Expected %q and got %q", ErrNotFound, err)
				}
			}

			if err == nil {
				t.Fatalf("Expected an error, got nil")
			}

			if err != want {
				t.Errorf("got %q, want %q", err, want)
			}
		})
	}
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
}
