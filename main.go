package main

import (
	"bufio"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type PageData struct {
	Greeting string
	Text     string
	Banner   string
	Result   string
	Error    string
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	renderTemplate(w, PageData{
		Greeting: "Abdulwasiucodes",
	})
}

func asciiHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "405 Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	text := r.FormValue("text")
	banner := r.FormValue("banner")

	data := PageData{
		Greeting: "Abdulwasiucodes",
		Text:     text,
		Banner:   banner,
	}

	if text == "" || banner == "" {
		data.Error = "Please fill in all fields"
		renderTemplate(w, data)
		return
	}

	data.Result = generateASCII(text, banner)

	renderTemplate(w, data)
}

func generateASCII(text string, banner string) string {
	font, err := loadBanner(banner + ".txt")
	if err != nil {
		return "Error loading banner file"
	}

	text = strings.ReplaceAll(text, "\r\n", "\n")

	var result strings.Builder

	lines := strings.Split(text, "\n")

	for i, line := range lines {
		result.WriteString(renderText(line, font))

		if i < len(lines)-1 {
			result.WriteString("\n")
		}
	}

	return result.String()
}

func renderText(text string, font map[rune][]string) string {
	var result strings.Builder

	for row := 0; row < 8; row++ {

		for _, char := range text {

			if ascii, ok := font[char]; ok {
				result.WriteString(ascii[row])
			}
		}

		result.WriteString("\n")
	}

	return result.String()
}

func loadBanner(path string) (map[rune][]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	font := make(map[rune][]string)

	if len(lines) > 0 && lines[0] == "" {
		lines = lines[1:]
	}

	index := 0

	for char := 32; char <= 126; char++ {

		if index+8 > len(lines) {
			break
		}

		font[rune(char)] = lines[index : index+8]

		index += 9
	}

	return font, nil
}

func renderTemplate(w http.ResponseWriter, data PageData) {
	templatePath := filepath.Join("templates", "index.html")

	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, "Template execution error", http.StatusInternalServerError)
	}
}

func main() {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/ascii-art", asciiHandler)

	log.Println("Server running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}