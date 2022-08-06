package main

import (
	"github.com/PuerkitoBio/goquery"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
)

const base_url = "https://gitlab.com"

var (
	username = "dkone"
	password = "MAmanko91**"
)

type App struct {
	Client *http.Client
}

type Token struct {
	Token string
}

type Project struct {
	Name string
}

func (a App) getToken() Token {
	// Request the HTML page.
	res, err := a.Client.Get("https://steering-tools-v3.novovitae.fr/login")
	if err != nil {
		log.Fatalln("Error fetching response. ", err)
	}
	defer func(Body io.ReadCloser) {
		if err := Body.Close(); err != nil {
			log.Fatal(err)
		}
	}(res.Body)

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal("Error loading HTTP response body. ", err)
	}

	// Catch token from page login
	token, _ := doc.Find("input[name='_csrf_token']").Attr("value")

	return Token{token}
}

func main() {
	jar, _ := cookiejar.New(nil)
	app := App{&http.Client{Jar: jar}}

	app.getToken()
}
