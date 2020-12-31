# How to contribute

Thank you for your interest in contribution to go-financial.
Here's the recommended process of contribution.



If you've changed APIs, update the documentation.
Make sure your code lints.
Issue that pull request!
1. Fork the repo and create your branch from master.
1. hack, hack, hack.
1. If you've added code that should be tested, add tests.
1. If the functions are exported, make sure they are documented with examples.
1. Update README if required.
1. Make sure that unit tests are passing `make test-unit`.
1. Make sure that lint-checks are also passing `make lint-check`.  
1. Make sure that your change follows best practices in Go
   - [Effective Go](https://golang.org/doc/effective_go.html)
   - [Go Code Review Comments](https://golang.org/wiki/CodeReviewComments)
1. Open a pull request in GitHub in this repo.

When you work on a larger contribution, it is also recommended that you get in touch
with us through the issue tracker.

## Code reviews

All submissions, including submissions by project members, require review.

## Releases

Releases are partially automated. To draft a new release, follow these steps:

1. Decide on a release version.
1. Create a new release/{version} branch and edit the `CHANGELOG.md`.
1. Make sure that the tests are passing.
1. Get the PR merged.
1. Tag the merge commit (`vX.Y.Z`) and create a release on GitHub for the tag. (repo maintainers)
1. (Required) Sit back and pat yourself on the back for a job well done :clap:.
