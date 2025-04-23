package main

import (
	"html/template"
	"log"
	"net/http"
	"sync"
	"time"
)

type Greeting struct {
	Author  string
	Content string
	Date    time.Time
}

var (
	greetings   []Greeting
	greetingsMu sync.RWMutex
)

func root(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	// Ensure no other goroutine is modifying the greetings slice while we read it
	// This is a read lock, allowing multiple readers but blocking writers
	greetingsMu.RLock()
	entries := make([]Greeting, len(greetings))
	copy(entries, greetings)
	greetingsMu.RUnlock()

	// Most recent entries first - reverse the slice
	for i, j := 0, len(entries)-1; i < j; i, j = i+1, j-1 {
		entries[i], entries[j] = entries[j], entries[i]
	}

	// render the template with the entries
	if err := guestbookTemplate.Execute(w, entries); err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

var guestbookTemplate = template.Must(template.New("book").Parse(`
<html>
  <head>
    <title>Go Guestbook 2025</title>
    <link rel="stylesheet" href="/static/style.css">
    <script src="/static/guestbook.js"></script>
  </head>
  <body>
    <h1>Go Guestbook</h1>

    <form action="/sign" method="post" onsubmit="handleFormSubmit(event)">
      <div><label for="author">Name (optional):</label></div>
      <div><input type="text" name="author" id="author" placeholder="Your Name"></div>
      <div><label for="content">Message:</label></div>
      <div><textarea name="content" id="content" rows="3" cols="60" required placeholder="Leave a message..."></textarea></div>
      <div><input type="submit" value="Sign Guestbook"></div>
    </form>

    <hr>

    <div id="entries"></div>

  </body>
</html>
`))

func sign(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	content := r.FormValue("content")
	if content == "" {
		http.Error(w, "Content cannot be empty", http.StatusBadRequest)
		return
	}

	g := Greeting{
		Author:  r.FormValue("author"),
		Content: content,
		Date:    time.Now(),
	}
	if g.Author == "" {
		g.Author = "Anonymous"
	}

	greetingsMu.Lock()
	greetings = append(greetings, g)
	greetingsMu.Unlock()

	http.Redirect(w, r, "/", http.StatusFound)
}

func main() {
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", root)
	http.HandleFunc("/sign", sign)

	log.Println("Starting server on http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
