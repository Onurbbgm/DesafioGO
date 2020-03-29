package main

type Err string

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

const (
	ErrNotFound                 = Err("Data type not found")
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

func main() {

}

func (e Err) Error() string {
	return string(e)
}
