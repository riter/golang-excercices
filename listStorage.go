package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"github.com/gorilla/mux"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type Document struct {
	Id   string
	Name string
	Size int64
}

var directoryPath = "./storage"
var routeListStorage = "/storages"

func main() {
	r := mux.NewRouter()
	r.HandleFunc(routeListStorage, GetDocuments).Methods("GET")
	log.Fatal(http.ListenAndServe(":9000", r))
}

func GetDocuments(w http.ResponseWriter, r *http.Request) {
	filesInfo := getFilesInfo(directoryPath)
	files := mapFilesInfoToDocument(filesInfo)

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(files)
}

func mapFilesInfoToDocument(filesInfo []os.FileInfo) []Document {
	var files []Document
	for _, f := range filesInfo {
		codeMD5, err := hashFileMD5(directoryPath + "/" + f.Name())
		if err != nil {
			log.Fatal(err)
		}
		files = append(files, Document{Id: codeMD5, Name: f.Name(), Size: f.Size()})
	}
	return files
}

func hashFileMD5(filePath string) (string, error) {
	var returnMD5String string

	file, err := os.Open(filePath)
	if err != nil {
		return returnMD5String, err
	}
	defer file.Close()

	hash := md5.New()

	if _, err := io.Copy(hash, file); err != nil {
		return returnMD5String, err
	}

	hashInBytes := hash.Sum(nil)[:16]
	returnMD5String = hex.EncodeToString(hashInBytes)
	return returnMD5String, nil
}

func getFilesInfo(directory string) []os.FileInfo {
	files, err := ioutil.ReadDir(directory)
	if err != nil {
		log.Fatal(err)
	}
	return files
}