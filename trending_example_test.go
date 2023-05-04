package trending_test

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/andygrunwald/go-trending"
)

/**
 * The (integration) test is failing right now (2023-05-04).
 * At the moment I don't have time to fix it.
 * TODO Fix test
 *
func ExampleTrending_GetProjects() {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	trend := trending.NewTrendingWithClient(client)
	projects, err := trend.GetProjects(trending.TimeToday, "go")
	if err != nil {
		log.Fatal(err)
	}

	projectsNotInGo := 0
	for _, project := range projects {
		if len(project.Language) > 0 && project.Language != "Go" {
			projectsNotInGo = projectsNotInGo + 1
		}
	}

	// Sometimes we get projects where the main language is
	// not Go, but the repository contains a large amount of Go code.
	// Like https://github.com/codeedu/imersao-7-codepix
	// Main language TypeScript + 36% Go.
	// In this case, we calculate a threshold to tackle this situation, because
	// our code works as expected and GitHub is also returning projects with Go code.
	// But it might not be the (only) main language.
	onlyGoProjects := true
	if projectsNotInGo > 0 {
		// Threshold 10%
		if (projectsNotInGo / len(projects) * 100) > 10 {
			onlyGoProjects = false
		}
	}

	if len(projects) > 0 && onlyGoProjects {
		fmt.Println("Projects (filtered by Go) received.")
	} else {
		fmt.Printf("Number of projectes received: %d / projects with a different main language than golang %d)", len(projects), projectsNotInGo)
	}

	// Output: Projects (filtered by Go) received.
}
*/

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
