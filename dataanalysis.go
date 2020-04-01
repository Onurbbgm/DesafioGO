package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"

	"github.com/gabriel-vasile/mimetype"
)

type Err string

type resultDone struct {
	bool
}

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
	TotalErrors      = "Total errors of "
)

//Consts of all the error messages of the system
const (
	ErrNotFound                 = Err("Data type not found")
	ErrColumNumber              = Err("Invalid column number")
	ErrNumberOfLines            = Err("Files have different number of lines, so can not be compared")
	ErrReadinLines              = Err("Erro when trying to read line")
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

var (
	totalErrMP = 0
	totalErrDP = 0
	//totalErrEN = 0
	totalErrL = 0
	//totalErrCN = 0
	totalErrRT = 0
	totalErrG  = 0
	totalErrED = 0
	totalErrTD = 0
	totalData  = 0
)

type DataAnalysisServer struct {
}

func (d *DataAnalysisServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodPost:
		getFiles(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func getFiles(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 << 20)

	files := r.MultipartForm

	if len(files.File) != 2 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Needs to get 2 files got ", len(files.File))
		//log.Fatalf("Needs to get 2 files got %d", len(files.File))
		return
	}
	for _, v := range files.File {
		for _, f := range v {
			file, err := f.Open()
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprint(w, "Got error ", err.Error())
				return
			}
			defer file.Close()

			mime, err := mimetype.DetectReader(file)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprint(w, "Got error ", err.Error())
				return
			}
			fmt.Println(mime.Extension())
			if mime.Extension() != ".csv" {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Println(mime)
				fmt.Fprint(w, "Expected contentType text/csv, got ", mime)
				return
			}
		}
	}
}

func getFileContentType(out *os.File) (string, error) {

	// Only the first 512 bytes are used to sniff the content type.
	buffer := make([]byte, 512)

	_, err := out.Read(buffer)
	if err != nil {
		return "", err
	}

	// Use the net/http package's handy DectectContentType function. Always returns a valid
	// content-type by returning "application/octet-stream" if no others seemed to match.
	contentType := http.DetectContentType(buffer)

	return contentType, nil
}

var wg sync.WaitGroup
var mu = sync.Mutex{}

func CheckCSV(fileOne, fileTwo string) error {
	csvOne, errOpen := os.Open(fileOne)
	if errOpen != nil {
		return errOpen
	}

	csvTwo, errOpen := os.Open(fileTwo)
	if errOpen != nil {
		return errOpen
	}

	resultCSV, errCreate := os.Create("result.csv")
	if errCreate != nil {
		return errCreate
	}

	readerOne := csv.NewReader(bufio.NewReader(csvOne))
	readerTwo := csv.NewReader(bufio.NewReader(csvTwo))

	numLinesOne, err := lineCounter(readerOne)
	if err != nil {
		return err
	}
	numLinesTwo, err := lineCounter(readerTwo)
	if err != nil {
		return err
	}

	totalData = numLinesOne - 2

	csvOne.Seek(0, io.SeekStart)
	csvTwo.Seek(0, io.SeekStart)

	if numLinesOne != numLinesTwo {
		return ErrNumberOfLines
	}

	defer csvOne.Close()
	defer csvTwo.Close()
	defer resultCSV.Close()

	writer := csv.NewWriter(resultCSV)
	defer writer.Flush()

	firstRow := []string{"", MedicalPlan, DentalPlan, Language, RelationshipType, Gender, EffectiveDate, TerminationDate}
	writer.Write(firstRow)

	//Jump first line of both files
	if _, err := readerOne.Read(); err != nil {
		return err
	}

	if _, err := readerTwo.Read(); err != nil {
		return err
	}

	wg.Add(numLinesOne)

	//count := 0
	done := false
	for {

		select {
		case resultRoutine := <-readAndWriteLines(readerOne, readerTwo, writer):
			if resultRoutine.bool == true {
				done = true
			}
		}

		if done {
			break
		}
	}

	totalsNamesRow, totalsRow := GenerateRowTotals()
	writer.Write([]string{"", "", "", "", "", "", "", "", ""})
	writer.Write(totalsNamesRow)
	writer.Write(totalsRow)

	return nil
}

