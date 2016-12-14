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

	if c.BaseURL == nil {
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
		website := getContentOfFile("./tests/github.com_trending_developers.html")
		fmt.Fprint(w, string(website))
	})

	developers, err := client.GetDevelopers(TimeToday, "")
	if err != nil {
		t.Errorf("GetDevelopers returned error: %v", err)
	}

	cloudsonURL, _ := url.Parse(server.URL + "/cloudson")
	cloudsonAvatar, _ := url.Parse("https://avatars1.githubusercontent.com/u/94096?v=3")
	zeitURL, _ := url.Parse(server.URL + "/zeit")
	zeitAvatar, _ := url.Parse("https://avatars3.githubusercontent.com/u/14985020?v=3")
	want := []Developer{
		{
			ID:          94096,
			DisplayName: "cloudson",
			FullName:    "Claudson Oliveira",
			URL:         cloudsonURL,
			Avatar:      cloudsonAvatar,
		},
		{
			ID:          14985020,
			DisplayName: "zeit",
			FullName:    "ZEIT",
			URL:         zeitURL,
			Avatar:      zeitAvatar,
		},
	}

	if !reflect.DeepEqual(developers, want) {
		t.Errorf("GetDevelopers returned %+v, want %+v", developers, want)
	}
}

func TestGetDevelopers_NoConent(t *testing.T) {
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

func TestGetTrendingLanguages(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/trending", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		website := getContentOfFile("./tests/github.com_trending.html")
		fmt.Fprint(w, string(website))
	})

	languages, err := client.GetTrendingLanguages()
	if err != nil {
		t.Errorf("GetLanguages returned error: %v", err)
	}

	uAll, _ := url.Parse("https://github.com/trending?since=daily")
	uUnknown, _ := url.Parse("https://github.com/trending/unknown?since=daily")
	uPHP, _ := url.Parse("https://github.com/trending/php?since=daily")
	uGo, _ := url.Parse("https://github.com/trending/go?since=daily")
	uJava, _ := url.Parse("https://github.com/trending/java?since=daily")
	uJavaScript, _ := url.Parse("https://github.com/trending/javascript?since=daily")
	want := []Language{
		{"All languages", "", uAll},
		{"Unknown languages", "unknown", uUnknown},
		{"Go", "go", uGo},
		{"Java", "java", uJava},
		{"JavaScript", "javascript", uJavaScript},
		{"PHP", "php", uPHP},
	}

	if !reflect.DeepEqual(languages, want) {
		t.Errorf("GetLanguages returned %+v, want %+v", languages, want)
	}
}

func TestGetTrendingLanguages_NoConent(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/trending", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
	})

	languages, err := client.GetTrendingLanguages()
	if err != nil {
		t.Errorf("GetLanguages returned error: %v", err)
	}

	var want []Language
	if !reflect.DeepEqual(languages, want) {
		t.Errorf("GetLanguages returned %+v, want %+v", languages, want)
	}
}

func TestGetLanguages(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/trending", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		website := getContentOfFile("./tests/github.com_trending.html")
		fmt.Fprint(w, string(website))
	})

	languages, err := client.GetLanguages()
	if err != nil {
		t.Errorf("GetLanguages returned error: %v", err)
	}

	uAbap, _ := url.Parse("https://github.com/trending/abap?since=daily")
	uActionScript, _ := url.Parse("https://github.com/trending/as3?since=daily")
	uAda, _ := url.Parse("https://github.com/trending/ada?since=daily")
	uAgda, _ := url.Parse("https://github.com/trending/agda?since=daily")
	uAGS, _ := url.Parse("https://github.com/trending/ags-script?since=daily")
	uAlloy, _ := url.Parse("https://github.com/trending/alloy?since=daily")
	uAMPL, _ := url.Parse("https://github.com/trending/ampl?since=daily")
	uANTLR, _ := url.Parse("https://github.com/trending/antlr?since=daily")

	want := []Language{
		{"ABAP", "abap", uAbap},
		{"ActionScript", "as3", uActionScript},
		{"Ada", "ada", uAda},
		{"Agda", "agda", uAgda},
		{"AGS Script", "ags-script", uAGS},
		{"Alloy", "alloy", uAlloy},
		{"AMPL", "ampl", uAMPL},
		{"ANTLR", "antlr", uANTLR},
	}

	if !reflect.DeepEqual(languages, want) {
		t.Errorf("GetLanguages returned %+v, want %+v", languages, want)
	}
}

func TestGetLanguages_NoConent(t *testing.T) {
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
		website := getContentOfFile("./tests/github.com_trending.html")
		fmt.Fprint(w, string(website))
	})

	projects, err := client.GetProjects(TimeToday, "go")
	if err != nil {
		t.Errorf("GetProjects returned error: %v", err)
	}

	uGitql, _ := url.Parse(server.URL + "/cloudson/gitql")
	uGitqlContributor, _ := url.Parse(server.URL + "/cloudson/gitql/graphs/contributors")
	uNeveragaindottech, _ := url.Parse(server.URL + "/neveragaindottech/neveragaindottech.github.io")
	uNeveragaindottechContributor, _ := url.Parse(server.URL + "/neveragaindottech/neveragaindottech.github.io/graphs/contributors")

	cloudsonURL, _ := url.Parse(server.URL + "/cloudson")
	cloudsonAvatar, _ := url.Parse("https://avatars2.githubusercontent.com/u/94096?v=3")
	emhoracekURL, _ := url.Parse(server.URL + "/emhoracek")
	emhoracekAvatar, _ := url.Parse("https://avatars3.githubusercontent.com/u/5353499?v=3")
	zestypingURL, _ := url.Parse(server.URL + "/zestyping")
	zestypingAvatar, _ := url.Parse("https://avatars1.githubusercontent.com/u/236086?v=3")
	saboURL, _ := url.Parse(server.URL + "/sabo")
	saboAvatar, _ := url.Parse("https://avatars0.githubusercontent.com/u/1196568?v=3")

	want := []Project{
		{
			Name:           "cloudson / gitql",
			Owner:          "cloudson",
			RepositoryName: "gitql",
			Description:    "A git query language",
			Language:       "Go",
			Stars:          1824,
			URL:            uGitql,
			ContributerURL: uGitqlContributor,
			Contributer: []Developer{
				{
					ID:          94096,
					DisplayName: "cloudson",
					FullName:    "",
					URL:         cloudsonURL,
					Avatar:      cloudsonAvatar,
				},
			},
		},
		{
			Name:           "neveragaindottech / neveragaindottech.github.io",
			Owner:          "neveragaindottech",
			RepositoryName: "neveragaindottech.github.io",
			Description:    "Source files for the neveragain.tech site",
			Language:       "HTML",
			Stars:          238,
			URL:            uNeveragaindottech,
			ContributerURL: uNeveragaindottechContributor,
			Contributer: []Developer{
				{
					ID:          5353499,
					DisplayName: "emhoracek",
					FullName:    "",
					URL:         emhoracekURL,
					Avatar:      emhoracekAvatar,
				},
				{
					ID:          236086,
					DisplayName: "zestyping",
					FullName:    "",
					URL:         zestypingURL,
					Avatar:      zestypingAvatar,
				},
				{
					ID:          1196568,
					DisplayName: "sabo",
					FullName:    "",
					URL:         saboURL,
					Avatar:      saboAvatar,
				},
			},
		},
	}

	if !reflect.DeepEqual(projects, want) {
		t.Errorf("GetProjects returned %+v, want %+v", projects, want)
	}
}
