package main

import (
    "fmt"
    "io"
    "net/http"
    "strconv"
    "os"
    "archive/zip"
    "encoding/json"
    "bytes"
    "mime/multipart"
    "path/filepath"
)

type Manga struct {
    Id string 
    Attributes struct {
        Title map[string]string
        AltTitles []map[string]string
    }
}

type Chapter struct {
    Id string
    Attributes struct {
        Chapter string 
        Pages int
        Title string
        TranslatedLanguage string
    }
}


/*
type TopLevel struct {
    Data []byte `json:"data"`
}
*/

func search_title(title string) []Manga {


    url := "https://api.mangadex.org/manga"
    req, _ := http.NewRequest("GET", url, nil)

    q := req.URL.Query()
    q.Add("title", title)
    req.URL.RawQuery = q.Encode()

    res, err := http.Get(req.URL.String())
    if err != nil {
        return nil
    }
    defer res.Body.Close()

    body, err := io.ReadAll(res.Body)

    if err != nil {
        return nil
    }

    type TopLevel struct {
        Data []Manga
    }

    var topLevel TopLevel
    err = json.Unmarshal(body, &topLevel)


    if err != nil {
        fmt.Println(err)
        return nil
    }

    return topLevel.Data
}

func get_chapters(manga_id string, limit int, offset int, language string) []Chapter {
    url := "https://api.mangadex.org/manga/" + manga_id + "/feed"
    req, _ := http.NewRequest("GET", url, nil)
    q := req.URL.Query()
    q.Add("limit", strconv.Itoa(limit))
    q.Add("offset", strconv.Itoa(offset))
    q.Add("order[volume]", "desc")
    q.Add("order[chapter]", "desc")

    if language != "" {
        q.Add("translatedLanguage[]", language)
    }

    req.URL.RawQuery = q.Encode()
    res, err := http.Get(req.URL.String())

    if err != nil {
        fmt.Println(err)
    }
    defer res.Body.Close()

    body, err := io.ReadAll(res.Body)
    if err != nil {
        fmt.Println(err)
    }

    type TopLevel struct {
        Data []Chapter
    }

    var topLevel TopLevel
    err = json.Unmarshal(body, &topLevel)
    if err != nil {
        fmt.Println(err)
    }

    return topLevel.Data
}

func all_chapters(manga_id string, language string) []Chapter {
    var all_chapters []Chapter
    limit := 100
    offset := 0
    for {
        chapters := get_chapters(manga_id, limit, offset, language)
        all_chapters = append(all_chapters, chapters...)
        if len(chapters) < limit {
            break
        }
        offset += limit
    }
    return all_chapters
}

type Page struct {
    Data []byte
    Page uint
}

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

func saveAsCBZ(outputFileName string, content []Page) error {
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

func download_chapter(chapter_id string, outputFileName string) {
    url := "https://api.mangadex.org/at-home/server/" + chapter_id + "?forcePort443=false"
    req, _ := http.NewRequest("GET", url, nil)
    res, err := http.Get(req.URL.String())
    if err != nil {
        fmt.Println(err)
    }
    defer res.Body.Close()

    body, err := io.ReadAll(res.Body)
    if err != nil {
        fmt.Println(err)
    }

    type ChapterData struct {
        Result string 
        BaseUrl string 
        Chapter struct {
            Hash string 
            Data []string 
        } 
    }

    var chapterData ChapterData
    json.Unmarshal(body, &chapterData)

    var contents []Page

    for i, fileName := range chapterData.Chapter.Data {
        imageURL := fmt.Sprintf("%s/data/%s/%s", chapterData.BaseUrl, chapterData.Chapter.Hash, fileName)
        image, _ := downloadImage(imageURL)
        var page Page
        page.Data = image
        page.Page = uint(i)
        contents = append(contents, page)
    }

    err = saveAsCBZ(outputFileName, contents)
    if err != nil {
        fmt.Println(err)
    }
}

func upload_file(file_path string) {
    body := &bytes.Buffer{}
    formWriter := multipart.NewWriter(body)

    err := formWriter.WriteField("sender", "web")

    //err = formWriter.WriteField("parent", "0")
    bodyWriter, err := formWriter.CreateFormFile("file", filepath.Base(file_path))

    b, err := os.ReadFile(file_path)

    _, err = bodyWriter.Write(b)

    formWriter.Close()

    url := "http://10.0.0.37:8085/api/library/upload"
    req, err := http.NewRequest("POST", url, body)

    req.Header.Set("Content-Type", formWriter.FormDataContentType())

    client := &http.Client{}
    res, err := client.Do(req)
    if err != nil {
        fmt.Println(err)
    }
    defer res.Body.Close()
}


func main() {
    title := "hunter x hunter"
    res := search_title(title)
    for i, manga := range res {
        fmt.Println(i, manga.Id, manga.Attributes.Title["en"])
    }
    choice := 0
    fmt.Scan(&choice)
    title = res[choice].Attributes.Title["en"]

    id := res[choice].Id
    chapters := get_chapters(id, 96, 0, "en")

    for i, chapter := range chapters {
        fmt.Println(i, chapter.Id, chapter.Attributes)
    }
    fmt.Scan(&choice)

    chapter_id := chapters[choice].Id
    chapter := chapters[choice].Attributes.Chapter 

    download_dir := "downloads"

    file_name := fmt.Sprintf("%s/%s_%s.cbz", download_dir, title, chapter)
    download_chapter(chapter_id, file_name)
    upload_file(file_name)
}
