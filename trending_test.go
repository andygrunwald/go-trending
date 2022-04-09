package trending

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"
)

var (
	// mux is the HTTP request multiplexer used with the test server.
	mux *http.ServeMux

	// client is the Trending client being tested.
	client *Trending

	// server is a test HTTP server used to provide mock API responses.
	server *httptest.Server
)

// setup sets up a test HTTP server along with a trending.Trending that is configured to talk to that test server.
// Tests should register handlers on mux which provide mock responses for the http call being tested.
func setup() {
	// test server
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)

	// trending client configured to use test server
	client = NewTrending()
	url, _ := url.Parse(server.URL)
	client.BaseURL = url
}

// teardown closes the test HTTP server.
func teardown() {
	server.Close()
}

// testMethod is a utility function to test the request method provided in want
func testMethod(t *testing.T, r *http.Request, want string) {
	if got := r.Method; got != want {
		t.Errorf("Request method: %v, want %v", got, want)
	}
}

// getContentOfFile is a utility method to open and return the content of fileName
func getContentOfFile(fileName string) []byte {
	content, err := ioutil.ReadFile(fileName)
	if err != nil {
		return []byte{}
	}

	return content
}

type values map[string]string

// testFormValues is a utility method to test the query values given in values
func testFormValues(t *testing.T, r *http.Request, values values) {
	want := url.Values{}
	for k, v := range values {
		want.Add(k, v)
	}

	r.ParseForm()
	if got := r.Form; !reflect.DeepEqual(got, want) {
		t.Errorf("Request parameters: %v, want %v", got, want)
	}
}

func TestNewTrending(t *testing.T) {
	c := NewTrending()
	if c == nil {
		t.Error("Trending client is nil. Expected trending.Trending structure")
	}

	if c != nil && c.BaseURL == nil {
		t.Error("Trending BaseURL is nil. Expected a URL")
	}
}

func TestGetDevelopers_Today(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/trending/developers", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{
			"since": "daily",
		})
		c := getContentOfFile("./testdata/github.com_trending_developers.html")
		fmt.Fprint(w, string(c))
	})

	developers, err := client.GetDevelopers(TimeToday, "")
	if err != nil {
		t.Errorf("GetDevelopers returned error: %v", err)
	}

	n := len(developers)
	if n == 0 {
		t.Error("GetDevelopers returned no developers at all")
	}

	if n <= 25 {
		t.Errorf("GetDevelopers returned %+v developers, expected > 25", n)
	}
}

func TestGetDevelopers_TodayCorrectContent(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/trending/developers", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{
			"since": "daily",
		})
		website := getContentOfFile("./testdata/github.com_trending_developers.html")
		fmt.Fprint(w, string(website))
	})

	developers, err := client.GetDevelopers(TimeToday, "")
	if err != nil {
		t.Errorf("GetDevelopers returned error: %v", err)
	}

	d := developers[0]
	if d.ID == 0 {
		t.Error("GetDevelopers returned no developer ID")
	}
	if len(d.DisplayName) == 0 {
		t.Error("GetDevelopers returned no developer DisplayName")
	}
	if len(d.URL.String()) == 0 {
		t.Error("GetDevelopers returned no developer URL")
	}
	if len(d.Avatar.String()) == 0 {
		t.Error("GetDevelopers returned no developer avatar URL")
	}
}

func TestGetDevelopers_NoContent(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/trending/developers", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{
			"since": "weekly",
			"l":     "go",
		})
	})

	developers, err := client.GetDevelopers(TimeWeek, "go")
	if err != nil {
		t.Errorf("GetDevelopers returned error: %v", err)
	}

	var want []Developer
	if !reflect.DeepEqual(developers, want) {
		t.Errorf("GetDevelopers returned %+v, want %+v", developers, want)
	}
}

func TestGetLanguages_NumberOfLanguages(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/trending", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		website := getContentOfFile("./testdata/github.com_trending.html")
		fmt.Fprint(w, string(website))
	})

	languages, err := client.GetLanguages()
	if err != nil {
		t.Errorf("GetLanguages returned error: %v", err)
	}

	// https://github.com/trending has multiple dropdowns on the page
	// one for the languages (mostly on the right side) and
	// one for the timeframe (today, this week, ...)
	// Here we check if we don't catch the timeframe one
	if languages[0].Name == "today" {
		t.Errorf("GetLanguages catches the timeframe dropdown on https://github.com/trending")
	}

	// Lets simple count the number of language that we got
	// Right now (2019-05-11) the https://github.com/trending
	// has 503 languages. We might use a different number
	// below to be save
	n := len(languages)
	nExpected := 450
	if n == 0 || n < nExpected {
		t.Errorf("GetLanguages returned to less languages (%+v), we expected more than %+v", n, nExpected)
	}
}

