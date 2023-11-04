package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

func main() {
    // Check if the required command-line arguments are provided.
    if len(os.Args) != 4 {
        fmt.Println("Usage: go run client.go <request method> <port> <resource path>")
        return
    }

    requestMethod := os.Args[1]
    port := os.Args[2]
    localFilePath := os.Args[3]

    // Define the URL with the specified port and resource path.
    url := fmt.Sprintf("http://localhost:" + port + "/" + localFilePath)


    if requestMethod == "GET" {
        // Handle GET request to retrieve a file.
        retrieveFile(url)
    } else if requestMethod == "POST" {
        // Handle POST request to send a file.
        sendFile(url, localFilePath, requestMethod)
    } else {
        fmt.Println("Invalid request method. Use 'GET' or 'POST'.")
    }

}

func getFileNameFromURL(url string) string {
    parts := strings.Split(url, "/")
    return parts[len(parts)-1]
}

func retrieveFile(url string) {
    // Create an HTTP client.
    client := &http.Client{}

    // Create an HTTP GET request with the specified URL.
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        fmt.Println("Error creating request:", err)
        return
    }

    // Send the GET request.
    res, err := client.Do(req)
    if err != nil {
        fmt.Println("Error sending request:", err)
        return
    }

    defer res.Body.Close()

    // Check if the response status is OK (200).
    if res.StatusCode != http.StatusOK {
        fmt.Println("Request failed with status:", res.Status)
        return
    }

    // Create a local file for writing the response.
    fileName := getFileNameFromURL(url)
    localFile, err := os.Create(fileName)
    if err != nil {
        fmt.Println("Error creating local file:", err)
        return
    }
    defer localFile.Close()

    // Copy the response body to the local file.
    _, err = io.Copy(localFile, res.Body)
    if err != nil {
        fmt.Println("Error copying response to file:", err)
        return
    }

    fmt.Printf("File saved as %s\n", fileName)
}

func sendFile(url string, localFilePath string, requestMethod string) {
    // Open and read the local file.
    localFile, err := os.Open(localFilePath)
    if err != nil {
        fmt.Println("Error opening local file:", err)
        return
    }
    defer localFile.Close()

    // Create an HTTP client.
    client := &http.Client{}

    // Create a buffer to store the file content.
    var buf bytes.Buffer
    _, err = io.Copy(&buf, localFile)
    if err != nil {
        fmt.Println("Error reading local file:", err)
        return
    }

    // Create an HTTP request with the specified method and the file content as the request body.
    req, err := http.NewRequest(requestMethod, url, &buf)
    if err != nil {
        fmt.Println("Error creating request:", err)
        return
    }

    // Send the request.
    res, err := client.Do(req)
    if err != nil {
        fmt.Println("Error sending request:", err)
        return
    }

    defer res.Body.Close()

    // Check if the response status is OK (200).
    if res.StatusCode != http.StatusCreated {
        fmt.Println("Request failed with status:", res.Status)
        return
    }

    fmt.Println("File successfully sent and processed by the server.")
}
