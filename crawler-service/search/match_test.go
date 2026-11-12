package search

import (
	"errors"
	"reflect"
	"testing"

	"cache_service/crawler-service/handler"
)

type fakeMatcher struct {
	result handler.ReturnList
	calls  int
}

func (f *fakeMatcher) Search() handler.ReturnList {
	f.calls++
	return f.result
}

func TestMatchReturnsMatcherResult(t *testing.T) {
	want := handler.ReturnList{List: []string{"https://example.com/a", "https://example.com/b"}}
	matcher := &fakeMatcher{result: want}

	got := Match(matcher)

	if matcher.calls != 1 {
		t.Errorf("Search called %d times, want 1", matcher.calls)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Match() = %v, want %v", got, want)
	}
}

func TestMatchPropagatesError(t *testing.T) {
	wantErr := errors.New("boom")
	got := Match(&fakeMatcher{result: handler.ReturnList{Error: wantErr}})

	if got.Error != wantErr {
		t.Errorf("Match() error = %v, want %v", got.Error, wantErr)
	}
	if len(got.List) != 0 {
		t.Errorf("Match() list = %v, want empty", got.List)
	}
}
