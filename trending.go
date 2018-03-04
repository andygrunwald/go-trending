package trending

import (
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// These are predefined constants to define the timerange of the requested repository or developer.
// If trending repositories or developer are requested, a timeframe has to be added.
// It is suggested to use this constants for this.
const (
	// TimeToday is limit of the current day.
	TimeToday = "daily"
	// TimeWeek will focus on the complete week
	TimeWeek = "weekly"
	// TimeMonth include the complete month
	TimeMonth = "monthly"

	// Base URL for the github website
	defaultBaseURL = "https://github.com"
	// Relative URL for trending repositories
	urlTrendingPath = "/trending"
	// Relative URL for trending developers
	urlDevelopersPath = "/developers"

	// Standard mode: github.com/trending
	modeRepositories = "repositories"
	// Developers mode: github.com/trending/developers
	modeDevelopers = "developers"
	// Language mode: Only query parameters will be added
	modeLanguages = "languages"
)

// Trending reflects the main datastructure of this package.
// It doesn`t provide an exported state, but based on this the methods are called.
// To receive a new instance just add
//
//		package main
//
//		import (
//			"github.com/andygrunwald/go-trending"
//		)
//
//		func main() {
//			trend := trending.NewTrending()
//			...
//		}
//
type Trending struct {
	// Base URL for requests.
	// Defaults to the public GitHub website, but can be set to a domain endpoint to use with GitHub Enterprise.
	// BaseURL should always be specified with a trailing slash.
	BaseURL *url.URL

	// Client to use for requests
	Client *http.Client
}

// Project reflects a single trending repository.
// It provides information as printed on the source website https://github.com/trending.
type Project struct {
	// Name is the name of the repository including user / organisation like "andygrunwald/go-trending" or "airbnb/javascript".
	Name string

	// Owner is the name of the user or organisation. "andygrunwald" in "andygrunwald/go-trending" or "airbnb" in "airbnb/javascript".
	Owner string

	// RepositoryName is the name of therepository. "go-trending" in "andygrunwald/go-trending" or "javascript" in "airbnb/javascript".
	RepositoryName string

	// Description is the description of the repository like "JavaScript Style Guide" (for "airbnb/javascript").
	Description string

	// Language is the determined programing language of the project (by Github).
	// Sometimes Language is an empty string, because Github can`t determine the (main) programing language (like for "google/deepdream").
	Language string

	// Stars is the number of github stars this project received in the given timeframe (see TimeToday / TimeWeek / TimeMonth constants).
	// This number don`t reflect the overall stars of the project.
	Stars int

	// URL is the http(s) address of the project reflected as url.URL datastructure like "https://github.com/Workiva/go-datastructures".
	URL *url.URL

	// ContributorURL is the http(s) address of the contributors page of the project reflected as url.URL datastructure like "https://github.com/Workiva/go-datastructures/graphs/contributors".
	ContributorURL *url.URL

	// Contributor are a collection of Developer.
	// Be aware that this collection don`t covers all contributor.
	// Only those who are mentioned at githubs trending page.
	Contributor []Developer
}

// Language reflects a single (programing) language offered by github for filtering.
// If you call "GetProjects" you are able to filter by programing language.
// For filter input you should use the URLName of Language.
type Language struct {
	// Name is the human readable name of the language like "Go" or "Web Ontology Language"
	Name string

	// URLName is the machine readable / usable name of the language used for filtering / url parameters like "go" or "web-ontology-language".
	// Please use URLName if you want to filter your requests.
	URLName string

	// URL is the filter URL for the language like "https://github.com/trending?l=go" for "go" or "https://github.com/trending?l=unknown" or "unknown".
	URL *url.URL
}

// Developer reflects a single trending developer / organisation.
// It provides information as printed on the source website https://github.com/trending/developers.
type Developer struct {
	// ID is the github`s unique identifier of the user / organisation like 1342004 (google) or 698437 (airbnb).
	ID int

	// // DisplayName is the username of the developer / organisation like "torvalds" or "apache".
	DisplayName string

	// FullName is the real name of the developer / organisation like "Linus Torvalds" (for "torvalds") or "The Apache Software Foundation" (for "apache").
	FullName string

	// URL is the http(s) address of the developer / organisation reflected as url.URL datastructure like https://github.com/torvalds.
	URL *url.URL

	// Avatar is the http(s) address of the developer / organisation avatar as url.URL datastructure like https://avatars1.githubusercontent.com/u/1024025?v=3&s=192.
	Avatar *url.URL
}

// NewTrending is the main entry point of the trending package.
// It provides access to the API of this package by returning a Trending datastructure.
// Usage:
//
//		trend := trending.NewTrending()
//		projects, err := trend.GetProjects(trending.TimeToday, "")
//		...
//
func NewTrending() *Trending {
	return NewTrendingWithClient(http.DefaultClient)
}

// NewTrendingWithClient allows providing a custom http.Client to use for fetching trending items.
// It allows setting timeouts or using 3rd party http.Client implementations, such as Google App Engine
// urlfetch.Client.
func NewTrendingWithClient(client *http.Client) *Trending {
	baseURL, _ := url.Parse(defaultBaseURL)
	t := Trending{
		BaseURL: baseURL,
		Client:  client,
	}
	return &t
}

// GetProjects provides a slice of Projects filtered by the given time and language.
//
// time can be filtered by applying by one of the Time* constants (e.g. TimeToday, TimeWeek, ...).
// If an empty string will be applied TimeToday will be the default (current default by Github).
//
// language can be filtered by applying a programing language by your choice.
// The input must be a known language by Github and be part of GetLanguages().
// Further more it must be the Language.URLName and not the human readable Language.Name.
// If language is an empty string "All languages" will be applied (current default by Github).
func (t *Trending) GetProjects(time, language string) ([]Project, error) {
	var projects []Project

	// Generate the correct URL to call
	u, err := t.generateURL(modeRepositories, time, language)
	if err != nil {
		return projects, err
	}

	// Receive document
	res, err := t.Client.Get(u.String())
	if err != nil {
		return projects, err
	}

	doc, err := goquery.NewDocumentFromResponse(res)
	if err != nil {
		return projects, err
	}

	// Query our information
	doc.Find("ol.repo-list li").Each(func(i int, s *goquery.Selection) {

		// Collect project information
		name := t.getProjectName(s.Find("h3 a").Text())

		// Split name (like "andygrunwald/go-trending") into owner ("andygrunwald") and repository name ("go-trending"")
		splittedName := strings.SplitAfterN(name, "/", 2)
		owner := splittedName[0][:len(splittedName[0])-1]
		owner = strings.TrimSpace(owner)
		repositoryName := strings.TrimSpace(splittedName[1])

		address, exists := s.Find("h3 a").First().Attr("href")
		projectURL := t.appendBaseHostToPath(address, exists)

		description := s.Find(".py-1 p").Text()
		description = strings.TrimSpace(description)

		language := s.Find("div.f6 span").Eq(1).Text()
		language = strings.TrimSpace(language)

		starsString := s.Find("div.f6 a").First().Text()
		starsString = strings.TrimSpace(starsString)
		starsString = strings.Replace(starsString, ",", "", 1)
		starsString = strings.Replace(starsString, ".", "", 1)
		stars, err := strconv.Atoi(starsString)
		if err != nil {
			stars = 0
		}

		contributerSelection := s.Find("div.f6 a").Eq(2)
		contributorPath, exists := contributerSelection.Attr("href")
		contributorURL := t.appendBaseHostToPath(contributorPath, exists)

		// Collect contributor
		var developer []Developer
		contributerSelection.Find("img").Each(func(j int, devSelection *goquery.Selection) {
			devName, exists := devSelection.Attr("title")
			linkURL := t.appendBaseHostToPath(devName, exists)

			avatar, exists := devSelection.Attr("src")
			avatarURL := t.buildAvatarURL(avatar, exists)

			developer = append(developer, t.newDeveloper(devName, "", linkURL, avatarURL))
		})

		p := Project{
			Name:           name,
			Owner:          owner,
			RepositoryName: repositoryName,
			Description:    description,
			Language:       language,
			Stars:          stars,
			URL:            projectURL,
			ContributorURL: contributorURL,
			Contributor:    developer,
		}
		projects = append(projects, p)
	})

	return projects, nil
}

// GetLanguages will return a slice of Language known by gitub.
// With the Language.URLName you can filter your GetProjects / GetDevelopers calls.
func (t *Trending) GetLanguages() ([]Language, error) {
	return t.generateLanguages(".col-md-3 div.select-menu .select-menu-list a")
}

// GetTrendingLanguages will return a slice of Language that are currently trending.
// Trending languages are displayed at https://github.com/trending on the right side.
// With the Language.URLName you can filter your GetProjects / GetDevelopers calls.
func (t *Trending) GetTrendingLanguages() ([]Language, error) {
	return t.generateLanguages("ul.filter-list a")
}

// generateLanguages will retrieve the languages out of the github document.
// Trending languages are shown on the right side as a small list.
// Other languages are hidden in a dropdown at this site
func (t *Trending) generateLanguages(mainSelector string) ([]Language, error) {
	var languages []Language

	// Generate the URL to call
	u, err := t.generateURL(modeLanguages, "", "")
	if err != nil {
		return languages, err
	}

	// Get document
	res, err := t.Client.Get(u.String())
	if err != nil {
		return languages, err
	}

	doc, err := goquery.NewDocumentFromResponse(res)
	if err != nil {
		return languages, err
	}

	// Query our information
	doc.Find(mainSelector).Each(func(i int, s *goquery.Selection) {
		languageAddress, _ := s.Attr("href")
		languageURLName := ""

		filterURL, _ := url.Parse(languageAddress)

		re := regexp.MustCompile("github.com/trending/([^/\\?]*)")
		if matches := re.FindStringSubmatch(languageAddress); len(matches) >= 2 && len(matches[1]) > 0 {
			languageURLName = matches[1]
		}

		language := Language{
			Name:    strings.TrimSpace(s.Text()),
			URLName: languageURLName,
			URL:     filterURL,
		}
		languages = append(languages, language)
	})

	return languages, nil
}

// GetDevelopers provides a slice of Developer filtered by the given time and language.
//
// time can be filtered by applying by one of the Time* constants (e.g. TimeToday, TimeWeek, ...).
// If an empty string will be applied TimeToday will be the default (current default by Github).
//
// language can be filtered by applying a programing language by your choice
// The input must be a known language by Github and be part of GetLanguages().
// Further more it must be the Language.URLName and not the human readable Language.Name.
// If language is an empty string "All languages" will be applied (current default by Github).
func (t *Trending) GetDevelopers(time, language string) ([]Developer, error) {
	var developers []Developer

	// Generate URL
	u, err := t.generateURL(modeDevelopers, time, language)
	if err != nil {
		return developers, err
	}

	// Get document
	res, err := t.Client.Get(u.String())
	if err != nil {
		return developers, err
	}

	doc, err := goquery.NewDocumentFromResponse(res)
	if err != nil {
		return developers, err
	}

	// Query information
	doc.Find(".explore-content li").Each(func(i int, s *goquery.Selection) {
		name := s.Find("h2 a").Text()
		name = strings.TrimSpace(name)
		name = strings.Split(name, " ")[0]
		name = strings.TrimSpace(name)

		fullName := s.Find("h2 a span").Text()
		fullName = t.trimBraces(fullName)

		linkHref, exists := s.Find("h2 a").Attr("href")
		linkURL := t.appendBaseHostToPath(linkHref, exists)

		avatar, exists := s.Find("a img").Attr("src")
		avatarURL := t.buildAvatarURL(avatar, exists)

		developers = append(developers, t.newDeveloper(name, fullName, linkURL, avatarURL))
	})

	return developers, nil
}

// newDeveloper is a utility function to create a new Developer
func (t *Trending) newDeveloper(name, fullName string, linkURL, avatarURL *url.URL) Developer {
	return Developer{
		ID:          t.getUserIDBasedOnAvatarURL(avatarURL),
		DisplayName: name,
		FullName:    fullName,
		URL:         linkURL,
		Avatar:      avatarURL,
	}
}

// trimBraces will remove braces "(" & ")" from the string
func (t *Trending) trimBraces(text string) string {
	text = strings.TrimSpace(text)
	text = strings.TrimLeft(text, "(")
	text = strings.TrimRight(text, ")")
	return text
}

// buildAvatarURL will build a url.URL out of the Avatar URL provided by Github
func (t *Trending) buildAvatarURL(avatar string, exists bool) *url.URL {
	if exists == false {
		return nil
	}

	avatarURL, err := url.Parse(avatar)
	if err != nil {
		return nil
	}

	// Remove s parameter
	// The "s" parameter controls the size of the avatar
	q := avatarURL.Query()
	q.Del("s")
	avatarURL.RawQuery = q.Encode()

	return avatarURL
}

// getUserIDBasedOnAvatarLink determines the UserID based on an avatar link avatarURL
func (t *Trending) getUserIDBasedOnAvatarURL(avatarURL *url.URL) int {
	id := 0
	if avatarURL == nil {
		return id
	}

	re := regexp.MustCompile("u/([0-9]+)")
	if matches := re.FindStringSubmatch(avatarURL.Path); len(matches) >= 2 && len(matches[1]) > 0 {
		id, _ = strconv.Atoi(matches[1])
	}
	return id
}

// appendBaseHostToPath will add the base host to a relative url urlStr.
//
// A urlStr like "/trending" will be returned as https://github.com/trending
func (t *Trending) appendBaseHostToPath(urlStr string, exists bool) *url.URL {
	if exists == false {
		return nil
	}

	rel, err := url.Parse(urlStr)
	if err != nil {
		return nil
	}

	return t.BaseURL.ResolveReference(rel)
}

// getProjectName will return the project name in format owner/repository
func (t *Trending) getProjectName(name string) string {
	trimmedNameParts := []string{}

	nameParts := strings.Split(name, "\n")
	for _, part := range nameParts {
		trimmedNameParts = append(trimmedNameParts, strings.TrimSpace(part))
	}

	return strings.Join(trimmedNameParts, "")
}

// generateURL will generate the correct URL to call the github site.
//
// Depending on mode, time and language it will set the correct pathes and query parameters.
func (t *Trending) generateURL(mode, time, language string) (*url.URL, error) {
	urlStr := urlTrendingPath
	if mode == modeDevelopers {
		urlStr += urlDevelopersPath
	}

	u := t.appendBaseHostToPath(urlStr, true)
	q := u.Query()
	if len(time) > 0 {
		q.Set("since", time)
	}

	if len(language) > 0 {
		q.Set("l", language)
	}

	u.RawQuery = q.Encode()

	return u, nil
}
