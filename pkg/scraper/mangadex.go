package scraper

import (
    "net/http"
    "io"
    "strconv"
    "fmt"
    "encoding/json"
    "main/pkg/downloader"
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


func SearchTitle(title string) []Manga {
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

func GetChapters(manga_id string, limit int, offset int, language string) []Chapter {
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

func AllChapters(manga_id string, language string) []Chapter {
    var all_chapters []Chapter
    limit := 100
    offset := 0
    for {
        chapters := GetChapters(manga_id, limit, offset, language)
        all_chapters = append(all_chapters, chapters...)
        if len(chapters) < limit {
            break
        }
        offset += limit
    }
    return all_chapters
}


func DownloadChapter(chapter_id string, outputFileName string) {
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

    var contents []downloader.Page
    
    var urls []string
    for _, fileName := range chapterData.Chapter.Data {
        imageURL := fmt.Sprintf("%s/data/%s/%s", chapterData.BaseUrl, chapterData.Chapter.Hash, fileName)
        urls = append(urls, imageURL)
    }

    contents = downloader.DownloadFiles(urls)

    err = downloader.SaveAsCBZ(outputFileName, contents)
    if err != nil {
        fmt.Println(err)
    }
}
