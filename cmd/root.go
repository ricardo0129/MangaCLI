package cmd

import (
    "fmt"
    "bufio"
    "os"
    "main/pkg/scraper"
    "strconv"
    "main/pkg/uploader"
    "main/pkg/utils"
)


func Run() {
    read := bufio.NewReader(os.Stdin)
    title, _ := read.ReadString('\n')

    res := scraper.SearchTitle(title)
    for i, manga := range res {
        fmt.Println(i, manga.Id, manga.Attributes.Title["en"])
    }
    choice := 0
    fmt.Scan(&choice)
    fmt.Println(res[choice].Attributes.Title["en"])
    title = res[choice].Attributes.Title["en"]

    id := res[choice].Id
    //chapters := scraper.GetChapters(id, 96, 0, "en")
    
    range_str, _ := read.ReadString('\n')
    chapter_range := utils.ParseRange(range_str)

    chapters := scraper.GetAllChapters(id, "en")

    download_dir := "downloads"
    url := "http://10.0.0.37:8085/api/library/upload"
    upload := true

    for i, chapter := range chapters {
        chapter_number, _ := strconv.Atoi(chapter.Attributes.Chapter)
        if utils.InsideRange(chapter_range, chapter_number) {
            fmt.Println(i, chapter.Id, chapter.Attributes)
            chapter_id := chapter.Id
            chapter_name := chapter.Attributes.Chapter 
            file_name := fmt.Sprintf("%s/%s_%s.cbz", download_dir, title, chapter_name)
            scraper.DownloadChapter(chapter_id, file_name)
            if upload {
                uploader.UploadFile(file_name, url)
            }
        }
    }
    //fmt.Scan(&choice)
}
