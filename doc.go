/*
Package trending provides access to github`s trending repositories and developers.
The data will be collected from githubs website at https://github.com/trending and https://github.com/trending/developers.

Construct a new trending client, then use the various functions on the client to
access the trending content from GitHub.
For example:

	trend := trending.NewTrending()

	// Get trending projects of language "go" for today.
	projects, err := trend.GetProjects(trending.TimeToday, "go")

GitHub Enterprise

If you are running a GitHub Enterprise yourself you can use this library as well.
For such cases you can change the BaseURL of the library:

	myURL, _ := url.Parse("https://my.github.enterprise.com")
	trend := trending.NewTrending()
	trend.BaseURL = myURL

	// Get trending projects of language "go" for today.
	projects, err := trend.GetProjects(trending.TimeToday, "go")

*/
package trending
