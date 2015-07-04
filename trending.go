package trending

import (
	"github.com/PuerkitoBio/goquery"
	"net/url"
	"strconv"
	"strings"
)

const (
	TimeToday = "daily"
	TimeWeek  = "weekly"
	TimeMonth = "monthly"

	baseHost = "https://github.com"
	basePath = "/trending"

	modeRepositories = "repositories"
	modeDevelopers   = "developers"
	modeLanguages    = "languages"
)

type Trending struct {
	document *goquery.Document
}

type Project struct {
	Name        string
	Description string
	Language    string
	Stars       int
	URL         *url.URL
}

type Language struct {
	Name, URLName string
}

type Developer struct {
	DisplayName string
	FullName    string
	URL         *url.URL
	Avatar      *url.URL
}

func NewTrending() *Trending {
	t := Trending{}
	return &t
}

func (t *Trending) GetProjects(time, language string) ([]Project, error) {
	var projects []Project

	u, err := t.generateURL(modeRepositories, time, language)
	if err != nil {
		return projects, err
	}

	doc, err := goquery.NewDocument(u.String())
	if err != nil {
		return projects, err
	}

	doc.Find(".repo-list-item").Each(func(i int, s *goquery.Selection) {

		name := t.getProjectName(s.Find(".repo-list-name a").Text())

		address, exists := s.Find(".repo-list-name a").First().Attr("href")
		projectURL := t.getProjectURL(address, exists)

		description := s.Find(".repo-list-description").Text()
		description = strings.TrimSpace(description)

		meta := s.Find(".repo-list-meta").Text()
		language, stars := t.getLanguageAndStars(meta)

		p := Project{
			Name:        name,
			Description: description,
			Language:    language,
			Stars:       stars,
			URL:         projectURL,
		}

		projects = append(projects, p)
	})

	return projects, nil
}

func (t *Trending) GetLanguages() ([]Language, error) {
	var languages []Language

	u, err := t.generateURL(modeLanguages, "", "")
	if err != nil {
		return languages, err
	}

	doc, err := goquery.NewDocument(u.String())
	if err != nil {
		return languages, err
	}

	doc.Find("div.select-menu-item a").Each(func(i int, s *goquery.Selection) {
		languageURLName, exists := s.Attr("href")
		if exists == false {
			languageURLName = ""
		}

		// TODO
		// language = href.match(/github.com\/trending\?l=(.+)/).to_a[1]
		//      languages << CGI.unescape(language) if language

		language := Language{
			Name:    s.Text(),
			URLName: languageURLName,
		}

		languages = append(languages, language)
	})

	return languages, nil
}

func (t *Trending) GetDevelopers(time, language string) ([]Developer, error) {
	var developers []Developer

	u, err := t.generateURL(modeDevelopers, time, language)
	if err != nil {
		return developers, err
	}

	doc, err := goquery.NewDocument(u.String())
	if err != nil {
		return developers, err
	}

	doc.Find(".user-leaderboard-list-item").Each(func(i int, s *goquery.Selection) {

		name := s.Find(".user-leaderboard-list-name a").Text()
		name = strings.TrimSpace(name)
		name = strings.Split(name, " ")[0]
		name = strings.TrimSpace(name)

		fullName := s.Find(".user-leaderboard-list-name .full-name").Text()
		fullName = strings.TrimSpace(fullName)
		fullName = strings.TrimLeft(fullName, "(")
		fullName = strings.TrimRight(fullName, ")")

		linkHref, exists := s.Find(".user-leaderboard-list-name a").Attr("href")
		var linkURL *url.URL
		if exists == true {
			linkURL, err = url.Parse(baseHost + linkHref)
			if err != nil {
				linkURL = nil
			}
		}

		avatar, exists := s.Find("img.leaderboard-gravatar").Attr("src")
		var avatarURL *url.URL

		if exists == true {
			avatarURL, err = url.Parse(avatar)
			if err != nil {
				avatarURL = nil
			}
		}

		d := Developer{
			DisplayName: name,
			FullName:    fullName,
			URL:         linkURL,
			Avatar:      avatarURL,
		}

		developers = append(developers, d)
	})

	return developers, nil
}

func (t *Trending) getLanguageAndStars(meta string) (string, int) {
	splittedMetaData := strings.Split(meta, string('•'))
	language := ""
	starsIndex := 1

	// If we got 2 parts we only got "stars" and "Built by", but no language
	if len(splittedMetaData) == 2 {
		starsIndex = 0
	} else {
		language = strings.TrimSpace(splittedMetaData[0])
	}

	stars := strings.TrimSpace(splittedMetaData[starsIndex])
	// "stars" contain now a string like
	// 105 stars today
	// 1,472 stars this week
	// 2,552 stars this month
	stars = strings.SplitN(stars, " ", 2)[0]
	stars = strings.Replace(stars, ",", "", 1)
	stars = strings.Replace(stars, ".", "", 1)

	starsInt, err := strconv.Atoi(stars)
	if err != nil {
		starsInt = 0
	}

	return language, starsInt
}

func (t *Trending) getProjectURL(address string, exists bool) *url.URL {
	if exists == false {
		return nil
	}

	u, err := url.Parse(baseHost)
	if err != nil {
		return nil
	}

	u.Path = address

	return u
}

func (t *Trending) getProjectName(name string) string {
	trimmedNameParts := []string{}

	nameParts := strings.Split(name, "\n")
	for _, part := range nameParts {
		trimmedNameParts = append(trimmedNameParts, strings.TrimSpace(part))
	}

	return strings.Join(trimmedNameParts, "")
}

func (t *Trending) generateURL(mode, time, language string) (*url.URL, error) {
	parseURL := baseHost + basePath
	if mode == modeDevelopers {
		parseURL += "/" + modeDevelopers
	}

	u, err := url.Parse(parseURL)
	if err != nil {
		return nil, err
	}

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
