package main

import (
	"crypto/rand"
	"fmt"
	"io/ioutil"
	"log"
	"mime"
	"net/http"
	"os"
	"io"
	"strconv"
	"path/filepath"
)

const maxUploadSize = 25 * 1024 * 1024 // 25 mb
const uploadPath = "./tmp"

// id -> user-facing filename, for /getallfiles
var filenames map[string]string
// id -> internal filename, for /files/{fileId}
var internalFilenames map[string]string

func main() {
	filenames = make(map[string]string, 0)
	internalFilenames = make(map[string]string, 0)

	http.HandleFunc("/upload", uploadFileHandler())
	http.HandleFunc("/hello", helloHandler())
	
    fs := http.FileServer(http.Dir(uploadPath))
	http.Handle("/files/", http.StripPrefix("/files", fs))

	http.HandleFunc("/allfiles", allFilesHandler())
	http.HandleFunc("/download", downloadFileHandler())

	// create ./tmp, overwriting any previous files that may have been present
	fmt.Printf("overwriting ./tmp")
	err1 := os.RemoveAll(uploadPath)
	if err1 != nil {
		log.Fatal("Couldn't remove ./tmp %v", err1)
		return
	}

	err2 := os.Mkdir(uploadPath, 0700)
	if err2 != nil {
		log.Fatal("Couldn't create ./tmp %v", err2)
		return
	}

	log.Print("Server started on localhost:8080, use /upload for uploading files and /files/{fileName} for downloading")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func downloadFileHandler() http.HandlerFunc {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		// filename is id + file-extension
		// internalFilenames[
		
	fileId := request.URL.Query().Get("file")
	if fileId == "" {
		//Get not set, send a 400 bad request
		http.Error(writer, "Get 'file' not specified in url.", 400)
		return
	}

	fileName := internalFilenames[fileId]
	fmt.Println("Client requests: " + fileId + ", " + fileName)

	//Check if file exists and open
	Openfile, err := os.Open(fileName)
	defer Openfile.Close() //Close after function return
	if err != nil {
		//File not found, send 404
		http.Error(writer, "File not found.", 404)
		return
	}

	//File is found, create and send the correct headers

	//Get the Content-Type of the file
	//Create a buffer to store the header of the file in
	FileHeader := make([]byte, 512)
	//Copy the headers into the FileHeader buffer
	Openfile.Read(FileHeader)
	//Get content type of file
	FileContentType := http.DetectContentType(FileHeader)

	//Get the file size
	FileStat, _ := Openfile.Stat()                     //Get info from file
	FileSize := strconv.FormatInt(FileStat.Size(), 10) //Get file size as a string

	//Send the headers
	writer.Header().Set("Content-Disposition", "attachment; filename="+fileName)
	writer.Header().Set("Content-Type", FileContentType)
	writer.Header().Set("Content-Length", FileSize)

	//Send the file
	//We read 512 bytes from the file already, so we reset the offset back to 0
	Openfile.Seek(0, 0)
	io.Copy(writer, Openfile) //'Copy' the file to the client
	return

	})
}

func helloHandler() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        log.Print("hello")
    	w.Write([]byte("hello world"))
	})
}

func allFilesHandler() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		responseStr := ""
		for id, filename := range filenames {
			responseStr = responseStr + id + ", " + filename + "\n" 
		}
		w.Write([]byte(responseStr))
	})
}

func uploadFileHandler() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Only handle POST
		if r.Method == "GET" {
			return
		}

		// Can we parse the form?
		if err := r.ParseMultipartForm(maxUploadSize); err != nil {
			fmt.Printf("Could not parse multipart form: %v\n", err)
			renderError(w, "CANT_PARSE_FORM", http.StatusInternalServerError)
			return
		}

		// Actually parse the form
		file, fileHeader, err := r.FormFile("uploadFile")
		if err != nil {
			fmt.Printf("Can't read file header: %v\n", err)
			renderError(w, "INVALID_FILE", http.StatusBadRequest)
			return
		}

		defer file.Close()
		
		// Get and validate filesize
		fileSize := fileHeader.Size
		fmt.Printf("File size (bytes): %v\n", fileSize)
		if fileSize > maxUploadSize {
			renderError(w, "FILE_TOO_BIG", http.StatusBadRequest)
			return
		}
		fileBytes, err := ioutil.ReadAll(file)
		if err != nil {
			fmt.Printf("Can't read the file: %v\n", err)
			renderError(w, "INVALID_FILE", http.StatusBadRequest)
			return
		}

		// filetype
		// check file type, detectcontenttype only needs the first 512 bytes
		detectedFileType := http.DetectContentType(fileBytes)
		// switch detectedFileType {
		// case "image/jpeg", "image/jpg":
		// case "image/gif", "image/png":
		// case "application/pdf":
		// 	break
		// default:
		// 	renderError(w, "INVALID_FILE_TYPE", http.StatusBadRequest)
		// 	return
		// }
		

		// Get filename and uuid (to differentiate duplicate filenames)
		fileName := fileHeader.Filename
		fileId := randToken(8)
		
		fileEndings, err := mime.ExtensionsByType(detectedFileType)
		if err != nil {
			renderError(w, "CANT_READ_FILE_TYPE", http.StatusInternalServerError)
			return
		}
		filePath := filepath.Join(uploadPath, fileId + fileEndings[0])
		fmt.Printf("Path: %s\n", filePath)

		fmt.Printf("File: %s, ID: %s\n", fileName, fileId)

		// write file
		newFile, err := os.Create(filePath)
		if err != nil {
			fmt.Printf("Can't write file %v\n", err)
			renderError(w, "CANT_WRITE_FILE", http.StatusInternalServerError)
			return
		}
		defer newFile.Close() // idempotent, okay to call twice
		if _, err := newFile.Write(fileBytes); err != nil || newFile.Close() != nil {
			renderError(w, "CANT_WRITE_FILE", http.StatusInternalServerError)
			return
		}
		filenames[fileId] = fileName
		internalFilenames[fileId] = filePath
		// filenames = append(filenames, fileName)
		w.Write([]byte("SUCCESS: " + fileId))
	})
}
// $ curl -F 'img_avatar=@/home/petehouston/hello.txt' 
// 		http://localhost/upload

func renderError(w http.ResponseWriter, message string, statusCode int) {
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte(message))
}

func randToken(len int) string {
	b := make([]byte, len)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}
