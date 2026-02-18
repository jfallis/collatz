# Unified Go Guidelines (NASA-Style)

## 1. Keep the control flow simple
- Use only `if`, `for`, `switch`.
- Avoid `goto`; forbid recursion (acyclic call graph).
- Keep `select` small and predictable (few cases, no nesting).
- Prefer clarity over cleverness; explain *why*, not *what*.

---

## 2. Bound all loops and goroutines
- Every loop must have a provable upper bound.
- Service loops must include documented shutdown (e.g., `context.Context`, closed channel).
- Ranged loops over containers are acceptable.

---

## 3. Minimise allocation and indirection
- After initialisation, avoid heap allocation (`make`, `new`, `append` that grow).
- Pre-size slices/maps; reuse buffers.
- Allow at most one level of pointer indirection.
- Forbid `unsafe` and reflection in critical paths; if used, isolate in a tiny reviewed package.

---

## 4. Write small, cohesive functions
- Keep functions to ~60 lines max (one screen).
- Split rather than nest logic.
- Average ~2 assertions per function (assertions must not mutate state).

---

## 5. Check all results and validate inputs
- Never discard `err`; always check return values.
- Validate inputs at function boundaries.
- Propagate errors with context (`fmt.Errorf("... %w", err)`).
- Justify any ignored results in comments.

---

## 6. Minimise scope and visibility
- Declare variables in the smallest possible scope.
- Avoid package-level state.
- Prefer unexported identifiers.
- Pass dependencies explicitly.

---

## 7. Enforce clarity, simplicity, and consistency
- Favour the simplest code that meets requirements.
- Use core language first, then stdlib, then internal libraries—before adding dependencies.
- Match existing code style; consistency beats novelty unless justified.

---

## 8. Follow canonical Go naming
- Use `MixedCaps`; no underscores (except `_test.go`, cgo, or generated code).
- Preserve case for initialisms (`ID`, `URL`, `API`).
- Avoid redundant names (`widget.New`, not `widget.NewWidget`).
- Short names for small scopes, longer for larger ones.

---

## 9. Package and file hygiene
- Package names: short, lowercase, descriptive; no `util` or `helper`.
- Group strongly related functionality; avoid monoliths and excess fragmentation.
- File names may include underscores; identifiers should not.
- In tests, name doubles clearly (`Stub`, `AlwaysCharges`, etc.).

---

## 10. Documentation and comments
- Write doc comments as full sentences, starting with the identifier name.
- Focus on rationale (*why*), not repetition of code (*what*).
- Provide runnable examples for public APIs.
- Wrap for readability (~80 cols), but no hard limit.

---

## 11. Imports and dependencies
- Organise imports in this order:
  1. Standard library
  1. External packages
  1. (Optional) protobuf imports (`foopb …`)
  1. Side-effect imports (`_ "…"`)
- Keep build tags and generated code minimal, each justified.

---

## 12. Testing and analysis
- Maintain comprehensive tests with clear diagnostics.
- Incorporate fuzzing tests for critical input validation and parsing logic using Go's built-in fuzzing `go test -fuzz=Fuzz` for functions named `Fuzz*`.
- Run `go vet ./...`, `golangci-lint run ./...`, and `go test -race ./...` with zero warnings.
- Add benchmarks to enforce allocation rules in hot paths.

---

## 13. Error handling
- Treat errors as values; design clear, user-friendly messages.
- Prefer direct returns over panics; use `panic` only for truly unrecoverable states.

---

## 14. Quick reference
- **Formatting**: always `gofmt`.
- **Identifiers**: MixedCaps, case-preserving initialisms.
- **Functions**: small, cohesive, bounded.
- **Errors**: always checked, propagated with context.
- **Safety**: no recursion, unsafe, or unchecked allocations.
- **Consistency**: align with existing code; deviations must be justified.

---

✅ This unified style blends **NASA's Power of Ten reliability rules** with **Google's Go style guide and best practices** into a concise, enforceable set for robust, maintainable Go code.
