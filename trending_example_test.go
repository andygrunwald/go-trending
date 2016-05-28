package trending_test

import (
	"fmt"
	"github.com/andygrunwald/go-trending"
	"log"
)

func ExampleTrending_GetProjects() {
	trend := trending.NewTrending()
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
	trend := trending.NewTrending()
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
	trend := trending.NewTrending()
	developers, err := trend.GetDevelopers(trending.TimeWeek, "")
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
	trend := trending.NewTrending()
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
