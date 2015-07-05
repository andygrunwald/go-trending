# go-trending

[![GoDoc](https://godoc.org/github.com/andygrunwald/go-trending?status.svg)](https://godoc.org/github.com/andygrunwald/go-trending)
[![Build Status](https://travis-ci.org/andygrunwald/go-trending.svg?branch=master)](https://travis-ci.org/andygrunwald/go-trending)

A package to retrieve [trending repositories](https://github.com/trending) and [developers](https://github.com/trending/developers) from Github written in [golang](https://golang.org/).

* TODO Add Screenshot

This package were inspired by [rochefort/git-trend](https://github.com/rochefort/git-trend) and [sheharyarn/github-trending](https://github.com/sheharyarn/github-trending).

## Features

repositories
developers
time filtering (day, week, month)
language filtering
TODO

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
	projects, err := trend.GetProjects(trending.TimeWeek, "go")
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

### List trending developers of this month for Swift

```go
package main

import (
	"fmt"
	"github.com/andygrunwald/go-trending"
	"log"
)

func main() {
	trend := trending.NewTrending()

	developers, err := trend.GetDevelopers(trending.TimeMonth, "swift")
	if err != nil {
		log.Fatal(err)
	}
	for index, developer := range developers {
		no := index + 1
		fmt.Printf("%d: %s (%s)\n", no, developer.DisplayName, developer.FullName)
	}
}
```

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

## TODO-List

* Languages url name is the full url. Change it.

## License

This project is released under the terms of the [MIT license](http://en.wikipedia.org/wiki/MIT_License).
