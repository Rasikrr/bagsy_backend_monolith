---
description: Generate or review Go tests following project conventions
argument-hint: <target> — path to package/file (e.g. internal/usecases/schedule) or "describe" to explain conventions
allowed-tools: Bash, Read, Edit, Write, Glob, Grep, Agent
---

# Go Testing Conventions

When `$ARGUMENTS` is provided:
- If it equals **"describe"** — explain the testing conventions from this document relevant to the target layer. Do NOT write code.
- Otherwise treat `$ARGUMENTS` as a **package or file path**. Detect which layer it belongs to (domain / usecases / handlers / repositories / pkg) and generate comprehensive tests following the conventions below. Read all source files in the target first.

Follow these conventions when writing, reviewing, or generating tests for this project.

---

## 1. Unit Tests — Domain Layer

Domain tests are **pure unit tests**. No external dependencies, no mocks — just constructors, methods, and assertions.

### File Placement

Place `_test.go` files next to the source file in the **same package** (white-box testing). Use black-box (`_test` suffix package) only when testing exclusively through the public API.

```
internal/domain/identity/
├── employee.go
├── employee_test.go      # same package: `package identity`
└── errors.go
```

### Test Naming

```go
// Function under test + scenario + expected outcome
func TestNewEmployee_EmptyPhone_ReturnsError(t *testing.T) { ... }

// For method tests: Type_Method_Scenario
func TestEmployee_Deactivate_AlreadyDeleted_ReturnsError(t *testing.T) { ... }
```

### Table-Driven Tests

Use table-driven tests when testing the same function with multiple inputs/outputs.

```go
func TestNewPhone(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        wantErr bool
    }{
        {"valid phone", "+77001234567", false},
        {"empty string", "", true},
        {"too short", "+7700", true},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            _, err := shared.NewPhone(tt.input)
            if tt.wantErr {
                require.Error(t, err)
            } else {
                require.NoError(t, err)
            }
        })
    }
}
```

**Rules for table-driven tests:**
- Keep the struct flat. If you need >5-6 fields or conditional logic inside the loop, split into separate test functions.
- Use descriptive `name` fields — they appear in `go test -v` output.
- Every case MUST have a `name` and use `t.Run(tt.name, ...)`.

### Subtests with t.Run

Use `t.Run` to group related assertions under a parent test when testing different behaviors of the same method.

```go
func TestOrganization_Activation(t *testing.T) {
    org, _ := organization.NewStubOrganization()

    t.Run("deactivate active organization", func(t *testing.T) {
        err := org.Deactivate()
        assert.NoError(t, err)
        assert.False(t, org.Active)
    })

    t.Run("activate on deleted organization", func(t *testing.T) {
        _ = org.Delete()
        err := org.Activate()
        assert.ErrorIs(t, err, organization.ErrOrganizationDeleted)
    })
}
```

### Test Helpers

Mark all helper functions with `t.Helper()` so failure locations point to the calling test, not the helper.

```go
func newTestEmployee(t *testing.T) *identity.Employee {
    t.Helper()
    emp, err := identity.NewEmployee(identity.CreateEmployeeParams{
        Phone: mustPhone(t, "+77001234567"),
        OrgID: uuid.New(),
    })
    require.NoError(t, err)
    return emp
}

func mustPhone(t *testing.T, raw string) shared.Phone {
    t.Helper()
    p, err := shared.NewPhone(raw)
    require.NoError(t, err)
    return p
}
```

**When helpers do NOT accept `*testing.T`** (e.g., used in benchmarks too), use `require`-style panics or return errors. Avoid package-level `must*` helpers that panic without `t.Helper()` context — they produce confusing stack traces.

### t.Cleanup

Use `t.Cleanup` instead of `defer` when the cleanup must run after subtests complete.

