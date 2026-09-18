## Development

Run `gofmt`, `go test ./...`, `go test -race ./...`, and `go vet ./...` before
opening a change. Add regression tests for every behavior change.

## Release

Update the changelog and create an annotated semantic-version tag:

```bash
git tag -a vX.Y.Z -m "release vX.Y.Z"
git push origin vX.Y.Z
```
