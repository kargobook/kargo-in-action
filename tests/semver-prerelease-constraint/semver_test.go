// Package semverconstraint empirically checks the claim made by Listing 4.4
// in "Kargo in Action" (chapter 4, "A Warehouse that tracks pre-release
// versions"):
//
//	constraint: ">=1.2.0-rc.0 <1.2.0"
//
// with the callout "The -rc.0 lower bound opts pre-release tags into the
// range."
//
// It uses github.com/Masterminds/semver/v3 at the version Kargo depends on
// and reproduces Kargo's own tag-parsing logic (pkg/controller/semver.Parse
// in akuity/kargo, which wraps semver.NewVersion / semver.StrictNewVersion)
// so the tests also demonstrate what a Warehouse image subscription's
// strictSemvers setting actually changes: tag *parsing* eligibility, not
// constraint *evaluation*. See ../README.md for how to run this.
package semverconstraint

import (
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/Masterminds/semver/v3"
)

// bookConstraint is the constraint printed in Listing 4.4.
const bookConstraint = ">=1.2.0-rc.0 <1.2.0"

// kargoParse mirrors pkg/controller/semver.Parse in akuity/kargo. strict=true
// additionally requires the tag to parse via semver.StrictNewVersion (full
// MAJOR.MINOR.PATCH, no missing components) after a leading "v" is trimmed.
func kargoParse(tag string, strict bool) *semver.Version {
	sv, err := semver.NewVersion(tag)
	if err != nil {
		return nil // tag wasn't a semantic version
	}
	if strict {
		if _, err := semver.StrictNewVersion(strings.TrimPrefix(tag, "v")); err != nil {
			return nil // tag wasn't a strict semantic version
		}
	}
	return sv
}

type tagCase struct {
	tag string
	// matches is whether bookConstraint admits the tag once parsed.
	matches bool
	// parsesLoose / parsesStrict are whether kargoParse accepts the tag with
	// strictSemvers false / true.
	parsesLoose  bool
	parsesStrict bool
	why          string
}

// cases covers the tags the book's prose reasons about, plus two that
// exercise the strictSemvers half of the claim rather than asserting it.
var cases = []tagCase{
	{"1.1.9", false, true, true, "below the lower bound"},
	{"1.2.0-rc.1", true, true, true, "the case the book is about"},
	{"1.2.0-rc.2", true, true, true, "the case the book is about"},
	{"1.2.0-beta.3", false, true, true, "beta sorts below rc.0, so it stays out"},
	{"1.2.0", false, true, true, "the final release is excluded by <1.2.0"},
	{"1.2.1", false, true, true, "above the upper bound"},
	{"1.3.0", false, true, true, "above the upper bound"},
	{"1.2", false, true, false, "strictSemvers rejects a missing PATCH component"},
	{"v1.2.0-rc.1", true, true, true, "a leading v is trimmed before the strict check, so strict accepts it"},
}

// TestBookConstraintAdmitsOnlyReleaseCandidates is the claim in Listing 4.4:
// the >=1.2.0-rc.0 lower bound opts 1.2.0-rc.* tags in and nothing else.
func TestBookConstraintAdmitsOnlyReleaseCandidates(t *testing.T) {
	c, err := semver.NewConstraint(bookConstraint)
	if err != nil {
		t.Fatalf("book constraint %q does not parse: %v", bookConstraint, err)
	}
	for _, tc := range cases {
		sv := kargoParse(tc.tag, false)
		if sv == nil {
			t.Errorf("tag %q: kargoParse(strict=false) rejected it, want parsed (%s)", tc.tag, tc.why)
			continue
		}
		if got := c.Check(sv); got != tc.matches {
			t.Errorf("tag %q: constraint %q match = %v, want %v (%s)",
				tc.tag, bookConstraint, got, tc.matches, tc.why)
		}
	}
}

// TestStrictSemversGovernsParsingNotEvaluation demonstrates the second claim:
// strictSemvers decides whether a tag is recognized as a semver at all, and
// changes nothing about how an already-parsed version is checked against a
// constraint.
func TestStrictSemversGovernsParsingNotEvaluation(t *testing.T) {
	c, err := semver.NewConstraint(bookConstraint)
	if err != nil {
		t.Fatalf("book constraint %q does not parse: %v", bookConstraint, err)
	}

	strictRejectedSomething := false
	for _, tc := range cases {
		loose := kargoParse(tc.tag, false)
		strict := kargoParse(tc.tag, true)

		if (loose != nil) != tc.parsesLoose {
			t.Errorf("tag %q: parses with strictSemvers=false = %v, want %v (%s)",
				tc.tag, loose != nil, tc.parsesLoose, tc.why)
		}
		if (strict != nil) != tc.parsesStrict {
			t.Errorf("tag %q: parses with strictSemvers=true = %v, want %v (%s)",
				tc.tag, strict != nil, tc.parsesStrict, tc.why)
		}
		if strict == nil {
			strictRejectedSomething = true
			continue
		}
		// Where both modes parse the tag, the constraint verdict must agree.
		if c.Check(loose) != c.Check(strict) {
			t.Errorf("tag %q: constraint verdict differs between strictSemvers modes (%v vs %v); "+
				"strictSemvers must affect parsing only",
				tc.tag, c.Check(loose), c.Check(strict))
		}
	}
	if !strictRejectedSomething {
		t.Error("no test case was rejected by strictSemvers=true, so this test proves nothing " +
			"about the strictSemvers half of the claim; add a non-strict tag such as \"1.2\"")
	}
}

var semverRequire = regexp.MustCompile(`(?m)^\s*(?:require\s+)?github\.com/Masterminds/semver/v3\s+(v\S+)`)

// TestSemverVersionMatchesKargo pins this module's Masterminds/semver version
// to Kargo's own, so the result above stays a statement about Kargo rather
// than about whatever version happens to be pinned here. CI fetches Kargo's
// go.mod at the release tag under test and passes it in KARGO_GO_MOD.
func TestSemverVersionMatchesKargo(t *testing.T) {
	kargoGoMod := os.Getenv("KARGO_GO_MOD")
	if kargoGoMod == "" {
		t.Skip("KARGO_GO_MOD not set; see ../README.md for how to point this at Kargo's go.mod")
	}
	kargoVersion := requiredSemverVersion(t, kargoGoMod)
	localVersion := requiredSemverVersion(t, "go.mod")
	if kargoVersion != localVersion {
		t.Errorf("Masterminds/semver/v3 is %s here but %s in Kargo's go.mod (%s); "+
			"bump this module so the test still describes Kargo's behavior",
			localVersion, kargoVersion, kargoGoMod)
	}
}

func requiredSemverVersion(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("reading %s: %v", path, err)
	}
	m := semverRequire.FindSubmatch(data)
	if m == nil {
		t.Fatalf("%s has no github.com/Masterminds/semver/v3 requirement", path)
	}
	return string(m[1])
}
