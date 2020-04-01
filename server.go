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

	if len(files.File) != 2 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Needs to get 2 files got ", len(files.File))
		//log.Fatalf("Needs to get 2 files got %d", len(files.File))
		return
	}
	// var filesCSV []os.FileInfo
	var fileNames []string
	for key, v := range files.File {
		for _, f := range v {
			file, err := f.Open()
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprint(w, "Got error ", err.Error())
				return
			}
			//defer file.Close()

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
			file.Close()
			//filesCSV = append(filesCSV, file)
		}
		fileNames = append(fileNames, key)
	}
	fileOne, err := files.File[fileNames[0]][0].Open()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Got error ", err.Error())
		return
	}
	fileTwo, err := files.File[fileNames[1]][0].Open()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Got error ", err.Error())
		return
	}

	CheckCSV(fileOne, fileTwo, w)
}