```go
func TestWithTempFile(t *testing.T) {
    f, err := os.CreateTemp("", "test-*")
    require.NoError(t, err)
    t.Cleanup(func() { os.Remove(f.Name()) })

    // ... test logic using f
}
```

### t.Parallel

Use `t.Parallel()` for tests that are **truly independent** — no shared mutable state.

```go
func TestSlugGeneration(t *testing.T) {
    tests := []struct {
        name  string
        input string
        want  string
    }{
        {"latin", "My Salon", "my-salon"},
        {"cyrillic", "Мой Салон", "moj-salon"},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            t.Parallel()
            got := shared.GenerateSlug(tt.input)
            assert.Equal(t, tt.want, got)
        })
    }
}
```

**Do NOT use `t.Parallel()` when:**
- Tests share mutable state (e.g., modifying the same entity across subtests).
- Tests depend on execution order.
- Tests modify global variables, files, or databases without isolation.

**Note on loop variable capture:** In Go 1.22+ the loop variable is per-iteration, so `tt := tt` is no longer required. For Go <1.22, always add `tt := tt` before `t.Run` in parallel table tests.

---

## 2. Assertions — testify/assert vs testify/require

This project uses `github.com/stretchr/testify`.

### require — Stop on Failure

Use `require` when a failure makes subsequent assertions meaningless (preconditions, nil checks, length checks).

```go
emp, err := identity.NewEmployee(params)
require.NoError(t, err)           // if this fails, emp is nil — no point continuing
require.NotNil(t, emp)
assert.Equal(t, "John", emp.Name) // safe to check fields now
```

### assert — Continue on Failure

Use `assert` when you want to see ALL failures in one test run (checking multiple fields of a result).

```go
assert.Equal(t, expected.Name, got.Name)
assert.Equal(t, expected.Phone, got.Phone)
assert.Equal(t, expected.Role, got.Role)
```

### Decision Rule

| Situation | Use |
|-----------|-----|
| Error check before using result | `require.NoError` / `require.NotNil` |
| Length check before indexing | `require.Len` / `require.NotEmpty` |
| Checking multiple independent fields | `assert.*` |
| Single critical assertion in a test | `require.*` |

### Error Assertions

```go
// Check specific domain error
assert.ErrorIs(t, err, booking.ErrSlotNotAvailable)

// Check error is non-nil without matching type
require.Error(t, err)

// NEVER compare error strings directly:
// BAD: assert.Equal(t, "slot not available", err.Error())
```

### go-cmp (Optional, for Complex Structs)

For comparing large structs with custom comparison logic (ignoring fields, custom transformers), use `github.com/google/go-cmp/cmp`:

```go
import "github.com/google/go-cmp/cmp"

if diff := cmp.Diff(want, got, cmpopts.IgnoreFields(Foo{}, "ID", "CreatedAt")); diff != "" {
    t.Errorf("mismatch (-want +got):\n%s", diff)
}
```

Use go-cmp when testify's `assert.Equal` produces unreadable diffs for deeply nested structs.

---

## 3. Mocking

### Interface-Based Mocking (Preferred)

Define narrow interfaces at the consumer (use case), not the provider (repository).

```go
// In use case package — only the methods this use case needs
type employeeReader interface {
    GetByID(ctx context.Context, id uuid.UUID) (*identity.Employee, error)
}
```

### Hand-Written Mocks (Default Choice)

For interfaces with 1-3 methods, write mocks manually. They are explicit, easy to read, and require no code generation.

```go
type mockEmployeeRepo struct {
    getByIDFn func(ctx context.Context, id uuid.UUID) (*identity.Employee, error)
}

func (m *mockEmployeeRepo) GetByID(ctx context.Context, id uuid.UUID) (*identity.Employee, error) {
    return m.getByIDFn(ctx, id)
}
```

Usage in tests:

```go
repo := &mockEmployeeRepo{
    getByIDFn: func(_ context.Context, id uuid.UUID) (*identity.Employee, error) {
        if id == knownID {
            return testEmployee, nil
        }
        return nil, identity.ErrEmployeeNotFound
    },
}
uc := usecase.New(repo)
```

