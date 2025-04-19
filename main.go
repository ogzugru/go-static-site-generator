package main

import (
	"fmt"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"gopkg.in/yaml.v2"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var md = goldmark.New(
	goldmark.WithExtensions(
		extension.GFM, // GitHub Flavored Markdown (tablolar, vs.)
	),
	goldmark.WithParserOptions(
		parser.WithAutoHeadingID(),
	),
)

type Config struct {
	SiteName  string `yaml:"site_name"`
	Author    string `yaml:"author"`
	OutputDir string `yaml:"output_dir"`
}

type FrontMatter struct {
	Title string `yaml:"title"`
	Date  string `yaml:"date"`
}

type Post struct {
	Title    string
	Date     string
	Author   string
	Body     template.HTML
	Content  template.HTML
	Filename string
}

func main() {
	// Config oku
	cfgData, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		log.Fatalf("Config dosyası okunamadı: %v", err)
	}
	var cfg Config
	if err := yaml.Unmarshal(cfgData, &cfg); err != nil {
		log.Fatalf("Config ayrıştırılamadı: %v", err)
	}

	// Template'leri yükle
	baseTpl := template.Must(template.ParseFiles("templates/base.html"))
	postTpl := template.Must(template.ParseFiles("templates/post.html"))
	indexTpl := template.Must(template.ParseFiles("templates/index.html"))

	// İçerikleri oku
	files, err := ioutil.ReadDir("content")
	if err != nil {
		log.Fatal(err)
	}

	os.MkdirAll(cfg.OutputDir, os.ModePerm)

	var posts []Post

	for _, file := range files {
		if filepath.Ext(file.Name()) != ".md" {
			continue
		}

		data, err := ioutil.ReadFile(filepath.Join("content", file.Name()))
		if err != nil {
			log.Printf("Dosya okunamadı %s: %v", file.Name(), err)
			continue
		}

		content := string(data)

		var fm FrontMatter
		var markdownContent string

		// Frontmatter ayrıştır
		if strings.HasPrefix(content, "---") {
			parts := strings.SplitN(content, "---", 3)
			if len(parts) >= 3 {
				yamlData := parts[1]
				markdownContent = parts[2]
				if err := yaml.Unmarshal([]byte(yamlData), &fm); err != nil {
					log.Printf("Frontmatter ayrıştırılamadı %s: %v", file.Name(), err)
				}
			}
		} else {
			markdownContent = content
		}

		// Markdown içeriği HTML'e çevir
		var buf strings.Builder
		if err := md.Convert([]byte(markdownContent), &buf); err != nil {
			log.Fatalf("Markdown dönüştürme hatası: %v", err)
		}
		html := buf.String()
		html = strings.ReplaceAll(html, "<table>", `<table class="table table-striped">`)

		filename := strings.TrimSuffix(file.Name(), ".md") + ".html"

		post := Post{
			Title:    fm.Title,
			Date:     fm.Date,
			Author:   cfg.Author,
			Body:     template.HTML(html),
			Filename: filename,
		}

		// Post içeriğini oluştur
		var body strings.Builder
		if err := postTpl.Execute(&body, post); err != nil {
			log.Printf("Post template hatası %s: %v", file.Name(), err)
			continue
		}
		post.Content = template.HTML(body.String())

		outputPath := filepath.Join(cfg.OutputDir, filename)
		f, err := os.Create(outputPath)
		if err != nil {
			log.Printf("Dosya oluşturulamadı %s: %v", outputPath, err)
			continue
		}
		defer f.Close()

		if err := baseTpl.Execute(f, post); err != nil {
			log.Printf("Base template hatası %s: %v", file.Name(), err)
		}

		fmt.Println("Oluşturuldu:", outputPath)
		posts = append(posts, post)
	}

	// Gönderileri tarihe göre tersten sırala (en yeni yukarıda olsun)
	for i, j := 0, len(posts)-1; i < j; i, j = i+1, j-1 {
		posts[i], posts[j] = posts[j], posts[i]
	}

	// Pagination
	postsPerPage := 5
	totalPages := (len(posts) + postsPerPage - 1) / postsPerPage

	for i := 0; i < totalPages; i++ {
		start := i * postsPerPage
		end := start + postsPerPage
		if end > len(posts) {
			end = len(posts)
		}
		pagePosts := posts[start:end]

		var prevPage, nextPage string
		if i > 0 {
			if i == 1 {
				prevPage = "/index.html"
			} else {
				prevPage = fmt.Sprintf("/page/%d/index.html", i)
			}
		}
		if i < totalPages-1 {
			nextPage = fmt.Sprintf("/page/%d/index.html", i+2)
		}

		outputDir := cfg.OutputDir
		if i > 0 {
			outputDir = filepath.Join(cfg.OutputDir, "page", fmt.Sprintf("%d", i+1))
			os.MkdirAll(outputDir, os.ModePerm)
		}

		outputFilePath := filepath.Join(outputDir, "index.html")
		outputFile, err := os.Create(outputFilePath)
		if err != nil {
			log.Fatalf("Sayfa oluşturulamadı: %v", err)
		}
		defer outputFile.Close()

		if err := indexTpl.Execute(outputFile, struct {
			SiteName string
			Posts    []Post
			PrevPage string
			NextPage string
		}{
			SiteName: cfg.SiteName,
			Posts:    pagePosts,
			PrevPage: prevPage,
			NextPage: nextPage,
		}); err != nil {
			log.Fatalf("Sayfa template hatası: %v", err)
		}

		fmt.Println("Sayfa oluşturuldu:", outputFilePath)
	}

	// Sunucu başlat
	fmt.Println("Sunucu başlatılıyor: http://localhost:8080")
	http.Handle("/", http.FileServer(http.Dir(cfg.OutputDir)))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
