package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"
)

func TestCheckCSV(t *testing.T) {

	t.Run("file creation success", func(t *testing.T) {
		fileOne, err := os.Open("testOne.csv")
		fileTwo, err := os.Open("testTwo.csv")
		if err != nil {
			t.Fatalf("Got an error %q and error should have been nil", err)
		}
		response := httptest.NewRecorder()

		CheckCSV(fileOne, fileTwo, response)

		assertStatus(t, response.Code, http.StatusOK)
		assertResponseBody(t, response.Body.String(), "File result created with report.")
	})
	t.Run("file wrong number of lines error", func(t *testing.T) {
		fileOne, err := os.Open("testNumLinesOne.csv")
		fileTwo, err := os.Open("testNumLinesTwo.csv")

		if err != nil {
			t.Fatalf("Got an error %q and error should have been nil", err)
		}

		response := httptest.NewRecorder()
		CheckCSV(fileOne, fileTwo, response)

		assertStatus(t, response.Code, http.StatusBadRequest)
		assertResponseBody(t, response.Body.String(), "Got error "+ErrNumberOfLines.Error())
	})
	t.Run("total errors match", func(t *testing.T) {
		fileOne, err := os.Open("testOne.csv")
		fileTwo, err := os.Open("testTwo.csv")
		if err != nil {
			t.Fatalf("Got an error %q and error should have been nil", err)
		}
		response := httptest.NewRecorder()

		CheckCSV(fileOne, fileTwo, response)
		_, got := GenerateRowTotals()
		want := []string{"", "9/10", "7/10", "7/10", "7/10", "6/10", "10/10", "10/10"}

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

func TestServer(t *testing.T) {
	t.Run("send wrong number of file(s)", func(t *testing.T) {
		numer := 1
		response := httptest.NewRecorder()
		CheckNumberOfFiles(numer, response)

		assertStatus(t, response.Code, http.StatusBadRequest)
		concatenated := fmt.Sprint("Needs to get 2 files got ", numer)
		assertResponseBody(t, response.Body.String(), concatenated)

	})
	t.Run("send file of wrong type", func(t *testing.T) {
		response := httptest.NewRecorder()
		extension := ".txt"
		CheckFileExtension(extension, response)
		assertStatus(t, response.Code, http.StatusBadRequest)
		assertResponseBody(t, response.Body.String(), "Expected contentType text/csv, got "+extension)
	})
	t.Run("send GET to server", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/checkFiles/", nil)
		response := httptest.NewRecorder()
		server := &DataAnalysisServer{}
		server.ServeHTTP(response, req)
		assertStatus(t, response.Code, http.StatusMethodNotAllowed)
	})
}

func assertResponseBody(t *testing.T, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("response body is wrong, got %q want %q", got, want)
	}
}

func assertStatus(t *testing.T, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("dig not get correct status, got %d, want %d", got, want)
	}
}

func BenchmarkCheckCVS(b *testing.B) {
	fileOne, _ := os.Open("testOne.csv")
	fileTwo, _ := os.Open("testTwo.csv")

	response := httptest.NewRecorder()

	CheckCSV(fileOne, fileTwo, response)
}
