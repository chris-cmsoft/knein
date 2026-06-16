package picker

import (
	"reflect"
	"testing"
)

func TestFilterContextsMatchesAllTerms(t *testing.T) {
	contexts := []string{
		"prd-euw1-main",
		"prd-use1-main",
		"dev-euw1-main",
		"shared-prd-euw1-admin",
	}

	got := FilterContexts(contexts, "prd euw1", 9)
	want := []string{"prd-euw1-main", "shared-prd-euw1-admin"}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("FilterContexts() = %#v, want %#v", got, want)
	}
}

func TestFilterContextsIsCaseInsensitive(t *testing.T) {
	contexts := []string{"PRD-EUW1-MAIN"}

	got := FilterContexts(contexts, "prd euw1", 9)
	want := []string{"PRD-EUW1-MAIN"}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("FilterContexts() = %#v, want %#v", got, want)
	}
}

func TestFilterContextsLimitsMatches(t *testing.T) {
	contexts := []string{"a", "b", "c"}

	got := FilterContexts(contexts, "", 2)
	want := []string{"a", "b"}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("FilterContexts() = %#v, want %#v", got, want)
	}
}
