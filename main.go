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
	username = "*****"
	password = "*****"
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

func (a App) GetToken() (Token, error) {
	// Request the HTML page.
	res, err := a.Client.Get(base_url + "/login")
	if err != nil {
		log.Fatalln("Error fetching response. ", err)
		return Token{}, err
	}
	defer res.Body.Close()

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal("Error loading HTTP response body. ", err)
		return Token{}, err
	}

	var token Token
	// Catch token from page login
	t, _ := doc.Find("input[name='authenticity_token']").Attr("value")
	token.Token = t

	return token, nil
}

func (a App) Login() error {
	token, err := a.GetToken()
	if err != nil {
		return err
	}
	data := url.Values{
		"username":           {username},
		"password":           {password},
		"authenticity_token": {token.Token},
	}
	res, err := a.Client.PostForm(base_url+"/login", data)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	_, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	return nil
}

func (a App) GetProjects() ([]Project, error) {
	// Request the HTML page.
	res, err := a.Client.Get(base_url + "/disco07?tab=repositories")
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

	if err := app.Login(); err != nil {
		log.Fatal(err)
	}

	projects, err := app.GetProjects()
	if err != nil {
		log.Fatal(err)
	}

	for index, project := range projects {
		fmt.Printf("%d: %s\n", index+1, project.Name)
	}
}
