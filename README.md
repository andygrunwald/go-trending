# go-trending
A package to retrieve [trending repositories](https://github.com/trending) and [developers](https://github.com/trending/developers) from Github written in [golang](https://golang.org/).

* TODO Add Screenshot
* TODO Add TravisCI
* TODO Add GoDoc
* TODO Add 

## Installation

TODO

## API

TODO

## Examples

### List trending repositories of today for all languages

```go
package main

import (
	"fmt"
	"github.com/andygrunwald/go-trending"
	"log"
)

func main() {
	trend := trending.NewTrending()
	
	// Show projects of today
	projects, err := trend.GetProjects(trending.TimeToday, "")
	if err != nil {
		log.Fatal(err)
	}
	for index, project := range projects {
		no := index + 1
		if len(project.Language) > 0 {
			fmt.Printf("%d: %s (written in %s with %d \xE2\xAD\x90 )\n", no, project.Name, project.Language, project.Stars)
		} else {
			fmt.Printf("%d: %s (with %d \xE2\xAD\x90 )\n", no, project.Name, project.Stars)
		}
	}
}
```

### List trending repositories of this week for golang

TODO

### List trending developers of this month for Swift

TODO

### List available languages

```go
package main

import (
	"fmt"
	"github.com/andygrunwald/go-trending"
	"log"
)

func main() {
	trend := trending.NewTrending()

	// Show languages
	languages, err := trend.GetLanguages()
	if err != nil {
		log.Fatal(err)
	}
	for index, language := range languages {
		no := index + 1
		fmt.Printf("%d: %s (%s)\n", no, language.Name, language.URLName)
	}
}

```

## License

This project is released under the terms of the [MIT license](http://en.wikipedia.org/wiki/MIT_License).
