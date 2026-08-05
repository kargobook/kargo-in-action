# Tests

Standalone verification programs, kept separate from the chapter examples so they can be re-run as Kargo and its dependencies move.

- [`semver-prerelease-constraint/`](semver-prerelease-constraint) — proves that the ch04 Listing 4.4 constraint `>=1.2.0-rc.0 <1.2.0` admits only `1.2.0-rc.*` pre-release tags (not `1.2.0-beta.*` or the final `1.2.0`), using the exact `Masterminds/semver/v3` version Kargo depends on, and that Kargo's `strictSemvers` setting affects tag *parsing*, not constraint *evaluation*. Run with `cd semver-prerelease-constraint && go run main.go`.