func readAndWriteLines(readerOne, readerTwo *csv.Reader, writer *csv.Writer) chan resultDone {
	ch := make(chan resultDone)

	// ch := make(chan struct{})
	go func(w *sync.WaitGroup) {
		mu.Lock()
		defer mu.Unlock()

		lineOne, err := readerOne.Read()
		if err == io.EOF {
			ch <- resultDone{true}
			close(ch)
			return
		} else if err != nil {
			ch <- resultDone{false}
			close(ch)
			log.Fatalf(err.Error())
			return
		}
		lineTwo, err := readerTwo.Read()
		if err == io.EOF {
			ch <- resultDone{true}
			close(ch)
			return
		} else if err != nil {
			ch <- resultDone{false}
			close(ch)
			log.Fatalf(err.Error())
			return
		}
		var row []string
		row = append(row, "Employee "+lineOne[3]+" with relationship to claimant "+lineOne[5]+" has the following problems")
		for i := 1; i <= 9; i++ {
			columName := getColumn(i)
			if columName == EmployeeName || columName == ClaimantName {
				continue
			}
			result, errData := VerifyData(lineOne[i], lineTwo[i], getColumn(i))
			if errData == nil {
				row = append(row, result)
			} else {
				addTotal(i)
				row = append(row, errData.Error()+": got "+lineOne[i]+" and "+lineTwo[i])
			}
		}
		writer.Write(row)
		w.Done()
		ch <- resultDone{false}
		close(ch)
	}(&wg)
	return ch
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

func addTotal(column int) {
	switch column {
	case 1:
		totalErrMP++
	case 2:
		totalErrDP++
	// case 3:
	// 	totalErrEN++
	case 4:
		totalErrL++
	// case 5:
	// 	totalErrCN++
	case 6:
		totalErrRT++
	case 7:
		totalErrG++
	case 8:
		totalErrED++
	case 9:
		totalErrTD++
	}
}

func GenerateRowTotals() ([]string, []string) {
	totalsNamesRow := []string{"", TotalErrors + MedicalPlan, TotalErrors + DentalPlan, TotalErrors + Language, TotalErrors + RelationshipType, TotalErrors + Gender, TotalErrors + EffectiveDate, TotalErrors + TerminationDate}
	totalsRow := []string{"", stringFormatTotal(totalErrMP), stringFormatTotal(totalErrDP), stringFormatTotal(totalErrL), stringFormatTotal(totalErrRT), stringFormatTotal(totalErrG), stringFormatTotal(totalErrED), stringFormatTotal(totalErrTD)}
	return totalsNamesRow, totalsRow
}

func ResetTotals() {
	totalErrMP = 0
	totalErrDP = 0
	//totalErrEN = 0
	totalErrL = 0
	//totalErrCN = 0
	totalErrRT = 0
	totalErrG = 0
	totalErrED = 0
	totalErrTD = 0
	totalData = 0
}

func stringFormatTotal(totalErr int) string {
	return strconv.FormatInt(int64(totalErr), 10) + "/" + strconv.FormatInt(int64(totalData), 10)
}

func lineCounter(r *csv.Reader) (int, error) {
	count := 0

	for {
		_, err := r.Read()
		count++
		switch {
		case err == io.EOF:
			return count, nil

		case err != nil:
			return count, err
		}
	}
}

// func main() {
// 	//Test data set
// 	//err := CheckCSV("testOne.csv", "testTwo.csv")
// 	//Official data set
// 	err := CheckCSV("clientData.csv", "ourData.csv")
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	ResetTotals()
// }

func (e Err) Error() string {
	return string(e)
}
