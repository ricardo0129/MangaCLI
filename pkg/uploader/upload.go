package uploader

import (
    "os"
    "bytes"
    "mime/multipart"
    "path/filepath"
    "net/http"
    "fmt"
)

func UploadFile(file_path string, url string) {
    body := &bytes.Buffer{}
    formWriter := multipart.NewWriter(body)

    err := formWriter.WriteField("sender", "web")

    //err = formWriter.WriteField("parent", "0")
    bodyWriter, err := formWriter.CreateFormFile("file", filepath.Base(file_path))

    b, err := os.ReadFile(file_path)

    _, err = bodyWriter.Write(b)

    formWriter.Close()

    req, err := http.NewRequest("POST", url, body)

    req.Header.Set("Content-Type", formWriter.FormDataContentType())

    client := &http.Client{}
    res, err := client.Do(req)
    if err != nil {
        fmt.Println(err)
    }
    defer res.Body.Close()
}
