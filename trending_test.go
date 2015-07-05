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
		fmt.Println("Projects (filtered by Go) recieved.")
	} else {
		fmt.Printf("Number of projectes recieved: %d (filtered by golang %v)", len(projects), onlyGoProjects)
	}

	// Output: Projects (filtered by Go) recieved.
}

func ExampleTrending_GetLanguages() {
	trend := trending.NewTrending()
	languages, err := trend.GetLanguages()
	if err != nil {
		log.Fatal(err)
	}

	if len(languages) > 0 {
		fmt.Println("Languages recieved.")
	} else {
		fmt.Printf("Number of languages recieved: %d", len(languages))
	}

	// Output: Languages recieved.
}

func ExampleTrending_GetDevelopers() {
	trend := trending.NewTrending()
	developers, err := trend.GetDevelopers(trending.TimeWeek, "")
	if err != nil {
		log.Fatal(err)
	}

	if len(developers) > 0 {
		fmt.Println("Developers recieved.")
	} else {
		fmt.Printf("Number of developer recieved: %d", len(developers))
	}

	// Output: Developers recieved.
}
