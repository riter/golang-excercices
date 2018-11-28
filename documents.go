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

const directoryPath = "./storage"

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/documents", GetDocuments).Methods("GET")
	r.HandleFunc("/documents/{id}",GetDocumentById).Methods("GET")
	r.HandleFunc("/documents",createDocument).Methods("POST")
	r.HandleFunc("/documents/{id}",deleteDocumentById).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":9000", r))
}

func GetDocumentById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	filesInfo := getFilesInfo(directoryPath)
	for _, f := range filesInfo {
		var document = mapFileInfoToDocument(f)
		if document.Id == params["id"] {
			json.NewEncoder(w).Encode(document)
			return
		}
	}
	http.Error(w, "Not found", http.StatusNotFound)
}

func createDocument(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Field File required", http.StatusBadRequest)
	}
	defer file.Close()

	out, err := os.Create(directoryPath + "/" + header.Filename)
	if err != nil {
		http.Error(w, "Unable to create the file for writing. Check your write access privilege", http.StatusInternalServerError)
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		http.Error(w, "upload fail", http.StatusInternalServerError)
	}

	fileInfo, _ := out.Stat()
	var document = mapFileInfoToDocument(fileInfo)
	json.NewEncoder(w).Encode(document)
}

func deleteDocumentById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	filesInfo := getFilesInfo(directoryPath)
	for _, f := range filesInfo {
		var document = mapFileInfoToDocument(f)
		if document.Id == params["id"] {
			var err = os.Remove(directoryPath + "/" + f.Name())
			if err != nil {
				http.Error(w, "Unable to remove the file.", http.StatusInternalServerError)
			}
			json.NewEncoder(w).Encode("Document deleted")
			return
		}
	}
	http.Error(w, "Not found", http.StatusNotFound)
}

func GetDocuments(w http.ResponseWriter, r *http.Request) {
	filesInfo := getFilesInfo(directoryPath)

	var files []Document
	for _, f := range filesInfo {
		var document = mapFileInfoToDocument(f)
		files = append(files, document)
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(files)
}

func mapFileInfoToDocument(fileInfo os.FileInfo) Document {
	codeMD5, err := hashFileMD5(directoryPath + "/" + fileInfo.Name())
	if err != nil {
		log.Fatal(err)
	}
	return Document{Id: codeMD5, Name: fileInfo.Name(), Size: fileInfo.Size()}
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