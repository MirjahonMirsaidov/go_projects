package choose_your_own_adventure

import (
	"encoding/json"
	"html/template"
	"io"
	"log"
	"net/http"
	"strings"
)

type Story map[string]Chapter

type Chapter struct {
	Title      string   `json:"title"`
	Paragraphs []string `json:"story"`
	Options    []Option `json:"options"`
}

type Option struct {
	Text    string `json:"text"`
	Chapter string `json:"arc"`
}

func JsonStory(r io.Reader) (Story, error) {
	d := json.NewDecoder(r)
	var story Story
	if err := d.Decode(&story); err != nil {
		return nil, err
	}
	return story, nil
}

type handler struct {
	story Story
}

func NewHandler(story Story) http.Handler {
	return handler{story}
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimSpace(r.URL.Path)
	if path == "" || path == "/" {
		path = "/intro"
	}
	title := path[1:]
	if destination, ok := h.story[title]; ok {
		tmpl := template.Must(template.ParseFiles("index.html"))
		err := tmpl.Execute(w, destination)
		if err != nil {
			log.Printf("Something went wrong, %v", err)
			http.Error(w, "Something went wrong.", http.StatusBadRequest)
		}
	} else {
		http.Error(w, "Page not found", http.StatusNotFound)
	}

}
