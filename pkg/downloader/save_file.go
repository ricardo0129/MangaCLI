package downloader

import (
    "os"
    "archive/zip"
    "fmt"
)

type Page struct {
    Data []byte
    Page uint
}
func SaveAsCBZ(outputFileName string, content []Page) error {
    cbzFile, err := os.Create(outputFileName)
    if err != nil {
        return err
    }
    defer cbzFile.Close()
    zipWriter := zip.NewWriter(cbzFile)
    defer zipWriter.Close()

    for i, page := range content {
        fileName := fmt.Sprintf("%03d.jpg", i)
        imageFile, err := zipWriter.Create(fileName)
        if err != nil {
            return err
        }
        _, err = imageFile.Write(page.Data)
        if err != nil {
            return err
        }
    }
    return nil
}
