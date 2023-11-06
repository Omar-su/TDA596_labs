package main

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
	"os"
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

	fmt.Println(" Proxy Server is listening on port ", port)

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
		getHandler(request, connection)

	default:
		responseStatus := "501 Not Implemented"
		contentType := "text/html"
		fileContent := []byte("Not Implemented \n")
		responseHandler(connection, responseStatus, contentType, fileContent)

	}
}

func getHandler(request  *http.Request, connection net.Conn) {
	// url := request.RequestURI
	//request.ParseForm()
	fmt.Println("Printing request")
	//print port and url 
	fmt.Println(request.Host)

	// Split the port and the ip from the Host and save each in a variable
	hostIp := strings.Split(request.Host, ":")
	fmt.Println(hostIp[0])
	fmt.Println(hostIp[1])
	// Split the port and the ip from the Host and save each in a variable
	// ip, host := strings.Split(request.Host, ":")

	fmt.Println(request.URL)
	// Get the path from the URL
	path := request.URL.Path 
	fmt.Println(path)
	// Create an HTTP GET request with the specified URL.
	// req, err := http.NewRequest("GET", url, nil)
	// if err != nil {
	// 		fmt.Println("Error creating request:", err)
	// 		return
	// }

	// Send the GET request.
	// res, err := client.Do(req)
	// if err != nil {
	// 		fmt.Println("Error sending request:", err)
	// 		return
	// }
}

func responseHandler(connection net.Conn, responseStatus string, contentType string, content []byte) {
	response := "HTTP/1.1 " + responseStatus + "\r\n" +
		"Content-Type: " + contentType + "\r\n" +
		"Content-Length: " + fmt.Sprintf("%d", len(content)) + "\r\n" +
		"\r\n" // Empty line separating headers and body..

	//Send headers followed by file content

	connection.Write([]byte(response))
	connection.Write(content)
}
