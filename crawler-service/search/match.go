package search

type Matcher interface {
	Search(url string) ([]string, error)
}

func Match(matcher Matcher, url string) ([]string, error) {
	return matcher.Search(url)
}
