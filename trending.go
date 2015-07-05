// Package trending provides access to github`s trending repositories and developers.
// The data will be collected from githubs website at https://github.com/trending and https://github.com/trending/developers.
package trending

import (
	"github.com/PuerkitoBio/goquery"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

// NewTrending is the main entry point of the trending package.
// It provides access to the API of this package by returning a Trending datastructure.
// Usage:
//
//		trend := trending.NewTrending()
//		projects, err := trend.GetProjects(trending.TimeToday, "")
//		...
func NewTrending() *Trending {
	t := Trending{}
	return &t
}

// GetProjects provides a slice of Project filtered by the given time and language.
// time can be filtered by applying by one of the Time* constants (e.g. TimeToday, TimeWeek, ...). If an empty string will be applied TimeToday will be the default.
// language can be filtered by applying a programing language by your choice. The input must be a known language by Github and be part of GetLanguages(). Further more it must be the Language.URLName and not the human readable Language.Name.
// If language is an empty string "All languages" will be applied.
func (t *Trending) GetProjects(time, language string) ([]Project, error) {
	// BUG(andygrunwald): time don`t get a default value if you apply an empty string. Default: TimeToday
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
		projectURL := t.appendBaseHostToPath(address, exists)

		description := s.Find(".repo-list-description").Text()
		description = strings.TrimSpace(description)

		meta := s.Find(".repo-list-meta").Text()
		language, stars := t.getLanguageAndStars(meta)

		contributerPath, exists := s.Find(".repo-list-meta a").First().Attr("href")
		contributerURL := t.appendBaseHostToPath(contributerPath, exists)

		p := Project{
			Name:           name,
			Description:    description,
			Language:       language,
			Stars:          stars,
			URL:            projectURL,
			ContributerURL: contributerURL,
		}

		projects = append(projects, p)
	})

	return projects, nil
}

// GetLanguages will return a slice of Language known by gitub.
// With the Language.URLName you can filter your GetProjects / GetDevelopers calls.
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

		filterURL, _ := url.Parse(languageURLName)

		re := regexp.MustCompile("github.com/trending\\?l=(.+)")
		if matches := re.FindStringSubmatch(languageURLName); len(matches) >= 2 && len(matches[1]) > 0 {
			languageURLName = matches[1]
		}

		language := Language{
			Name:    s.Text(),
			URLName: languageURLName,
			URL:     filterURL,
		}

		languages = append(languages, language)
	})

	return languages, nil
}

// GetDevelopers provides a slice of Developer filtered by the given time and language.
// time can be filtered by applying by one of the Time* constants (e.g. TimeToday, TimeWeek, ...). If an empty string will be applied TimeToday will be the default.
// language can be filtered by applying a programing language by your choice. The input must be a known language by Github and be part of GetLanguages(). Further more it must be the Language.URLName and not the human readable Language.Name.
// If language is an empty string "All languages" will be applied.
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
		fullName = t.trimBraces(fullName)

		linkHref, exists := s.Find(".user-leaderboard-list-name a").Attr("href")
		linkURL := t.appendBaseHostToPath(linkHref, exists)

		avatar, exists := s.Find("img.leaderboard-gravatar").Attr("src")
		avatarURL := t.buildAvatarURL(avatar, exists)

		developers = append(developers, t.newDeveloper(name, fullName, linkURL, avatarURL))
	})

	return developers, nil
}

// newDeveloper creates a new Developer
func (t *Trending) newDeveloper(name, fullName string, linkURL, avatarURL *url.URL) Developer {
	return Developer{
		ID:          t.getUserIDBasedOnAvatarURL(avatarURL),
		DisplayName: name,
		FullName:    fullName,
		URL:         linkURL,
		Avatar:      avatarURL,
	}
}

func (t *Trending) trimBraces(text string) string {
	text = strings.TrimSpace(text)
	text = strings.TrimLeft(text, "(")
	text = strings.TrimRight(text, ")")

	return text
}

func (t *Trending) buildAvatarURL(avatar string, exists bool) *url.URL {
	var avatarURL *url.URL
	var err error

	if exists == true {
		avatarURL, err = url.Parse(avatar)
		if err != nil {
			avatarURL = nil
		}
	}

	return avatarURL
}

// getUserIDBasedOnAvatarLink determines the UserID based on an avatar link avatarURL
func (t *Trending) getUserIDBasedOnAvatarURL(avatarURL *url.URL) int {
	id := 0
	if avatarURL != nil {
		re := regexp.MustCompile("u/([0-9]+)")
		if matches := re.FindStringSubmatch(avatarURL.Path); len(matches) >= 2 && len(matches[1]) > 0 {
			id, _ = strconv.Atoi(matches[1])
		}
	}

	return id
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

func (t *Trending) appendBaseHostToPath(address string, exists bool) *url.URL {
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
