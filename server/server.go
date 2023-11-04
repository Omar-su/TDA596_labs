package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var port string

const maxConnections int = 10

var semaphore = make(chan struct{}, maxConnections)

func main() {

	if len(os.Args) != 2 {
		fmt.Println("Program requires port number as input")
		return
	}

	port = os.Args[1]

	listener, err := net.Listen("tcp", ":"+port)

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Server is listening on port ", port)

	for {

		semaphore <- struct{}{}

		connection, err := listener.Accept()
		fmt.Println("Semaphore size: ", len(semaphore))
		if err != nil {
			fmt.Println(err)
			continue
		}
		go handleConnection(connection)
	}

}

func handleConnection(connection net.Conn) {

	defer func() {
		connection.Close()
		<-semaphore
	}()
	reader := bufio.NewReader(connection)
	request, err := http.ReadRequest(reader)

	if err != nil {
		fmt.Println(err)
		return
	}
	switch request.Method {
	case "GET":
		fmt.Println("GET request received")
		getHandler(request,connection)

	case "POST":
		fmt.Println("POST request received")
		postHandler(request, connection)

	default:
		responseStatus := "501 Not Implemented"
		contentType := "text/html"
		fileContent := []byte("Not Implemented \n")
		responseHandler(connection,responseStatus,contentType, fileContent)

	}

}

func getHandler(request *http.Request, connection net.Conn) {	
	var fileContent []byte
	err := error(nil)
	uri := request.RequestURI
	fP := "." + uri // Assuming requested files are in the current directory.
	
	contentType := getContentType(request, fP) 
	var responseStatus string
	switch contentType	{

	case "notFound":
		responseStatus = "404 Not Found"
		contentType = "text/html"
		fileContent = []byte("Not Found \n")
	
	case "badRequest":
		responseStatus = "400 Bad Request"
		contentType = "text/html"
		fileContent = []byte("Bad Request\n")

	
	default: 
		responseStatus = "200 OK"
		fileContent, err = ioutil.ReadFile(fP)
		if err != nil {
			responseHandler(connection,"500 Internal Server Error","text/html", []byte("Internal Server Error\n"))
			fmt.Println(err)
			return
		}
	
	}
	
	
	responseHandler(connection,responseStatus,contentType, fileContent)
}

func postHandler(request *http.Request, connection net.Conn) {
	
	uri := request.RequestURI
	fileName := getFileNameFromURL(uri)
	localFile, err := os.Create(fileName)
	if err != nil {
		fmt.Println("Error creating local file:", err)
		responseHandler(connection,"500 Internal Server Error","text/html", []byte("Internal Server Error\n"))
		return
	}
	defer localFile.Close()

	// Copy the response body to the local file.
	_, err = io.Copy(localFile, request.Body)
	if err != nil {
		fmt.Println("Error copying response to file:", err)
		responseHandler(connection,"500 Internal Server Error","text/html", []byte("Internal Server Error\n"))
		return
	}

	responseHandler(connection,"201 Created","text/html", []byte("File saved as "+ fileName + "\n"))
}

func getFileNameFromURL(url string) string {
	parts := strings.Split(url, "/")
	return parts[len(parts)-1]
}

func responseHandler(connection net.Conn, responseStatus string, contentType string, content []byte ) {
	response := "HTTP/1.1 "+ responseStatus + "\r\n" + 
		"Content-Type: " + contentType + "\r\n" +
	"Content-Length: " + fmt.Sprintf("%d", len(content)) +"\r\n" + 
	"\r\n" // Empty line separating headers and body..
	
	
	//Send headers followed by file content

	connection.Write([]byte(response))
	connection.Write(content)
}

func getContentType(request *http.Request, fP string ) string {
	var contentType string 
	// Check if the requested file exists.
	
	extension := filepath.Ext(fP)
	switch strings.ToLower(extension) {
	case ".html":
		contentType = "text/html"
	case ".txt":
		contentType = "text/plain"
	case ".gif":
		contentType = "image/gif"
	case ".jpeg", ".jpg":
		contentType =  "image/jpeg"
	case ".css":
		contentType = "text/css"
	default:
		return  "badRequest"
	}
	if _, err := os.Stat(fP); os.IsNotExist(err) {
		// Respond with a 404 "Not Found" code.
		contentType = "notFound"
	}
	return contentType
}
