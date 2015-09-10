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

	gloomysonURL, _ := url.Parse(server.URL + "/gloomyson")
	gloomysonAvatar, _ := url.Parse("https://avatars2.githubusercontent.com/u/13479175?v=3")
	backScreenURL, _ := url.Parse(server.URL + "/black-screen")
	backScreenAvatar, _ := url.Parse("https://avatars0.githubusercontent.com/u/14174343?v=3")
	want := []Developer{
		Developer{
			ID:          13479175,
			DisplayName: "gloomyson",
			FullName:    "Ryuta",
			URL:         gloomysonURL,
			Avatar:      gloomysonAvatar,
		},
		Developer{
			ID:          14174343,
			DisplayName: "black-screen",
			FullName:    "Black Screen",
			URL:         backScreenURL,
			Avatar:      backScreenAvatar,
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

	uAll, _ := url.Parse("https://github.com/trending")
	uUnknown, _ := url.Parse("https://github.com/trending?l=unknown")
	uCSS, _ := url.Parse("https://github.com/trending?l=css")
	uGo, _ := url.Parse("https://github.com/trending?l=go")
	uJava, _ := url.Parse("https://github.com/trending?l=java")
	want := []Language{
		Language{"All languages", "", uAll},
		Language{"Unknown languages", "unknown", uUnknown},
		Language{"CSS", "css", uCSS},
		Language{"Go", "go", uGo},
		Language{"Java", "java", uJava},
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

	uAbap, _ := url.Parse("https://github.com/trending?l=abap")
	uActionScript, _ := url.Parse("https://github.com/trending?l=as3")
	uAda, _ := url.Parse("https://github.com/trending?l=ada")
	uAgda, _ := url.Parse("https://github.com/trending?l=agda")
	uAGS, _ := url.Parse("https://github.com/trending?l=ags-script")
	uAlloy, _ := url.Parse("https://github.com/trending?l=alloy")
	uAMPL, _ := url.Parse("https://github.com/trending?l=ampl")
	uANTLR, _ := url.Parse("https://github.com/trending?l=antlr")

	want := []Language{
		Language{"ABAP", "abap", uAbap},
		Language{"ActionScript", "as3", uActionScript},
		Language{"Ada", "ada", uAda},
		Language{"Agda", "agda", uAgda},
		Language{"AGS Script", "ags-script", uAGS},
		Language{"Alloy", "alloy", uAlloy},
		Language{"AMPL", "ampl", uAMPL},
		Language{"ANTLR", "antlr", uANTLR},
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

func TestGetProjects_NoConent(t *testing.T) {
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
			"since": "monthly",
			"l":     "java",
		})
		website := getContentOfFile("./tests/github.com_trending.html")
		fmt.Fprint(w, string(website))
	})

	projects, err := client.GetProjects(TimeMonth, "java")
	if err != nil {
		t.Errorf("GetProjects returned error: %v", err)
	}

	uStarcraft, _ := url.Parse(server.URL + "/gloomyson/StarCraft")
	uStarcraftContributer, _ := url.Parse(server.URL + "/gloomyson/StarCraft/graphs/contributors")
	uBlackScreen, _ := url.Parse(server.URL + "/black-screen/black-screen")
	uBlackScreenContributer, _ := url.Parse(server.URL + "/black-screen/black-screen/graphs/contributors")
	uNode, _ := url.Parse(server.URL + "/nodejs/node")
	uNodeContributer, _ := url.Parse(server.URL + "/nodejs/node/graphs/contributors")

	gloomysonURL, _ := url.Parse(server.URL + "/gloomyson")
	gloomysonAvatar, _ := url.Parse("https://avatars0.githubusercontent.com/u/13479175?v=3")
	shockoneURL, _ := url.Parse(server.URL + "/shockone")
	shockoneAvatar, _ := url.Parse("https://avatars3.githubusercontent.com/u/188928?v=3")
	g07chaURL, _ := url.Parse(server.URL + "/G07cha")
	g07chaAvatar, _ := url.Parse("https://avatars3.githubusercontent.com/u/6943514?v=3")
	robertoUaURL, _ := url.Parse(server.URL + "/RobertoUa")
	robertoUaAvatar, _ := url.Parse("https://avatars1.githubusercontent.com/u/1307169?v=3")
	ryURL, _ := url.Parse(server.URL + "/ry")
	ryAvatar, _ := url.Parse("https://avatars1.githubusercontent.com/u/80?v=3")
	bnoordhuisURL, _ := url.Parse(server.URL + "/bnoordhuis")
	bnoordhuisAvatar, _ := url.Parse("https://avatars1.githubusercontent.com/u/275871?v=3")

	want := []Project{
		Project{
			Name:           "gloomyson/StarCraft",
			Owner:          "gloomyson",
			RepositoryName: "StarCraft",
			Description:    "HTML5 version for StarCraft game",
			Language:       "JavaScript",
			Stars:          1624,
			URL:            uStarcraft,
			ContributerURL: uStarcraftContributer,
			Contributer: []Developer{
				Developer{
					ID:          13479175,
					DisplayName: "gloomyson",
					FullName:    "",
					URL:         gloomysonURL,
					Avatar:      gloomysonAvatar,
				},
			},
		},
		Project{
			Name:           "black-screen/black-screen",
			Owner:          "black-screen",
			RepositoryName: "black-screen",
			Description:    "A terminal emulator for the 21st century.",
			Language:       "JavaScript",
			Stars:          624,
			URL:            uBlackScreen,
			ContributerURL: uBlackScreenContributer,
			Contributer: []Developer{
				Developer{
					ID:          188928,
					DisplayName: "shockone",
					FullName:    "",
					URL:         shockoneURL,
					Avatar:      shockoneAvatar,
				},
				Developer{
					ID:          6943514,
					DisplayName: "G07cha",
					FullName:    "",
					URL:         g07chaURL,
					Avatar:      g07chaAvatar,
				},
				Developer{
					ID:          1307169,
					DisplayName: "RobertoUa",
					FullName:    "",
					URL:         robertoUaURL,
					Avatar:      robertoUaAvatar,
				},
			},
		},
		Project{
			Name:           "nodejs/node",
			Owner:          "nodejs",
			RepositoryName: "node",
			Description:    "✨Future Node.js releases will be from this repo. ✨",
			Language:       "JavaScript",
			Stars:          348,
			URL:            uNode,
			ContributerURL: uNodeContributer,
			Contributer: []Developer{
				Developer{
					ID:          80,
					DisplayName: "ry",
					FullName:    "",
					URL:         ryURL,
					Avatar:      ryAvatar,
				},
				Developer{
					ID:          275871,
					DisplayName: "bnoordhuis",
					FullName:    "",
					URL:         bnoordhuisURL,
					Avatar:      bnoordhuisAvatar,
				},
			},
		},
	}

	if !reflect.DeepEqual(projects, want) {
		t.Errorf("GetProjects returned %+v, want %+v", projects, want)
	}
}
