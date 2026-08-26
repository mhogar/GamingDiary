package main

import (
	"fmt"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/binarysoupdev/go-extensions/json"
)

const (
	PATH  = "/media/ryan/Seagate Portable Drive/Videos/Gaming Diary/(1) The Sunshine Story"
	REGEX = `meta([0-9]+).txt$`
	GLOB  = "meta*[0-9].txt"
)

type Data struct {
	Title    string `json:"title"`
	Dates    string `json:"dates"`
	Prologue Video
	Videos   []Video
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

	data.Prologue = buildVideo(files[0], "00")
	files = files[1:]

	data.Videos = make([]Video, len(files))
	for i, file := range files {
		data.Videos[i] = buildVideo(file, regex.FindStringSubmatch(file)[1])
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

func buildVideo(path string, index string) Video {
	meta, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	lines := strings.Split(string(meta), "\n")

	return Video{
		Title:       fmt.Sprintf("Chapter %s | %s\n", index, strings.SplitN(lines[2], " | ", 2)[0]),
		Description: lines[5],
		Thumbnail:   fmt.Sprintf("src/t%s.png", index),
		Video:       fmt.Sprintf("src/v%s.mp4", index),
	}
}
