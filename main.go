package main

import (
	"fmt"
	"github.com/russross/blackfriday/v2"
	"gopkg.in/yaml.v2"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type Config struct {
	SiteName  string `yaml:"site_name"`
	Author    string `yaml:"author"`
	OutputDir string `yaml:"output_dir"`
}

type Post struct {
	Title   string
	Date    string
	Author  string
	Body    template.HTML
	Content template.HTML
}

func main() {
	// Config oku
	cfgData, _ := ioutil.ReadFile("config.yaml")
	var cfg Config
	yaml.Unmarshal(cfgData, &cfg)

	// Template yükle
	baseTpl := template.Must(template.ParseFiles("templates/base.html"))
	postTpl := template.Must(template.ParseFiles("templates/post.html"))

	files, err := ioutil.ReadDir("content")
	if err != nil {
		log.Fatal(err)
	}

	os.MkdirAll(cfg.OutputDir, os.ModePerm)

	for _, file := range files {
		data, _ := ioutil.ReadFile("content/" + file.Name())
		content := string(data)

		// Basit frontmatter ayrıştırma (örnek amaçlı)
		lines := strings.Split(content, "\n")
		var title, date string
		if len(lines) > 2 && lines[0] == "---" {
			title = strings.TrimPrefix(lines[1], "title: ")
			date = strings.TrimPrefix(lines[2], "date: ")
			content = strings.Join(lines[4:], "\n")
		}

		html := blackfriday.Run([]byte(content))

		post := Post{
			Title:  title,
			Date:   date,
			Author: cfg.Author,
			Body:   template.HTML(html),
		}

		var body strings.Builder
		postTpl.Execute(&body, post)
		post.Content = template.HTML(body.String())

		// Sonuç HTML dosyasını yaz
		outputFile := filepath.Join(cfg.OutputDir, strings.TrimSuffix(file.Name(), ".md")+".html")
		f, _ := os.Create(outputFile)
		baseTpl.Execute(f, post)
		fmt.Println("Oluşturuldu:", outputFile)
	}
}
