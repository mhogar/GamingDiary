package main

import (
	"fmt"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"regexp"

	"github.com/binarysoupdev/go-extensions/json"
)

const (
	PATH  = "/media/ryan/Seagate Portable Drive/Videos/Gaming Diary/(1) The Sunshine Story"
	REGEX = `meta([0-9]+).txt$`
	GLOB  = "meta*[0-9].txt"
)

type Data struct {
	Title  string `json:"title"`
	Dates  string `json:"dates"`
	Videos []Video
}

type Video struct {
	Title       string
	Description string
	Thumbnail   string
	Video       string
}

func main() {
	t := template.Must(template.ParseFiles("template.gohtml"))
	regex := regexp.MustCompile(REGEX)

	//-- load data
	data, err := json.UnmarshalFile[Data](filepath.Join(PATH, "index.json"))
	if err != nil {
		log.Fatal(err)
	}

	files, err := filepath.Glob(filepath.Join(PATH, "src", GLOB))
	if err != nil {
		log.Fatal(err)
	}

	data.Videos = make([]Video, len(files))
	for i, file := range files {
		index := regex.FindStringSubmatch(file)[1]

		data.Videos[i] = Video{
			Title:     fmt.Sprintf("Chapter %s", index),
			Thumbnail: fmt.Sprintf("src/t%s.png", index),
			Video:     fmt.Sprintf("src/v%s.mp4", index),
		}
	}

	//-- execute the template
	out := filepath.Join(PATH, "index.html")

	file, err := os.Create(out)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	if err := t.Execute(file, data); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("+ %s\n", out)
}
