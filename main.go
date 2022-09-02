package main

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
)

const baseUrl = "https://github.com/"

var (
	username = "disco07"
	password = "MAmanko91"
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

// getToken retrieve token on website to login later
// return Token and error
func (a App) getToken() (Token, error) {
	// Request the HTML page.
	res, err := a.Client.Get(baseUrl + "login")
	if err != nil {
		fmt.Println("Error fetching response. ", err)
		return Token{}, err
	}
	defer res.Body.Close()

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		fmt.Println("Error loading HTTP response body. ", err)
		return Token{}, err
	}

	var token Token
	// Catch token from page login
	t, found := doc.Find("input[name='authenticity_token']").Attr("value")
	if !found {
		return Token{}, errors.New("value not found")
	}
	token.Token = t

	return token, nil
}

// login func to connect to a website and return an error.
func (a App) login() error {
	token, err := a.getToken()
	if err != nil {
		return err
	}
	data := url.Values{
		"login":              {username},
		"password":           {password},
		"authenticity_token": {token.Token},
	}
	res, err := a.Client.PostForm(baseUrl+"login", data)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	fmt.Println(res.StatusCode)

	if res.StatusCode != 200 {
		return errors.New("failed to authenticate")
	}

	return nil
}

// GetProjects retrieve a name of projects on website
// return an array Project and an error
func (a App) getProjects() ([]Project, error) {
	// Request the HTML page.
	res, err := a.Client.Get(baseUrl + "disco07?tab=repositories")
	if err != nil {
		fmt.Println("Error fetching response. ", err)
		return nil, err
	}
	defer res.Body.Close()

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		fmt.Println("Error loading HTTP response body. ", err)
		return nil, err
	}

	var projects []Project

	doc.Find("#user-repositories-list li").Each(func(i int, s *goquery.Selection) {
		title := s.Find("h3 a").Text()
		project := Project{
			Name: strings.TrimSpace(title),
		}
		projects = append(projects, project)
	})

	return projects, nil
}

func main() {
	jar, err := cookiejar.New(nil)
	if err != nil {
		log.Fatal(err)
	}
	app := App{&http.Client{Jar: jar}}

	if err := app.login(); err != nil {
		log.Fatal(err)
	}

	projects, err := app.getProjects()
	if err != nil {
		log.Fatal(err)
	}

	for index, project := range projects {
		fmt.Printf("%d: %s\n", index+1, project.Name)
	}
}
