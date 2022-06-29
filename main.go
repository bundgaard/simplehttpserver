package main

import (
	"embed"
	"flag"
	"fmt"
	"html/template"
	"log"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var (
	port   = flag.String("port", "8080", "port to listen on")
	fp     = flag.String("path", "./", "path to expose")
	secure = flag.Bool("secure", false, "should we expose HTTPS")
)

//go:embed public/*
var videopage embed.FS

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func readDir(fp string) map[string][]string {
	dirEntries, err := os.ReadDir(fp)
	if err != nil {
		log.Println(err)
	}

	filesByExtension := make(map[string][]string)
	for _, entry := range dirEntries {
		if entry.IsDir() {
			continue
		}
		fileType := mime.TypeByExtension(filepath.Ext(entry.Name()))
		commaIdx := strings.Index(fileType, ";")
		if commaIdx > 0 {
			fileType = fileType[:commaIdx]
		}
		filesByExtension[fileType] = append(filesByExtension[fileType], entry.Name())
	}

	return filesByExtension
}

func videoPlayer(w http.ResponseWriter, r *http.Request) {

	tpl, err := template.ParseFS(videopage, "public/video.html")
	if err != nil {
		log.Fatal(err)
	}
	files := readDir(*fp)
	if err := tpl.Execute(w, files); err != nil {
		log.Fatal(err)
	}
}
func main() {
	flag.Parse()

	fmt.Printf("listening in %s on %s\n", *fp, *port)
	if *secure {
		fmt.Println("not implemented yet")
		return
	}

	readDir(*fp)

	root := http.NewServeMux()

	root.Handle("/", http.FileServer(http.Dir(*fp)))
	root.Handle("/static/", http.StripPrefix("/static", http.FileServer(http.FS(videopage))))
	root.HandleFunc("/video/", videoPlayer)

	if err := http.ListenAndServe(fmt.Sprintf(":%s", *port), root); err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		os.Exit(1)
	}

}
