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

	greetingsMu.RLock()
	entries := make([]Greeting, len(greetings))
	copy(entries, greetings)
	greetingsMu.RUnlock()

	for i, j := 0, len(entries)-1; i < j; i, j = i+1, j-1 {
		entries[i], entries[j] = entries[j], entries[i]
	}

	if err := guestbookTemplate.Execute(w, entries); err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

var guestbookTemplate = template.Must(template.New("book").Parse(`
<html>
  <head>
    <title>Modern Go Guestbook</title>
    <link rel="stylesheet" href="/static/style.css">
  </head>
  <body>
    <h1>Go Guestbook</h1>

    <form action="/sign" method="post">
      <div><label for="author">Name (optional):</label></div>
      <div><input type="text" name="author" id="author" placeholder="Your Name"></div>
      <div><label for="content">Message:</label></div>
      <div><textarea name="content" id="content" rows="3" cols="60" required placeholder="Leave a message..."></textarea></div>
      <div><input type="submit" value="Sign Guestbook"></div>
    </form>

    <hr>

    {{range .}}
      <div class="entry">
        {{with .Author}}
          <p><b>{{.}}</b> wrote:</p>
        {{else}}
          <p><b>An anonymous person</b> wrote:</p>
        {{end}}
        <pre>{{.Content}}</pre>
        <small>{{.Date.Format "Jan 2, 2006 15:04:05"}}</small>
      </div>
    {{else}}
        <p>No entries yet. Be the first!</p>
    {{end}}

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
