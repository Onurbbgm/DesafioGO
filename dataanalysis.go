package main

type DictionaryErr string

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
	ErrNotFound                 = DictionaryErr("Data type not found")
	ErrMedicalPlanNotMatch      = DictionaryErr("Medical plan does not match in both datasets")
	ErrDentalPlanNotMatch       = DictionaryErr("Dental plan does not match in both datasets")
	ErrEmployeeNameNotMatch     = DictionaryErr("Employee name does not match in both datasets")
	ErrLanguageNotMatch         = DictionaryErr("Language does not match in both datasets")
	ErrClaimantNameNotMatch     = DictionaryErr("Claimant name does not match in both datasets")
	ErrRelationshipTypeNotMatch = DictionaryErr("Relationship type does not match in both datasets")
	ErrGenderTypeNotMatch       = DictionaryErr("Gender type does not match in both datasets")
	ErrEffectiveDateNotMatch    = DictionaryErr("Effective date does not match in both datasets")
	ErrTerminationDateNotMatch  = DictionaryErr("Termination date does not match in both datasets")
)

func VerifyData(valueOne, valoueTwo string, dataType string) (string, error) {
	if valueOne == valoueTwo {
		return Success, nil
	}
	return "", getErrorMessage(dataType)
}

func getErrorMessage(dataType string) error {
	switch dataType {
	case MedicalPlan:
		return ErrMedicalPlanNotMatch
	case DentalPlan:
		return ErrDentalPlanNotMatch
	case EmployeeName:
		return ErrEmployeeNameNotMatch
	case Language:
		return ErrLanguageNotMatch
	case ClaimantName:
		return ErrClaimantNameNotMatch
	case RelationshipType:
		return ErrRelationshipTypeNotMatch
	case Gender:
		return ErrGenderTypeNotMatch
	case EffectiveDate:
		return ErrEffectiveDateNotMatch
	case TerminationDate:
		return ErrTerminationDateNotMatch
	default:
		return ErrNotFound
	}
}

func main() {

}

func (e DictionaryErr) Error() string {
	return string(e)
}
