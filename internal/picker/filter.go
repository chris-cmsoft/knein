package picker

import "strings"

// FilterContexts returns up to limit contexts matching all whitespace-separated
// terms in query. Matching is case-insensitive and substring-based.
func FilterContexts(contexts []string, query string, limit int) []string {
	if limit < 1 {
		return nil
	}

	terms := strings.Fields(strings.ToLower(query))
	matches := make([]string, 0, min(limit, len(contexts)))

	for _, context := range contexts {
		if matchesTerms(context, terms) {
			matches = append(matches, context)
			if len(matches) == limit {
				break
			}
		}
	}

	return matches
}

func matchesTerms(context string, terms []string) bool {
	name := strings.ToLower(context)
	for _, term := range terms {
		if !strings.Contains(name, term) {
			return false
		}
	}
	return true
}
