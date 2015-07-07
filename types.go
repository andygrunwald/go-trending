package trending

import (
	"github.com/PuerkitoBio/goquery"
	"net/url"
)

// These are predefined constants to define the timerange of the requested repository or developer.
// If trending repositories or developer are requested a timeframe has to be added.
// It is suggested to use this constants for this.
// TimeToday is limit of the current day.
// TimeWeek will focus on the complete week
// TimeMonth include the complete month
const (
	TimeToday = "daily"
	TimeWeek  = "weekly"
	TimeMonth = "monthly"
)

// Internal used constants related to github`s website / structure.
const (
	baseHost       = "https://github.com"
	basePath       = "/trending"
	developersPath = "/developers"
)

// Internal used constants to determine the requested resource.
// The trending page of github provides repositories, developers and languages.
// These constants are used (amongst others) to generate the correct url
// to recieve the resources.
// We don`t export these constants, because the trending package provides
// dedicated methods to differ between repositories, developers and languages.
const (
	modeRepositories = "repositories"
	modeDevelopers   = "developers"
	modeLanguages    = "languages"
)

// Trending reflects the main datastructure of this package.
// It doesn`t provide an exported state, but based on this the methods are called.
// To recieve a new instance just add
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
	document *goquery.Document
}

// Project reflects a single trending repository.
// It provides information as printed on the source website https://github.com/trending.
// Name is the name of the repository including user / organisation like "andygrunwald/go-trending" or "airbnb/javascript".
// Description is the description of the repository like "JavaScript Style Guide" (for "airbnb/javascript").
// Language is the determined programing language of the project (by Github). Sometimes Language is an empty string, because Github can`t determine the (main) programing language (like for "google/deepdream").
// Stars is the number of github stars this project recieved in the given timeframe (see TimeToday / TimeWeek / TimeMonth constants). This number don`t reflect the overall stars of the project.
// URL is the http(s) address of the project reflected as url.URL datastructure like "https://github.com/Workiva/go-datastructures".
// ContributerURL is the http(s) address of the contributers page of the project reflected as url.URL datastructure like "https://github.com/Workiva/go-datastructures/graphs/contributors".
// Contributer are a collection of Developer. Be aware that this collection don`t covers all contributer. Only those who are mentioned at githubs trending page.
type Project struct {
	Name           string
	Description    string
	Language       string
	Stars          int
	URL            *url.URL
	ContributerURL *url.URL
	Contributer    []Developer
}

// Language reflects a single (programing) language offered by github for filtering.
// If you call "GetProjects" you are able to filter by programing language.
// For filter input you should use the URLName of Language.
// Name is the human readable name of the language like "Go" or "Web Ontology Language"
// URLName is the machine readable / usable name of the language used for filtering / url parameters like "go" or "web-ontology-language". Please use URLName if you want to filter your requests.
// URL is the filter URL for the language like "https://github.com/trending?l=go" for "go" or "https://github.com/trending?l=unknown" or "unknown".
type Language struct {
	Name, URLName string
	URL           *url.URL
}

// Developer reflects a single trending developer / organisation.
// It provides information as printed on the source website https://github.com/trending/developers.
// ID is the github`s unique identifier of the user / organisation like 1342004 (google) or 698437 (airbnb).
// DisplayName is the username of the developer / organisation like "torvalds" or "apache".
// FullName is the real name of the developer / organisation like "Linus Torvalds" (for "torvalds") or "The Apache Software Foundation" (for "apache").
// URL is the http(s) address of the developer / organisation reflected as url.URL datastructure like https://github.com/torvalds.
// Avatar is the http(s) address of the developer / organisation avatar as url.URL datastructure like https://avatars1.githubusercontent.com/u/1024025?v=3&s=192.
type Developer struct {
	ID          int
	DisplayName string
	FullName    string
	URL         *url.URL
	Avatar      *url.URL
}
