package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strings"
)

var port string

//removed connection cap as we are already limiting the number of connections on the real server. Can reinclude if needed
//const maxConnections int = 10

var serverPort string
var serverIp string
var requestParam string

//removed connection cap as we are already limiting the number of connections on the real server. Can reinclude if needed
//var semaphore = make(chan struct{}, maxConnections)

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

		//removed connection cap as we are already limiting the number of connections on the real server. Can reinclude if needed
		//semaphore <- struct{}{}

		connection, err := listener.Accept()

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
		//removed connection cap as we are already limiting the number of connections on the real server. Can reinclude if needed
		//<-semaphore
	}()
	reader := bufio.NewReader(connection)
	request, err := http.ReadRequest(reader)

	if err != nil {
		fmt.Println(err)
		return
	}
	urlParser(request)
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

func urlParser(request *http.Request) {

	// Split the port and the ip from the Host and save each in a variable
	hostUrl := strings.Split(request.Host, ":")
	serverPort = hostUrl[1]
	serverIp = hostUrl[0]

	fmt.Println(hostUrl[0])
	fmt.Println(hostUrl[1])

	// Get the path from the URL
	requestParam = request.URL.Path
	fmt.Println(requestParam)

}

func getHandler(request *http.Request, connection net.Conn) {

	response, err := http.Get("http://" + serverIp + ":" + serverPort + requestParam)
	fmt.Print(response)

	if err != nil {
		fmt.Println("Helllo this is error", err)
		return
	}

	// Send response to client
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Print("body")

	responseStatus := response.Status
	contentType := response.Header.Get("Content-Type")
	responseHandler(connection, responseStatus, contentType, body)

}

func responseHandler(connection net.Conn, responseStatus string, contentType string, content []byte) {
	response := "HTTP/1.1 " + responseStatus + "\r\n" +
		"Content-Type: " + contentType + "\r\n" +
		"Content-Length: " + fmt.Sprintf("%d", len(content)) + "\r\n" +
		"\r\n" // Empty line separating headers and body..

	//Send headers followed by file content
	fmt.Println("Response is ready to be sent")
	connection.Write([]byte(response))
	connection.Write(content)
	fmt.Println("Response sent")
}
