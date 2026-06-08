package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strings"
)

// PageData holds data for the template
type PageData struct {
	Text   string
	Banner string
	Result string
}

// BannerMap stores loaded banner files
var BannerMap = make(map[string][]string)

func main() {
	// Load all banner files at startup
	banners := []string{"standard", "shadow", "thinkertoy"}
	for _, name := range banners {
		if err := loadBanner(name); err != nil {
			fmt.Printf("Error loading %s: %v\n", name, err)
			os.Exit(1)
		}
	}

	// Set up routes
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/ascii-art", asciiArtHandler)

	// Start server
	fmt.Println("Server running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

// loadBanner reads a banner file into memory
func loadBanner(name string) error {
	data, err := os.ReadFile("banners/" + name + ".txt")
	if err != nil {
		return err
	}
	// Split by newline and store
	BannerMap[name] = strings.Split(string(data), "\n")
	return nil
}

// homeHandler serves the main page (GET /)
func homeHandler(w http.ResponseWriter, r *http.Request) {
	// Check for wrong path
	if r.URL.Path != "/" {
		http.Error(w, "404 Not Found", http.StatusNotFound)
		return
	}
	// Only allow GET
	if r.Method != http.MethodGet {
		http.Error(w, "405 Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Load and execute template
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, PageData{})
}

// asciiArtHandler processes the form (POST /ascii-art)
func asciiArtHandler(w http.ResponseWriter, r *http.Request) {
	// Only allow POST
	if r.Method != http.MethodPost {
		http.Error(w, "405 Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse form data
	r.ParseForm()
	text := r.FormValue("text")
	banner := r.FormValue("banner")

	// Validate input
	if text == "" {
		http.Error(w, "400 Bad Request", http.StatusBadRequest)
		return
	}

	// Default to standard if no banner selected
	if banner == "" {
		banner = "standard"
	}

	// Check if banner exists
	bannerLines, ok := BannerMap[banner]
	if !ok {
		http.Error(w, "404 Not Found", http.StatusNotFound)
		return
	}

	// Generate ASCII art using your existing render logic
	result := render(text, bannerLines)

	// Load template and show result on same page
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
		return
	}

	data := PageData{
		Text:   text,
		Banner: banner,
		Result: result,
	}

	tmpl.Execute(w, data)
}

// render is YOUR original function - unchanged logic
func render(input string, banner []string) string {
	if input == "" {
		return ""
	}

	var result strings.Builder
	parts := strings.Split(input, "\\n")

	for i, part := range parts {

		if part == "" {
			if i < len(parts)-1 {
				result.WriteString("\n")
			}
			continue
		}

		for row := 0; row < 8; row++ {
			for _, ch := range part {

				if ch < 32 || ch > 126 {
					ch = ' '
				}

				start := (int(ch)-32)*9 + 1

				if start+row < len(banner) {
					result.WriteString(banner[start+row])
				}
			}
			result.WriteString("\n")
		}
	}

	return result.String()
}
