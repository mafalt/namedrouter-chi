## Description

Describe what this pull request changes and why.

## Type of change

- [ ] Bug fix
- [ ] New feature
- [ ] Refactoring
- [ ] Documentation
- [ ] Chi compatibility
- [ ] Other

## Related issue

Closes #

## Changes

Describe the main changes introduced by this pull request.

## Architecture

Does this change belong entirely in the Chi adapter?

If the change affects NamedRouter itself or the common `Adapter` interface,
explain why.

Consider:

- Does the change introduce NamedRouter-specific behavior into the adapter?
- Does it expose Chi-specific functionality through the common API?
- Does it affect route parameter parsing or application?
- Does it affect middleware or subrouter behavior?

## Chi compatibility

If applicable, describe which Chi behavior or API this change relies on.

- Chi version:
- Chi-specific behavior:

## Testing

- [ ] Existing tests pass
- [ ] New tests added
- [ ] Existing tests updated where necessary
- [ ] `go vet ./...` passes
- [ ] Code formatted with `gofmt`

## Documentation

- [ ] Documentation updated
- [ ] Public API documentation updated if necessary
- [ ] No documentation changes required

## Additional notes

Add any additional information that reviewers should know.