### Generated Mocks (When Justified)

Use `mockgen` (uber-go/mock) or `moq` only when:
- Interface has >5 methods.
- You need to verify call counts, argument matchers, or call ordering.
- Multiple use cases share the same large interface.

Place generated mocks in `internal/usecases/{context}/mocks/` or alongside the test file.

### When NOT to Mock

- **Domain layer**: Never mock. Domain entities are pure Go — test them directly.
- **Value Objects**: Never mock. Construct real instances.
- **Simple data transformations**: Test with real data, not mocks.
- **Repository layer**: Use integration tests with testcontainers instead of mocking the database driver.

---

## 4. Use Case Layer Tests

Use case tests verify orchestration: load entity, call domain method, save. Mock the repository interfaces.

```go
func TestFireEmployee(t *testing.T) {
    emp := newTestEmployee(t)

    repo := &mockEmployeeRepo{
        getByIDFn: func(_ context.Context, _ uuid.UUID) (*identity.Employee, error) {
            return emp, nil
        },
        saveFn: func(_ context.Context, e *identity.Employee) error {
            assert.True(t, e.IsDeleted())
            return nil
        },
    }

    uc := identity.NewUseCase(repo)
    err := uc.FireEmployee(context.Background(), emp.ID)
    require.NoError(t, err)
}
```

**Rules:**
- One test file per use case: `fire_employee_test.go`.
- Test the happy path + each domain error path + infra error propagation.
- Verify that the use case does NOT contain business logic — if a domain rule is complex, it should be tested in the domain layer.

---

## 5. HTTP Handler Layer Tests

Use `net/http/httptest` to test handlers without starting a real server.

```go
func TestGetEmployee_NotFound(t *testing.T) {
    uc := &mockUseCase{
        getEmployeeFn: func(_ context.Context, _ uuid.UUID) (*identity.Employee, error) {
            return nil, identity.ErrEmployeeNotFound
        },
    }

    h := handler.NewEmployeeHandler(uc)

    req := httptest.NewRequest(http.MethodGet, "/api/v1/employees/"+uuid.New().String(), nil)
    rec := httptest.NewRecorder()

    h.GetEmployee(rec, req)

    assert.Equal(t, http.StatusNotFound, rec.Code)

    var resp map[string]string
    require.NoError(t, json.NewDecoder(rec.Body).Decode(&resp))
    assert.Equal(t, "employee_not_found", resp["code"])
}
```

**Rules:**
- Test HTTP status codes and response body structure, not internal implementation.
- Test error mapping: each domain error should map to the correct HTTP status + slug code.
- Test request validation: malformed JSON, missing fields, invalid UUIDs.
- Use `chi.RouteContext` or middleware chains when handlers depend on URL params from the router.

---

## 6. Integration Tests — Repository Layer

Use `testcontainers-go` for testing against real PostgreSQL/Redis.

### Build Tags

Separate integration tests with build tags:

```go
//go:build integration

package repository_test
```

Run with: `go test -tags=integration ./internal/repositories/...`

### TestMain for Shared Container

```go
//go:build integration

package repository_test

var testDB *pgxpool.Pool

func TestMain(m *testing.M) {
    ctx := context.Background()

    pgContainer, err := postgres.Run(ctx,
        "postgres:16-alpine",
        postgres.WithDatabase("test_db"),
        postgres.WithUsername("test"),
        postgres.WithPassword("test"),
        testcontainers.WithWaitStrategy(
            wait.ForLog("database system is ready to accept connections").
                WithOccurrence(2).
                WithStartupTimeout(5*time.Second),
        ),
    )
    if err != nil {
        log.Fatal(err)
    }
    defer pgContainer.Terminate(ctx)

    connStr, _ := pgContainer.ConnectionString(ctx, "sslmode=disable")
    testDB, err = pgxpool.New(ctx, connStr)
    if err != nil {
        log.Fatal(err)
    }

    // Run migrations
    runMigrations(connStr)

    os.Exit(m.Run())
}
```

