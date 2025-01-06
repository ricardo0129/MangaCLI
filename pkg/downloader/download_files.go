package downloader

import (
    "fmt"
    "io"
    "net/http"
)

func downloadImage(url string) ([]byte, error) {
    res, err := http.Get(url)
    if err != nil {
        return nil, err
    }
    defer res.Body.Close()
    if res.StatusCode != 200 {
        return nil, fmt.Errorf("bad status: %s", res.Status)
    }

    return io.ReadAll(res.Body)
}

func DownloadFiles(urls []string) []Page {
    var pages []Page
    for i, url := range urls {
        var page Page
        image , _ := downloadImage(url)
        page.Data = image
        page.Page = uint(i)
        pages = append(pages, page)
    }
    return pages
}
