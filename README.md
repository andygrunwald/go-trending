# go-trending

[![GoDoc](https://godoc.org/github.com/andygrunwald/go-trending?status.svg)](https://godoc.org/github.com/andygrunwald/go-trending)
[![Go Report Card](https://goreportcard.com/badge/github.com/andygrunwald/go-trending)](https://goreportcard.com/report/github.com/andygrunwald/go-trending)

A package to retrieve [trending repositories](https://github.com/trending) and [developers](https://github.com/trending/developers) from Github written in [Go](https://go.dev/).

[![trending package showcase](./img/go-trending-shrinked.png "trending package showcase")](https://raw.githubusercontent.com/andygrunwald/go-trending/master/img/go-trending-shrinked.png)

## Features

* Get trending repositories
* Get trending developers
* Get all programming languages known by GitHub
* Filtering by time and (programming) language
* Support for [GitHub Enterprise](https://enterprise.github.com/)

## Installation

It is go gettable

    $ go get github.com/andygrunwald/go-trending

or using/updating to the latest master

	$ go get -u github.com/andygrunwald/go-trending@master

## API

Please have a look at the [package documentation](https://pkg.go.dev/github.com/andygrunwald/go-trending) for a detailed API description.

## Examples

A few examples how the API can be used.
More examples are available in the [GoDoc examples section](https://pkg.go.dev/github.com/andygrunwald/go-trending#readme-examples).

### List trending repositories of today for all languages

```go
package main

import (
	"fmt"

	"github.com/andygrunwald/go-trending"
)

func main() {
	trend := trending.NewTrending()

	// Show projects of today
	projects, err := trend.GetProjects(trending.TimeToday, "")
	if err != nil {
		panic(err)
	}

	for index, project := range projects {
		i := index + 1
		if len(project.Language) > 0 {
			fmt.Printf("%d: %s (written in %s with %d ★ )\n", i, project.Name, project.Language, project.Stars)
		} else {
			fmt.Printf("%d: %s (with %d ★ )\n", i, project.Name, project.Stars)
		}
	}
}
```

### List trending repositories of this week for Go

```go
package main

import (
	"fmt"

	"github.com/andygrunwald/go-trending"
)

func main() {
	trend := trending.NewTrending()

	// Show projects of today
	projects, err := trend.GetProjects(trending.TimeWeek, "go")
	if err != nil {
		panic(err)
	}

	for index, project := range projects {
		i := index + 1
		if len(project.Language) > 0 {
			fmt.Printf("%d: %s (written in %s with %d ★ )\n", i, project.Name, project.Language, project.Stars)
		} else {
			fmt.Printf("%d: %s (with %d ★ )\n", i, project.Name, project.Stars)
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
)

func main() {
	trend := trending.NewTrending()

	developers, err := trend.GetDevelopers(trending.TimeMonth, "swift")
	if err != nil {
		panic(err)
	}

	for index, developer := range developers {
		i := index + 1
		fmt.Printf("%d: %s (%s)\n", i, developer.DisplayName, developer.FullName)
	}
}
```

### List available languages

```go
package main

import (
	"fmt"

	"github.com/andygrunwald/go-trending"
)

func main() {
	trend := trending.NewTrending()

	// Show languages
	languages, err := trend.GetLanguages()
	if err != nil {
		panic(err)
	}

	for index, language := range languages {
		i := index + 1
		fmt.Printf("%d: %s (%s)\n", i, language.Name, language.URLName)
	}
}
```

## Implementations

* [sikang99/hub-trend](https://github.com/sikang99/hub-trend/)
* [andygrunwald/TrendingGithub](https://github.com/andygrunwald/TrendingGithub) - [@TrendingGithub](https://twitter.com/TrendingGithub)

## Inspired by

* [rochefort/git-trend](https://github.com/rochefort/git-trend) (Ruby)
* [sheharyarn/github-trending](https://github.com/sheharyarn/github-trending) (Ruby)

## License

This project is released under the terms of the [MIT license](http://en.wikipedia.org/wiki/MIT_License).
