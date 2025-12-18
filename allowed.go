package dynamic

import "regexp"

type AllowedType int

const (
	AllowedTypeKeyword AllowedType = iota
	AllowedTypePath
	AllowedTypeURL
)

var allowedRe = map[AllowedType]string{
	AllowedTypeKeyword: `^[a-z0-9][a-z0-9-]*$`,
	// Path supports Unix absolute (/opt/x), Unix relative (./x, ../x),
	// Windows drive paths (C:\x or C:/x), and UNC paths (\\server\share).
	AllowedTypePath: `^(?:(?:(?:[A-Za-z]:[\\/])|(?:\\\\)|/|\./|\.\./).+)?$`,
	// URL matches common scheme URLs like https://, s3://, file://, etc.
	AllowedTypeURL: `^(?:[A-Za-z][A-Za-z0-9+.-]*://\S+)?$`,
}

var allowedReCompiled = map[AllowedType]*regexp.Regexp{
	AllowedTypeKeyword: regexp.MustCompile(allowedRe[AllowedTypeKeyword]),
	AllowedTypePath:    regexp.MustCompile(allowedRe[AllowedTypePath]),
	AllowedTypeURL:     regexp.MustCompile(allowedRe[AllowedTypeURL]),
}

type Allowed struct{}

var allowed = NewAllowed()

func NewAllowed() *Allowed {
	return &Allowed{}
}

func (a *Allowed) Match(t AllowedType, s string) bool {
	re, ok := allowedReCompiled[t]
	if !ok {
		return false
	}
	return re.MatchString(s)
}

func (a *Allowed) IsKeyword(s string) bool {
	return a.Match(AllowedTypeKeyword, s)
}

func (a *Allowed) IsPath(s string) bool {
	return a.Match(AllowedTypePath, s)
}

func (a *Allowed) IsURL(s string) bool {
	return a.Match(AllowedTypeURL, s)
}

// Detect returns the first matched AllowedType in the order: URL -> Path -> Keyword.
func (a *Allowed) Detect(s string) (AllowedType, bool) {
	if a.Match(AllowedTypeURL, s) {
		return AllowedTypeURL, true
	}
	if a.Match(AllowedTypePath, s) {
		return AllowedTypePath, true
	}
	if a.Match(AllowedTypeKeyword, s) {
		return AllowedTypeKeyword, true
	}
	return 0, false
}
