// Command semver-prerelease-constraint empirically checks the claim made by
// Listing 4.4 in "Kargo in Action" (chapter 4, "A Warehouse that tracks
// pre-release versions"):
//
//	constraint: ">=1.2.0-rc.0 <1.2.0"
//
// with the callout "The -rc.0 lower bound opts pre-release tags into the
// range."
//
// It uses github.com/Masterminds/semver/v3 at the exact version Kargo v1.11.0
// depends on (v3.4.0, per Kargo's go.mod), and reproduces Kargo's own
// tag-parsing logic (pkg/controller/semver.Parse in akuity/kargo at tag
// v1.11.0, which wraps semver.NewVersion / semver.StrictNewVersion) so the
// test also shows what the Warehouse image subscription's strictSemvers
// setting actually changes: tag *parsing* eligibility, not constraint
// *evaluation*. See ../README.md for how to run this and what it proves.
package main

import (
	"fmt"
	"strings"

	"github.com/Masterminds/semver/v3"
)

// kargoParse mirrors pkg/controller/semver.Parse in akuity/kargo (tag
// v1.11.0). strict=true additionally requires the tag to parse via
// semver.StrictNewVersion (full MAJOR.MINOR.PATCH, no missing components).
func kargoParse(tag string, strict bool) *semver.Version {
	sv, err := semver.NewVersion(tag)
	if err != nil {
		return nil
	}
	if strict {
		if _, err := semver.StrictNewVersion(strings.TrimPrefix(tag, "v")); err != nil {
			return nil
		}
	}
	return sv
}

func main() {
	const bookConstraint = ">=1.2.0-rc.0 <1.2.0"

	// The book's claim, tag by tag: which of these does the constraint admit?
	// (Expectations are what the surrounding prose in ch04 implies: rc tags
	// for the same 1.2.0 release are opted in, everything else stays out.)
	expected := map[string]bool{
		"1.1.9":        false,
		"1.2.0-rc.1":   true,
		"1.2.0-rc.2":   true,
		"1.2.0-beta.3": false,
		"1.2.0":        false,
		"1.2.1":        false,
		"1.3.0":        false,
	}
	tags := []string{"1.1.9", "1.2.0-rc.1", "1.2.0-rc.2", "1.2.0-beta.3", "1.2.0", "1.2.1", "1.3.0"}

	fmt.Println("=== Masterminds/semver v3.4.0 (Kargo v1.11.0's exact go.mod dependency) ===")
	fmt.Printf("Constraint under test (book Listing 4.4): %q\n\n", bookConstraint)

	c, err := semver.NewConstraint(bookConstraint)
	if err != nil {
		panic(err)
	}

	allOK := true
	fmt.Println("--- constraint.Check() vs. the book's implied expectation ---")
	for _, t := range tags {
		v, err := semver.NewVersion(t)
		if err != nil {
			fmt.Printf("  %-14s -> PARSE ERROR: %v\n", t, err)
			allOK = false
			continue
		}
		got := c.Check(v)
		ok := got == expected[t]
		allOK = allOK && ok
		status := "OK"
		if !ok {
			status = "MISMATCH"
		}
		fmt.Printf("  %-14s -> match=%-5v expected=%-5v [%s]\n", t, got, expected[t], status)
	}

	fmt.Println()
	fmt.Println("--- Same check, through Kargo's own tag-parse-then-constraint-check path ---")
	fmt.Println("    (kargoParse + constraint.Check, for strictSemvers=false and strictSemvers=true)")
	strictChangesResult := false
	for _, strict := range []bool{false, true} {
		fmt.Printf("\n  strictSemvers=%v:\n", strict)
		for _, t := range tags {
			sv := kargoParse(t, strict)
			if sv == nil {
				fmt.Printf("    %-14s -> REJECTED at parse stage (not a valid tag under strict=%v)\n", t, strict)
				strictChangesResult = true
				continue
			}
			fmt.Printf("    %-14s -> parsed OK, constraint match=%v\n", t, c.Check(sv))
		}
	}

	fmt.Println()
	if allOK {
		fmt.Println("VERDICT: the book's constraint and callout are ACCURATE for all 7 tags.")
	} else {
		fmt.Println("VERDICT: the book's constraint or callout is INACCURATE for at least one tag above.")
	}
	if strictChangesResult {
		fmt.Println("VERDICT: strictSemvers changed the parse-stage outcome for at least one tag above.")
	} else {
		fmt.Println("VERDICT: strictSemvers did NOT change the outcome for any of these 7 tags.")
		fmt.Println("         (All 7 are canonical MAJOR.MINOR.PATCH[-prerelease] strings, so strict")
		fmt.Println("         and non-strict parsing agree. strictSemvers governs whether a tag is")
		fmt.Println("         recognized as a semver at all -- e.g. it would reject a bare \"1.2\" or a")
		fmt.Println("         \"v\"-prefixed tag under strict mode -- it does not change how an already-")
		fmt.Println("         parsed version is checked against a constraint.)")
	}
}
