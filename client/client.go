
package main

import (
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
    resourcePath := os.Args[3]

    // Define the URL with the specified port and resource path.
    url := fmt.Sprintf("http://localhost:%s%s", port, resourcePath)

    // Create an HTTP client.
    client := &http.Client{}

    // Create an HTTP request with the specified method and an optional request body.
    var reqBody io.Reader // You can set the request body if needed.
    req, err := http.NewRequest(requestMethod, url, reqBody)
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

func getFileNameFromURL(url string) string {
    parts := strings.Split(url, "/")
    return parts[len(parts)-1]
}
