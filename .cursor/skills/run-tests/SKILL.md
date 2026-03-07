---
name: run-tests
description: Run tests for Digital Leatherman. Use when testing code, running test suites, checking coverage, or verifying changes. Triggers on "run tests", "test", "testing", "go test", "npm test".
---

# Running Tests

## Backend (Go)

### Run All Tests

```bash
cd backend
go test ./...
```

### Run Handler Tests Only

```bash
cd backend
go test ./handlers/...
```

### Run with Coverage

```bash
cd backend
go test ./handlers/... -cover
```

### Run Specific Test

```bash
cd backend
go test ./handlers/... -run TestURLEncode
```

### Verbose Output

```bash
cd backend
go test ./... -v
```

### Run with Race Detection

```bash
cd backend
go test ./... -race
```

## Frontend

### Lint Check

```bash
cd frontend
npm run lint
```

### Type Check

```bash
cd frontend
npx tsc --noEmit
```

### Build (includes type check)

```bash
cd frontend
npm run build
```

## Test File Locations

- Handler tests: `backend/handlers/*_test.go`
- String handler tests: `backend/handlers/string_test.go`
- Lorem ipsum tests: `backend/handlers/lorem_test.go`
- JSON handler tests: `backend/handlers/json_test.go`

## Writing New Tests

Use the table-driven test pattern:

```go
func TestMyHandler(t *testing.T) {
	cases := []struct {
		name        string
		method      string
		body        string
		wantStatus  int
		wantResult  string
		checkResult bool
	}{
		{"method not allowed", "GET", "", http.StatusMethodNotAllowed, "", false},
		{"invalid JSON", "POST", "{", http.StatusBadRequest, "", false},
		{"valid input", "POST", `{"value":"test"}`, http.StatusOK, "expected", true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			status, body := runHandler(t, MyHandler, tc.method, tc.body)
			if status != tc.wantStatus {
				t.Errorf("status = %d, want %d", status, tc.wantStatus)
			}
			if tc.checkResult {
				got := parseResult(t, body)
				if got != tc.wantResult {
					t.Errorf("result = %q, want %q", got, tc.wantResult)
				}
			}
		})
	}
}
```

## Test Helpers

Available in `backend/handlers/string_test.go`:

- `runHandler(t, handler, method, body)` - Execute handler, returns (status, body)
- `parseResult(t, body)` - Extract `result` field from JSON response

## Quick Verification

After making changes:

```bash
# Backend
cd backend && go test ./... && cd ..

# Frontend
cd frontend && npm run lint && npm run build && cd ..
```
