# go-trending
A package to retrieve [trending repositories](https://github.com/trending) and [developers](https://github.com/trending/developers) from Github written in [golang](https://golang.org/).

## Installation

TODO

## API

TODO

## Examples

### List trending repositories of today for all languages

TODO

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
		fmt.Printf("%d: %s (%s)\n", index, language.Name, language.URLName)
	}
}

```

## License

This project is released under the terms of the [MIT license](http://en.wikipedia.org/wiki/MIT_License).
