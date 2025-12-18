package dynamic_test

import (
	"testing"

	dynamic "github.com/aura-studio/dynamic"
)

func TestAllowed_Detect(t *testing.T) {
	a := dynamic.NewAllowed()

	cases := []struct {
		in   string
		want dynamic.AllowedType
		ok   bool
	}{
		{"", dynamic.AllowedTypeURL, true},
		{"default", dynamic.AllowedTypeKeyword, true},
		{"aura-studio-dynamic", dynamic.AllowedTypeKeyword, true},
		{"./plugins/libgo_x.so", dynamic.AllowedTypePath, true},
		{"../plugins/libgo_x.so", dynamic.AllowedTypePath, true},
		{"/opt/warehouse/x", dynamic.AllowedTypePath, true},
		{"C:/warehouse/x", dynamic.AllowedTypePath, true},
		{"C:\\warehouse\\x", dynamic.AllowedTypePath, true},
		{"\\\\server\\share\\x", dynamic.AllowedTypePath, true},
		{"https://example.com/a/b", dynamic.AllowedTypeURL, true},
		{"s3://bucket/path", dynamic.AllowedTypeURL, true},
		{"not a value", 0, false},
	}

	for _, tc := range cases {
		got, ok := a.Detect(tc.in)
		if ok != tc.ok {
			t.Fatalf("Detect(%q) ok=%v want %v", tc.in, ok, tc.ok)
		}
		if ok && got != tc.want {
			t.Fatalf("Detect(%q)=%v want %v", tc.in, got, tc.want)
		}
	}
}

func TestAllowed_IsHelpers(t *testing.T) {
	a := dynamic.NewAllowed()
	if a.IsKeyword("") {
		t.Fatalf("IsKeyword(\"\") should be false")
	}
	if !a.IsKeyword("abc-123") {
		t.Fatalf("IsKeyword should accept lowercase/digits/hyphen")
	}
	if !a.IsPath("") {
		t.Fatalf("IsPath(\"\") should be true")
	}
	if !a.IsURL("") {
		t.Fatalf("IsURL(\"\") should be true")
	}
}