### Per-Test Isolation

Use transactions that rollback after each test:

```go
func withTx(t *testing.T, db *pgxpool.Pool) pgx.Tx {
    t.Helper()
    ctx := context.Background()
    tx, err := db.Begin(ctx)
    require.NoError(t, err)
    t.Cleanup(func() { tx.Rollback(ctx) })
    return tx
}
```

---

## 7. Test Fixtures and testdata/

### testdata/ Directory

Go ignores `testdata/` directories when building. Use them for:
- Golden files (expected outputs).
- JSON fixtures for complex test inputs.
- SQL seed data for integration tests.

```
internal/domain/booking/
├── slots.go
├── slots_test.go
└── testdata/
    └── golden/
        └── TestGenerateSlots_FullDay.golden
```

### Golden File Pattern

```go
func TestRenderResponse(t *testing.T) {
    got := renderResponse(input)

    golden := filepath.Join("testdata", "golden", t.Name()+".golden")

    if *update {
        os.WriteFile(golden, got, 0644)
        return
    }

    want, err := os.ReadFile(golden)
    require.NoError(t, err)
    assert.Equal(t, string(want), string(got))
}
```

Use golden files for Swagger specs, large JSON responses, or rendered templates. Update with `-update` flag.

---

## 8. Benchmarks

```go
func BenchmarkGenerateSlots(b *testing.B) {
    b.ReportAllocs()
    locSlots := []*schedule.LocationScheduleSlot{locWorkSlot(testDate, 0, 0, 24, 0)}
    dur := mustDuration(30)
    step := mustDuration(15)

    for b.Loop() {
        GenerateSlots(
            location.ScheduleTypeFixed,
            locSlots, nil, nil,
            dur, step,
            testDate, testDate,
            makeTime(testDate, 0, 0),
        )
    }
}
```

**Rules:**
- Use `b.ReportAllocs()` to track allocations.
- Use `b.ResetTimer()` after expensive setup that should not count.
- Run with `-benchmem` and `-count=10` for statistical significance.
- Place benchmarks in the same `_test.go` file as related unit tests.
- Use `b.Loop()` (Go 1.24+) instead of `for i := 0; i < b.N; i++` when available.

---

## 9. Fuzz Testing

Use for parsing, validation, and serialization functions — especially Value Objects.

```go
func FuzzNewPhone(f *testing.F) {
    // Seed corpus with representative inputs
    f.Add("+77001234567")
    f.Add("")
    f.Add("+1")
    f.Add("not-a-phone")

    f.Fuzz(func(t *testing.T, input string) {
        phone, err := shared.NewPhone(input)
        if err != nil {
            return // invalid input is fine — we just want no panics
        }
        // Round-trip: valid phones must survive String() -> NewPhone()
        phone2, err := shared.NewPhone(phone.String())
        require.NoError(t, err)
        assert.Equal(t, phone, phone2)
    })
}
```

**When to fuzz:**
- Value Object constructors (`NewPhone`, `NewMoney`, `NewSlug`, `NewDuration`).
- JSON/request parsing functions.
- Any function that must never panic on arbitrary input.

**When NOT to fuzz:**
- Business logic with complex preconditions.
- Functions requiring valid aggregate state.

Corpus files go in `testdata/fuzz/<FuzzTestName>/`.

---

## 10. Test Isolation and Parallel Safety

### Global State

- NEVER use package-level mutable variables in tests. If you need shared test data, use constants or create fresh instances per test.
- Exception: immutable test fixtures like `var testDate = time.Date(...)` are safe.

### Cleanup Order

`t.Cleanup` runs in LIFO order (last registered, first called). Register cleanups immediately after resource creation:

