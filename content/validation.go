package content

import (
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/polyglottis/platform/i18n"
)

func ValidExtractType(t ExtractType) bool {
	for _, candidate := range AllExtractTypes {
		if t == candidate {
			return true
		}
	}
	return false
}

func ValidFlavorType(t FlavorType) bool {
	return t == Audio || t == Text || t == Transcript
}

var validSlug = regexp.MustCompile(`^[A-Za-z0-9_]*$`)

// ValidSlug returns true if slug has at least 5 characters and matches ^[A-Za-z0-9_]*$.
// Otherwise an (error) message is returned.
func ValidSlug(slug string) (bool, i18n.Key) {
	if len(slug) < 5 {
		return false, "Slug too short."
	}
	if !validSlug.MatchString(slug) {
		return false, "Only unaccented letters, numbers and underscores are allowed."
	}
	return true, ""
}

var space = regexp.MustCompile(`\s+`)

// NormalizeSlug attempts to create a valid slug based on the seed.
// There is no guarantee that the result is a valid slug.
func NormalizeSlug(seed string) string {
	seed = strings.TrimSpace(seed)
	return space.ReplaceAllString(seed, "_")
}

func ValidLanguageComment(comment string) (bool, i18n.Key) {
	commentLength := utf8.RuneCountInString(comment)
	if commentLength < 5 {
		return false, "Please enter a longer comment."
	}
	if commentLength > 40 {
		return false, "This comment is too long (maximum 40 characters)."
	}
	return true, ""
}

func ValidSummary(summary string) (bool, i18n.Key) {
	summaryLength := utf8.RuneCountInString(summary)
	if summaryLength < 10 {
		return false, "Please enter a longer summary."
	}
	if summaryLength > 150 {
		return false, "This summary is too long (maximum 150 characters)."
	}
	return true, ""
}

func ValidTitle(title string) (bool, i18n.Key) {
	if title == "" {
		return false, "Please enter a title."
	}
	return true, ""
}
