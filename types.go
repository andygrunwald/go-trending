package trending

import (
	"net/url"
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

	// ContributerURL is the http(s) address of the contributors page of the project reflected as url.URL datastructure like "https://github.com/Workiva/go-datastructures/graphs/contributors".
	ContributerURL *url.URL

	// Contributer are a collection of Developer.
	// Be aware that this collection don`t covers all contributor.
	// Only those who are mentioned at githubs trending page.
	Contributer []Developer
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
