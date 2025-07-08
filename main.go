package main

import (
	"database/sql"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"os"

	"github.com/lib/pq"
)

type Database struct {
	db *sql.DB
}

type Tag struct {
	TagName string
}

type App struct {
	database  *Database
	templates *template.Template
}

func initApp(dbPath string) (*App, error) {
	db, err := sql.Open("postgres", dbPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error opening database connection!")
		return nil, err
	}

	templates := template.Must(template.ParseFiles("tmpl/add.html", "tmpl/view.html", "tmpl/home.html"))
	app := &App{database: &Database{db}, templates: templates}

	return app, db.Ping()
}

func (app *App) renderTemplate(w http.ResponseWriter, tmpl string, n any) {
	err := app.templates.ExecuteTemplate(w, tmpl+".html", n)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (tag *Tag) save(db *sql.DB) error {
	_, err := db.Exec("INSERT INTO tags VALUES ($1)", tag.TagName)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" {
				fmt.Fprintln(os.Stderr, "Tag already exists")
				return nil
			} else {
				fmt.Fprintf(os.Stderr, "Error saving tag: %s\n", pqErr.Code)
				return err
			}
		}
	}

	return err
}

func (app *App) loadTags(w http.ResponseWriter) ([]Tag, error) {
	rows, err := app.database.db.Query("SELECT tag_name FROM tags")
	if err != nil {
		http.Error(w, "Database query failed", http.StatusInternalServerError)
		return nil, err
	}
	defer rows.Close()

	tags := []Tag{}
	for rows.Next() {
		var tag Tag
		rows.Scan(&tag.TagName)
		tags = append(tags, tag)
	}

	return tags, nil
}

func (app *App) homeHandler(w http.ResponseWriter, r *http.Request) {
	app.renderTemplate(w, "home", nil)
}

func (app *App) viewHandler(w http.ResponseWriter, r *http.Request) {
	tags, err := app.loadTags(w)
	if err != nil {
		http.Error(w, "Failed to load tags", http.StatusInternalServerError)
		return
	}

	app.renderTemplate(w, "view", tags)
}

func (app *App) addHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodGet:
		app.renderTemplate(w, "add", nil)

	case http.MethodPost:
		message := r.FormValue("body")
		n := &Tag{TagName: message}
		err := n.save(app.database.db)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/view", http.StatusFound)

	default:
		http.Error(w, "Unsupported Method Request", http.StatusMethodNotAllowed)
	}
}

func main() {
	app, err := initApp("user=postgres password=super dbname=booru sslmode=disable") // grab this from config file later
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error starting application: %s\n", err.Error())
		os.Exit(1)
	}
	defer app.database.db.Close()

	mux := http.NewServeMux()
	mux.HandleFunc("/", app.homeHandler)
	mux.HandleFunc("/view", app.viewHandler)
	mux.HandleFunc("/add", app.addHandler)

	err = http.ListenAndServe(":8080", mux)
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server %s\n", err)
		os.Exit(1)
	}
}
