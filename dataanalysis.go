package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
)

type Err string

//Columns names e success message
const (
	Success          = "Ok"
	MedicalPlan      = "medical_plan"
	DentalPlan       = "dental_plan"
	EmployeeName     = "employee_name"
	Language         = "language"
	ClaimantName     = "claimant_name"
	RelationshipType = "relationship_type"
	Gender           = "gender"
	EffectiveDate    = "effective_date"
	TerminationDate  = "termination_date"
)

//Consts of all the error messages of the system
const (
	ErrNotFound                 = Err("Data type not found")
	ErrColumNumber              = Err("Invalid column number")
	ErrMedicalPlanNotMatch      = Err("Medical plan does not match in both datasets")
	ErrDentalPlanNotMatch       = Err("Dental plan does not match in both datasets")
	ErrEmployeeNameNotMatch     = Err("Employee name does not match in both datasets")
	ErrLanguageNotMatch         = Err("Language does not match in both datasets")
	ErrClaimantNameNotMatch     = Err("Claimant name does not match in both datasets")
	ErrRelationshipTypeNotMatch = Err("Relationship type does not match in both datasets")
	ErrGenderTypeNotMatch       = Err("Gender type does not match in both datasets")
	ErrEffectiveDateNotMatch    = Err("Effective date does not match in both datasets")
	ErrTerminationDateNotMatch  = Err("Termination date does not match in both datasets")
)

//Errors a dictionary for the error for the appropriated column
var Errors = map[string]Err{
	MedicalPlan:      ErrMedicalPlanNotMatch,
	DentalPlan:       ErrDentalPlanNotMatch,
	EmployeeName:     ErrEmployeeNameNotMatch,
	Language:         ErrLanguageNotMatch,
	ClaimantName:     ErrClaimantNameNotMatch,
	RelationshipType: ErrRelationshipTypeNotMatch,
	Gender:           ErrGenderTypeNotMatch,
	EffectiveDate:    ErrEffectiveDateNotMatch,
	TerminationDate:  ErrTerminationDateNotMatch,
}

func CheckCSV() error {
	csvOne, err := os.Open("ourData.csv")
	if err != nil {
		return err
	}

	csvTwo, err := os.Open("clientData.csv")
	if err != nil {
		return err
	}

	resultCSV, err := os.Create("result.csv")
	if err != nil {
		return err
	}

	readerOne := csv.NewReader(bufio.NewReader(csvOne))
	readerTwo := csv.NewReader(bufio.NewReader(csvTwo))

	defer csvOne.Close()
	defer csvTwo.Close()
	defer resultCSV.Close()

	writer := csv.NewWriter(resultCSV)
	defer writer.Flush()

	//Jump first line of both files
	if _, err := readerOne.Read(); err != nil {
		return err
	}

	if _, err := readerTwo.Read(); err != nil {
		return err
	}

	for {
		lineOne, err := readerOne.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		lineTwo, err := readerTwo.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		var row []string
		for i := 1; i <= 9; i++ {
			result, errData := VerifyData(lineOne[i], lineTwo[i], getColumn(i))
			if errData == nil {
				row = append(row, result)
			} else {
				row = append(row, errData.Error())
			}
		}
		writer.Write(row)
	}

	return nil
}

//VerifyData gets twi values and checks if they are the same or gives an error of the appropriated type
func VerifyData(valueOne, valoueTwo string, dataType string) (string, error) {
	_, found := Errors[dataType]
	if !found {
		return "", ErrNotFound
	}
	if valueOne == valoueTwo {
		return Success, nil
	}
	return "", getErrorMessage(dataType)
}

func getErrorMessage(dataType string) error {
	return Errors[dataType]
}

func getColumn(num int) string {
	switch num {
	case 1:
		return MedicalPlan
	case 2:
		return DentalPlan
	case 3:
		return EmployeeName
	case 4:
		return Language
	case 5:
		return ClaimantName
	case 6:
		return RelationshipType
	case 7:
		return Gender
	case 8:
		return EffectiveDate
	case 9:
		return TerminationDate
	default:
		log.Fatal(ErrColumNumber)
		return ErrColumNumber.Error()
	}
}

func main() {
	err := CheckCSV()
	if err != nil {
		fmt.Println(err)
	}
}

func (e Err) Error() string {
	return string(e)
}
