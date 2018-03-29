package trending_test

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/andygrunwald/go-trending"
)

func ExampleTrending_GetProjects() {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	trend := trending.NewTrendingWithClient(client)
	projects, err := trend.GetProjects(trending.TimeToday, "go")
	if err != nil {
		log.Fatal(err)
	}

	onlyGoProjects := true
	for _, project := range projects {
		if len(project.Language) > 0 && project.Language != "Go" {
			onlyGoProjects = false
		}
	}

	if len(projects) > 0 && onlyGoProjects == true {
		fmt.Println("Projects (filtered by Go) received.")
	} else {
		fmt.Printf("Number of projectes received: %d (filtered by golang %v)", len(projects), onlyGoProjects)
	}

	// Output: Projects (filtered by Go) received.
}

func ExampleTrending_GetLanguages() {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	trend := trending.NewTrendingWithClient(client)
	languages, err := trend.GetLanguages()
	if err != nil {
		log.Fatal(err)
	}

	// We need more as 15 languages, because 9 are trending languages
	if len(languages) > 15 {
		fmt.Println("Languages received.")
	} else {
		fmt.Printf("Number of languages received: %d", len(languages))
	}

	// Output: Languages received.
}

func ExampleTrending_GetDevelopers() {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	trend := trending.NewTrendingWithClient(client)
	developers, err := trend.GetDevelopers(trending.TimeToday, "")
	if err != nil {
		log.Fatal(err)
	}

	if len(developers) > 0 {
		fmt.Println("Developers received.")
	} else {
		fmt.Printf("Number of developer received: %d", len(developers))
	}

	// Output: Developers received.
}

func ExampleTrending_GetTrendingLanguages() {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	trend := trending.NewTrendingWithClient(client)
	languages, err := trend.GetTrendingLanguages()
	if err != nil {
		log.Fatal(err)
	}

	if len(languages) > 0 {
		fmt.Println("Trending Languages received.")
	} else {
		fmt.Printf("Number of languages received: %d", len(languages))
	}

	// Output: Trending Languages received.
}
