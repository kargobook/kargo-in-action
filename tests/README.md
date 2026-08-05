# Tests

Standalone verification programs, kept separate from the chapter examples so they can be re-run as Kargo and its dependencies move. They run in CI (the `semver-claim` job in [`.github/workflows/verify.yml`](../.github/workflows/verify.yml)), so a claim that stops holding turns the build red instead of printing a verdict nobody reads.

## `semver-prerelease-constraint/`

Checks the two claims chapter 4's Listing 4.4 makes about the constraint `>=1.2.0-rc.0 <1.2.0`:

1. The `-rc.0` lower bound admits `1.2.0-rc.*` tags and nothing else: not `1.2.0-beta.*`, not the final `1.2.0`, and nothing outside the `1.2.0` release.
2. A Warehouse image subscription's `strictSemvers` setting governs whether a tag is recognized as a semantic version at all, not how an already-parsed version is checked against a constraint. The tag list includes `1.2`, which `strictSemvers: true` rejects for its missing `PATCH` component, and `v1.2.0-rc.1`, which it *accepts*, because Kargo trims a leading `v` before the strict check.

It uses `github.com/Masterminds/semver/v3` and reproduces Kargo's own `pkg/controller/semver.Parse`, which wraps `semver.NewVersion` and `semver.StrictNewVersion`.

Run it:

```bash
cd semver-prerelease-constraint
go test -v ./...
```

`TestSemverVersionMatchesKargo` keeps the result a statement about Kargo rather than about whatever version happens to be pinned here: it compares this module's `Masterminds/semver/v3` requirement against the one in Kargo's own `go.mod`. Point it at a copy of Kargo's `go.mod` to run it:

```bash
curl -fsSL -o /tmp/kargo-go.mod \
  https://raw.githubusercontent.com/akuity/kargo/v1.11.0/go.mod
KARGO_GO_MOD=/tmp/kargo-go.mod go test -count=1 -v ./...
```

Without `KARGO_GO_MOD` the test skips, so a local `go test` needs no network. CI always sets it, at the release tag the book is written against.
