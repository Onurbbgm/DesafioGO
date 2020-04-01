package main

import (
	"fmt"
	"net/http"

	"github.com/gabriel-vasile/mimetype"
)

//DataAnalysisServer setup of the server
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

//Get files and check if they are valid in server request
func getFiles(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 << 20)

	files := r.MultipartForm

	if !CheckNumberOfFiles(len(files.File), w) {
		return
	}

	var fileNames []string
	for key, v := range files.File {
		for _, f := range v {
			file, err := f.Open()
			if CheckErrorExist(err, w, "Got error") {
				return
			}

			mime, err := mimetype.DetectReader(file)
			if CheckErrorExist(err, w, "Got error") {
				return
			}

			if !CheckFileExtension(mime.Extension(), w) {
				return
			}
			file.Close()
		}
		fileNames = append(fileNames, key)
	}

	fileOne, _ := files.File[fileNames[0]][0].Open()
	fileTwo, _ := files.File[fileNames[1]][0].Open()

	CheckCSV(fileOne, fileTwo, w)
}

//CheckFileExtension check if it is a validate extension
func CheckFileExtension(extension string, w http.ResponseWriter) bool {
	if extension != ".csv" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Expected contentType text/csv, got ", extension)
		return false
	}
	return true
}

//CheckErrorExist se if there is an error e send response
func CheckErrorExist(err error, w http.ResponseWriter, text string) bool {
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, text+err.Error())
		return true
	}
	return false
}

//CheckNumberOfFiles see if got the rigth number of files to check
func CheckNumberOfFiles(number int, w http.ResponseWriter) bool {
	if number != 2 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Needs to get 2 files got ", number)
		return false
	}
	return true
}
