package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
)

const base_url = "https://github.com/"

var (
	username = "****"
	password = "****"
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

func (a App) GetToken() Token {
	// Request the HTML page.
	res, err := a.Client.Get(base_url + "/login")
	if err != nil {
		log.Fatalln("Error fetching response. ", err)
	}
	defer res.Body.Close()

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal("Error loading HTTP response body. ", err)
	}

	// Catch token from page login
	token, _ := doc.Find("input[name='authenticity_token']").Attr("value")

	return Token{token}
}

func (a App) Login() {
	token := a.GetToken()
	data := url.Values{
		"username":    {username},
		"password":    {password},
		"_csrf_token": {token.Token},
	}
	res, err := a.Client.PostForm(base_url+"/login", data)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	_, err = ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatalln(err)
	}
}

func (a App) GetProjects() []Project {
	// Request the HTML page.
	res, err := a.Client.Get(base_url + "/disco07?tab=repositories")
	if err != nil {
		log.Fatalln("Error fetching response. ", err)
	}
	defer res.Body.Close()

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal("Error loading HTTP response body. ", err)
	}

	var projects []Project

	doc.Find(".js-user-profile-bio").Each(func(i int, s *goquery.Selection) {
		name := strings.TrimSpace(s.Text())
		project := Project{
			Name: name,
		}
		projects = append(projects, project)
	})

	return projects
}

func main() {
	jar, _ := cookiejar.New(nil)
	app := App{&http.Client{Jar: jar}}

	app.Login()
	projects := app.GetProjects()

	for index, project := range projects {
		fmt.Printf("%d: %s\n", index+1, project.Name)
	}
}
