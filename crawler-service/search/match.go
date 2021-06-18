package search

type Matcher interface {
	Search() ([]string, error)
}

func Match(matcher Matcher) ([]string, error) {
	return matcher.Search()
}
