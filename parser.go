package main

import (
  "log"
  "net/http"
  "github.com/PuerkitoBio/goquery"
  "strconv"
  "encoding/json"
  "io/ioutil"
  "os"
);

func checkFile(filename string) error {
  _, err := os.Stat(filename)
  if os.IsNotExist(err) {
    _, err := os.Create(filename)
    if err != nil {
      return err
    }
  }
  return nil
}

func main() {

  type  CompanyDetail struct {
    Title       string `json:"title"`
    Location    string `json:"loc"`
    JobPosted   string `json:"jposted"`
    Website     string `json:"website"`
  }

  limit := 256
  filename := "myFile.json"

  err := checkFile(filename)
  if err != nil {
    log.Fatal(err)
  }

  for i := 1; i <= limit; i++ {
    response, err := http.Get("https://weworkremotely.com/remote-companies?page=" + strconv.Itoa(i))
    if err != nil {
      log.Fatal(err)
    }
    defer response.Body.Close()

    document, err := goquery.NewDocumentFromReader(response.Body)
    if err != nil {
      log.Fatal("Error loading Http respose body. ", err)
    }

    var data = []CompanyDetail{}

    document.Find("#remote-companies").Each(func(i int, s *goquery.Selection) {

      link := s.Find(".tooltip a").AttrOr("href", "none")
      res, err := http.Get("https://weworkremotely.com" + link)
      if err != nil {
        log.Fatal(err)
      }
      defer res.Body.Close()

      detailDocument, err := goquery.NewDocumentFromReader(res.Body)
      if err != nil {
        log.Fatal("Error loading Http respose body. ", err)
      }

      detailDocument.Find(".listing-header-container").Each(func(i int, sel *goquery.Selection) {
        detailDocument.Find(".tools").Children().Remove().End()

        //read file
        file, err := ioutil.ReadFile(filename)
        if err != nil {
          log.Fatal(err)
        }
        json.Unmarshal(file, &data)

        site := sel.Find(".listing-tools a").AttrOr("href", "none")
        newStruct := &CompanyDetail {
          Title: sel.Find("h1").Text(),
          Location: sel.Find("h3").Text(),
          JobPosted: sel.Find("h4").Text(),
          Website: site,
        }

        data = append(data, *newStruct)
        dataBytes, err := json.Marshal(data)
        if err != nil {
          log.Fatal(err)
        }

        err = ioutil.WriteFile(filename, dataBytes, 0644)
        if err != nil {
          log.Fatal(err)
        }
      })
    })
  }

}
