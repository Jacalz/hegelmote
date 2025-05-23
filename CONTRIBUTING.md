Thank you for being interested in contributing to Hegelmote.
Any help we can get is appreciated. Reporting bugs, suggesting improvements or new featuresm contributing code etc. are all of great help.

## Reporting a bug

If you've found an issue with the application, please report it to help us fix it as soon as possible.
When reporting a bug, please follow the guidelines below:

1. Check the [issue list](https://github.com/Jacalz/hegelmote/issues) to see if it's already been reported. If so, update the existing issue with any additional information that you have (if necessary).
2. If not, then create a new issue using the issue template for reporting bugs.
3. Stay involved in the conversation on the issue and answer any questions that might arise. More information can sometimes be necessary.

## Code Contribution

Great! You have either found a bug to fix or a new feature to implement.
Follow the steps below to increase the chance of the changes being accepted quickly.

1. Read and follow the guidelines in the [Code standards](#Code-standards) section further down this page.
2. Consider how to structure your code so that it is readable, clean, and can be easily tested.
4. Write the code changes and create a new commit for your change.
5. Run the tests and make sure everything still works as expected using `go test ./...`.
6. Please refrain from force pushing or squashing. This makes it easier to review, and squashing can instead be done automatically when merging.

### Code standards

We aim to maintain a very high standard of code through design, testing, and implementation.
To manage this, we have various checks and processes in place that everyone should follow, including:

- For a more strict standard Go format, we use [gofumpt](https://github.com/mvdan/gofumpt).
- The code should pass the code quality checks by [staticcheck](https://staticcheck.io/) and [gosec](https://github.com/securego/gosec).
- The cyclomatic complexity of each function should be below 16.
