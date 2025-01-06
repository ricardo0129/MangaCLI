package cmd

import (
    "fmt"
    "main/pkg/scraper"
    "main/pkg/uploader"
)


func Run() {
    title := "hunter x hunter"
    res := scraper.SearchTitle(title)
    for i, manga := range res {
        fmt.Println(i, manga.Id, manga.Attributes.Title["en"])
    }
    choice := 0
    fmt.Scan(&choice)
    title = res[choice].Attributes.Title["en"]

    id := res[choice].Id
    chapters := scraper.GetChapters(id, 96, 0, "en")

    for i, chapter := range chapters {
        fmt.Println(i, chapter.Id, chapter.Attributes)
    }
    fmt.Scan(&choice)

    chapter_id := chapters[choice].Id
    chapter := chapters[choice].Attributes.Chapter 

    download_dir := "downloads"

    file_name := fmt.Sprintf("%s/%s_%s.cbz", download_dir, title, chapter)
    scraper.DownloadChapter(chapter_id, file_name)

    url := "http://10.0.0.37:8085/api/library/upload"
    uploader.UploadFile(file_name, url)
}