func TestGetLanguages_CorrectContent(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/trending", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		website := getContentOfFile("./testdata/github.com_trending.html")
		fmt.Fprint(w, string(website))
	})

	languages, err := client.GetLanguages()
	if err != nil {
		t.Errorf("GetLanguages returned error: %v", err)
	}

	// Today (2022-04-09), we have 637 languages
	// 500 is a random number, that (i assume) will not drop that fast.
	if len(languages) <= 500 {
		t.Errorf("GetLanguages returned %+v languages, expected > 500", len(languages))
	}
	// Might be dirty, but hey ...
	// a) it works
	// b) how high is the chance that HTML is not the 2nd language here?
	// -> Very high :D (until the next testdata update)
	secondLanguage := languages[1]
	expectedLanguage := "HTML"
	if secondLanguage.Name != expectedLanguage {
		t.Errorf("GetLanguages returned %+v, want %+v", secondLanguage.Name, expectedLanguage)
	}

	secondLanguageURL := "https://github.com/trending/html?since=daily"
	if languages[1].URL.String() != secondLanguageURL {
		t.Errorf("GetLanguages returned %+v, want %+v", languages[1].URL.String(), secondLanguageURL)
	}
}

func TestGetLanguages_NoContent(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/trending", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
	})

	languages, err := client.GetLanguages()
	if err != nil {
		t.Errorf("GetLanguages returned error: %v", err)
	}

	var want []Language
	if !reflect.DeepEqual(languages, want) {
		t.Errorf("GetLanguages returned %+v, want %+v", languages, want)
	}
}

func TestGetProjects_NoContent(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/trending", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{
			"since": "monthly",
		})
	})

	projects, err := client.GetProjects(TimeMonth, "")
	if err != nil {
		t.Errorf("GetProjects returned error: %v", err)
	}

	var want []Project
	if !reflect.DeepEqual(projects, want) {
		t.Errorf("GetProjects returned %+v, want %+v", projects, want)
	}
}

func TestGetProjects(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/trending", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{
			"since": "daily",
			"l":     "go",
		})
		website := getContentOfFile("./testdata/github.com_trending.html")
		fmt.Fprint(w, string(website))
	})

	projects, err := client.GetProjects(TimeToday, "go")
	if err != nil {
		t.Errorf("GetProjects returned error: %v", err)
	}

	// Lets simple count the number of language that we got
	// Right now (2019-05-11) the https://github.com/trending
	// has 503 languages. We might use a different number
	// below to be save
	n := len(projects)
	nExpected := 25
	if n == 0 || n < nExpected {
		t.Errorf("GetProjects returned to less projects (%+v), we expected %+v projects", n, nExpected)
	}
}

func TestGetProjects_CorrectContent(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/trending", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{
			"since": "daily",
			"l":     "go",
		})
		website := getContentOfFile("./testdata/github.com_trending.html")
		fmt.Fprint(w, string(website))
	})

	projects, err := client.GetProjects(TimeToday, "go")
	if err != nil {
		t.Errorf("GetProjects returned error: %v", err)
	}

	p := projects[0]
	if len(p.Name) == 0 {
		t.Error("GetProjects returns an empty project name.")
	}
	if len(p.Owner) == 0 {
		t.Error("GetProjects returns an empty project owner.")
	}
	if len(p.RepositoryName) == 0 {
		t.Error("GetProjects returns an empty repository name.")
	}
	if len(p.Language) == 0 {
		t.Error("GetProjects returns an empty language.")
	}
	if p.Stars == 0 {
		t.Error("GetProjects returns a trending project without stars.")
	}
	if len(p.URL.String()) == 0 {
		t.Error("GetProjects returns an empty project URL.")
	}
	if len(p.ContributorURL.String()) == 0 {
		t.Error("GetProjects returns an empty contributor URL.")
	}

	if len(p.Contributor[0].DisplayName) == 0 {
		t.Error("GetProjects returns an empty contributor.")
	}
}