```go
db := openTestDB(t)
t.Cleanup(func() { db.Close() })

tx := beginTx(t, db)
t.Cleanup(func() { tx.Rollback(ctx) })
// tx cleanup runs first, then db cleanup
```

### Race Detection

Always run tests with `-race` in CI:

```bash
go test -race ./...
```

---

## 11. Coverage

### When Coverage Matters

- **Domain layer**: Aim for high coverage (80%+). This is pure business logic.
- **Use case layer**: Cover happy path + each error branch.
- **Handler layer**: Cover each error mapping + request validation.
- **Repository layer**: Covered by integration tests, not unit test coverage.

### When Coverage Does NOT Matter

- Do not chase 100%. Covering trivial getters or generated code is waste.
- Never write tests solely to increase coverage numbers.
- Untested code is better than tests that verify nothing useful.

### Commands

```bash
go test -cover ./internal/domain/...           # Quick coverage check
go test -coverprofile=coverage.out ./...        # Full profile
go tool cover -html=coverage.out               # Visual report
```

---

## 12. Anti-Patterns to Avoid

### Over-Mocking
If your test has more mock setup than actual assertions, the design needs work. Mock at the boundary (repository interface), not deep inside.

### Testing Implementation Details
```go
// BAD: tests internal method call sequence
mockRepo.AssertCalled(t, "Save", mock.Anything)
mockRepo.AssertNumberOfCalls(t, "Save", 1)

// GOOD: tests observable outcome
emp, err := uc.FireEmployee(ctx, id)
require.NoError(t, err)
assert.True(t, emp.IsDeleted())
```

### Complicated Table Tests
If a table test has >6 fields, conditional setup logic, or `if tt.wantSpecialCase` branches inside the loop, break it into separate named test functions.

### Evergreen Tests
Tests that can never fail are useless. After writing a test, verify it fails when the condition is violated (red-green discipline).

### Asserting on Irrelevant Details
```go
// BAD: entire struct comparison when only status matters
assert.Equal(t, expectedEmployee, gotEmployee)

// GOOD: assert only what matters
assert.True(t, gotEmployee.IsDeleted())
assert.NotNil(t, gotEmployee.DeletedAt)
```

### Violating Encapsulation for Tests
Never export private methods just to test them. Test through the public API. If a private function is complex enough to need its own tests, consider extracting it into a separate (internal) package.

### Ignoring Test Errors
```go
// BAD: silent failure
result, _ := doSomething()

// GOOD: explicit check
result, err := doSomething()
require.NoError(t, err)
```

---

## 13. Project Test Structure Summary

```
internal/
├── domain/{context}/
│   ├── entity.go
│   ├── entity_test.go          # Pure unit tests, same package
│   └── testdata/               # Golden files, fixtures
├── usecases/{context}/
│   ├── create_foo.go
│   ├── create_foo_test.go      # Unit tests with mocked repos
│   └── mocks/                  # Generated mocks (if needed)
├── ports/http/handlers/{context}/
│   ├── handler.go
│   └── handler_test.go         # httptest-based tests
├── repositories/{context}/
│   ├── repository.go
│   └── repository_integration_test.go  # //go:build integration
└── testutil/                   # Shared test helpers across packages
    ├── db.go                   # testcontainers setup
    └── fixtures.go             # Common test entity builders
```

---

## Quick Reference

| Layer | Test Type | Mocking | Framework |
|-------|-----------|---------|-----------|
| Domain | Unit | None | testify |
| Use Case | Unit | Hand-written mocks for repos | testify |
| Policy | Unit | Hand-written mocks | testify |
| Handler | Unit/Integration | Mock use case interface | testify + httptest |
| Repository | Integration | None (real DB) | testify + testcontainers |
| Pkg | Unit | None or hand-written | testify |
| Workers | Unit + Integration | Mock repos for unit, real DB for integration | testify |
