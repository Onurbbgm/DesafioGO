package main

import (
	"reflect"
	"testing"
)

func TestCheckCSV(t *testing.T) {
	t.Run("file creation success", func(t *testing.T) {
		err := CheckCSV("testOne.csv", "testTwo.csv")

		if err != nil {
			t.Fatalf("Got an error %q and error should have been nil", err)
		}
	})
	t.Run("file wrong number of lines error", func(t *testing.T) {
		err := CheckCSV("testNumLinesOne.csv", "testNumLinesTwo.csv")

		if err == nil {
			t.Fatalf("Expected an error, got nil")
		}

		if err != ErrNumberOfLines {
			t.Errorf("got %q, want %q", err, ErrNumberOfLines)
		}
	})
	t.Run("file wrong number of lines error", func(t *testing.T) {
		err := CheckCSV("testNotExisting.csv", "testNumLinesTwo.csv")

		if err == nil {
			t.Fatalf("Expected an error, got nil")
		}
	})
	t.Run("total errors match", func(t *testing.T) {
		ResetTotals()
		err := CheckCSV("testOne.csv", "testTwo.csv")
		_, got := GenerateRowTotals()
		want := []string{"9/10", "7/10", "0/10", "7/10", "0/10", "7/10", "6/10", "10/10", "10/10"}
		if err != nil {
			t.Fatalf("Got an error %q and error should have been nil", err)
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v want %v", got, want)
		}
	})
}

func TestVerifyData(t *testing.T) {
	verifyDataErrorsTests := []struct {
		name     string
		valueOne string
		valueTwo string
		dataType string
	}{
		{"data type does not exist", "Ipe", "Ipe", "Something"},
		{"error in medical_plan", "Ipe", "Unimed", MedicalPlan},
		{"error in dental_plan", "Ipe", "Unimed", DentalPlan},
		{"error in employee_name", "Pedro", "Carlos", EmployeeName},
		{"error in language", "en-CA", "fr-CA", Language},
		{"error in claimant_name", "Carlos", "Pedro", ClaimantName},
		{"error in relationship_type", "Son", "Daughter", RelationshipType},
		{"error in gender", "Female", "Male", Gender},
		{"error in effective_date", "6/20/2019", "8/27/2019", EffectiveDate},
		{"error in termination_date", "9/20/2019", "11/27/2019", EffectiveDate},
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
			t.Fatalf("Got an error and error should have been nil")
		}

		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})
}
