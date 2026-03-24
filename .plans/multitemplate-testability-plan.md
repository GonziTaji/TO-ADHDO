Plan: make `getDomainData` easy to test

Goal
- Refactor `getDomainData` in `server/multitemplate/multitemplate.go` so its behavior can be covered with small, deterministic unit tests.
- Preserve current behavior where possible, while fixing obvious correctness issues discovered during the refactor.

Current problems
- The function reads the real filesystem directly via `os.Stat`, `os.ReadDir`, and `filepath.WalkDir`.
- It depends on package-global state (`cfg`) instead of explicit inputs.
- It writes warnings through the global logger, which makes assertions awkward.
- It mixes multiple responsibilities: page discovery, layout resolution, component discovery, and resource discovery.
- Error handling is inconsistent; for example, `os.ReadDir` is called and its error is not handled before iteration.
- There are likely correctness bugs in the helper predicates and walk logic, so tests should lock behavior down before and after cleanup.

Refactor strategy

1) Freeze current expected behavior with characterization tests
- Add tests around the public behavior of `getDomainData` using `t.TempDir()` and small fixture trees.
- Cover at least:
  - domain with default layout and one valid page
  - page with its own layout overriding the default
  - page directory missing `page.html` gets skipped
  - domain missing default layout logs or reports a warning
  - resources and component templates are discovered
- These tests can use the real filesystem initially; the goal is to create a safety net before deeper refactoring.

2) Introduce explicit dependencies
- Replace direct package calls with injected dependencies collected in a small struct, for example:
  - file existence check
  - read directory
  - walk directory
  - warning sink/logger
- Keep the first version minimal; do not abstract more than `getDomainData` actually needs.
- Recommended shape:
  - `type fileSystem interface { Stat(...); ReadDir(...); WalkDir(...) }`
  - `type domainLoader struct { fs fileSystem; cfg Config; warn func(string, ...any) }`
- Keep a package-level default implementation for production code so callers do not become noisy.

3) Move config out of globals for this path
- Extract the relevant fields from the global `cfg` into a dedicated immutable config type passed into the loader.
- Keep the existing global `cfg` as the default source for production wiring, but stop reading it directly inside the core logic.
- This allows table-driven tests to use tiny custom configs without mutating package globals.

4) Split the function by responsibility
- Break `getDomainData` into small helpers with narrow outputs:
  - `getDefaultLayoutPath(...)`
  - `collectPages(...)`
  - `collectComponentsAndResources(...)`
  - small predicates like `isTemplateFile(...)` and `isResourceFile(...)`
- Each helper should return data and warnings/errors rather than logging internally.
- Keep one thin orchestration function that assembles `DomainTemplatesData`.

5) Replace logging side effects with returned diagnostics
- Prefer returning warnings from the core logic, or passing in a warning callback.
- Production code can still log warnings, but tests should be able to assert on collected warnings without hijacking the global logger.
- If changing the function signature is too disruptive, add an internal helper that returns warnings and keep `getDomainData` as a wrapper.

6) Fix correctness issues during the refactor
- Review and correct the extension helpers:
  - `isTemplate` currently checks resource extensions
  - `isResource` currently checks template extensions
- Review the page-folder existence check; it currently stats `default_layout_path` twice instead of checking the pages directory.
- Review `filepath.WalkDir` path handling:
  - avoid joining `path` with `d.Name()` when `path` already includes the entry name
  - ensure page/layout files are skipped correctly
  - ensure the page directory exclusion uses the actual pages directory path
- Add targeted tests for each bug fixed so the behavior stays stable.

7) Keep one integration-style test with real files
- Even after dependency injection is added, keep at least one end-to-end test using `t.TempDir()` to validate real path behavior.
- This protects against path-joining mistakes that mocks may hide.

8) Add focused unit tests with a fake filesystem
- Once the filesystem is injected, add table-driven unit tests for edge cases that are hard to create on disk:
  - `Stat` errors other than not-exist
  - `ReadDir` failures
  - `WalkDir` callback error propagation
  - empty domain
  - mixed-case file extensions
- Use the fake only where it improves coverage or simplicity; do not replace every integration-style test.

9) Wire production code through the new seam
- Update `processDomain` to call the new loader/wrapper without changing its external behavior.
- Keep logging in `processDomain` or in a thin production adapter, not in the pure discovery helpers.

10) Definition of done
- `getDomainData` core logic no longer depends directly on `os`, `filepath.WalkDir`, or the global logger.
- Tests cover the happy path, skip behavior, missing-layout behavior, discovery behavior, and major error paths.
- Global config is not mutated in tests.
- The current domain loading flow still works through `processDomain` and `loadTemplatesAuto`.

Suggested implementation order
- Add characterization tests with `t.TempDir()`.
- Extract page collection and component/resource collection into helpers.
- Introduce injected config and warning sink.
- Introduce filesystem abstraction only where needed.
- Fix discovered bugs with tests added alongside each fix.
- Run the relevant Go tests for `server/multitemplate`.

Notes
- A full virtual filesystem may be unnecessary at first; simple function injection or a narrow interface is enough.
- If the package has few callers, it is acceptable to change `getDomainData` into an internal helper and expose a thin compatibility wrapper.
