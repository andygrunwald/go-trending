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
//
func NewTrending() *Trending {
	baseURL, _ := url.Parse(defaultBaseURL)
	t := Trending{
		BaseURL: baseURL,
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
	doc, err := goquery.NewDocument(u.String())
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
			ContributerURL: contributorURL,
			Contributer:    developer,
		}
		projects = append(projects, p)
	})

	return projects, nil
}

// GetLanguages will return a slice of Language known by gitub.
// With the Language.URLName you can filter your GetProjects / GetDevelopers calls.
func (t *Trending) GetLanguages() ([]Language, error) {
	return t.generateLanguages("div.one-fourth div.select-menu a")
}

// GetTrendingLanguages will return a slice of Language that are currently trending.
// Trending languages are displayed at https://github.com/trending on the right side.
// With the Language.URLName you can filter your GetProjects / GetDevelopers calls.
func (t *Trending) GetTrendingLanguages() ([]Language, error) {
	return t.generateLanguages("ul.language-filter-list a")
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
	doc, err := goquery.NewDocument(u.String())
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
	doc, err := goquery.NewDocument(u.String())
	if err != nil {
		return developers, err
	}

	// Query information
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
