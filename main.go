package main

import (
	"fmt"
	"html/template"
	"log"
	"os"
	"path/filepath"

	"github.com/binarysoupdev/go-extensions/json"
)

const (
	PATH = "/media/ryan/Seagate Portable Drive/Videos/Ryans Gaming Diary/(1) The Sunshine Story"
	OUT  = "out.html"
)

type Data struct {
	Title string `json:"title"`
	Dates string `json:"dates"`
}

func main() {
	t := template.Must(template.ParseFiles("template.gohtml"))

	data, err := json.UnmarshalFile[Data](filepath.Join(PATH, "index.json"))
	if err != nil {
		log.Fatal(err)
	}

	file, err := os.Create(OUT)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	if err := t.Execute(file, data); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("+ %s\n", OUT)
}
