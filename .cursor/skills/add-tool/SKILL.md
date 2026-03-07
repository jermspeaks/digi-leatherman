---
name: add-tool
description: Guide for adding new tools to Digital Leatherman. Use when creating a new string tool, JSON tool, generator, or any new utility. Triggers on "add tool", "new tool", "create tool", "add feature".
---

# Add Tool Workflow

Follow this checklist to add a new tool to both backend and frontend.

## Checklist

Copy and track progress:

```
- [ ] 1. Backend handler
- [ ] 2. Register route in main.go
- [ ] 3. Write handler tests
- [ ] 4. Frontend API client
- [ ] 5. Sidebar config entry
- [ ] 6. Tool component config
- [ ] 7. Add route in App.tsx (if new page)
- [ ] 8. Test end-to-end
```

## 1. Backend Handler

Create or add to a handler file in `backend/handlers/`.

**For string tools**, add to `string.go` or create a new file:

```go
func MyTool(w http.ResponseWriter, r *http.Request) {
	req, ok := decodeBody(w, r)
	if !ok {
		return
	}
	result := processValue(req.Value) // your logic here
	writeJSON(w, StringResponse{Result: result})
}
```

**For tools with custom request/response**, define new types:

```go
type MyToolRequest struct {
	Value   string `json:"value"`
	Option  string `json:"option"`
}

type MyToolResponse struct {
	Result string `json:"result"`
	Count  int    `json:"count,omitempty"`
}
```

## 2. Register Route

In `backend/main.go`, add the route with CORS wrapper:

```go
mux.HandleFunc("/api/category/my-tool", cors(handlers.MyTool))
```

Route naming:
- String tools: `/api/string/tool-name`
- JSON tools: `/api/json/tool-name`
- Generators: `/api/lorem-ipsum/tool-name` or new category

## 3. Write Handler Tests

In `backend/handlers/my_tool_test.go` or existing test file, use table-driven tests:

```go
func TestMyTool(t *testing.T) {
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
			status, body := runHandler(t, MyTool, tc.method, tc.body)
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

Run tests: `cd backend && go test ./handlers/...`

## 4. Frontend API Client

Add to appropriate file in `frontend/src/api/` (e.g., `stringTools.ts`):

```typescript
export async function myTool(value: string): Promise<StringResult> {
  return postString('/api/category/my-tool', value);
}
```

For custom payloads:

```typescript
export async function myTool(value: string, option: string): Promise<MyResult> {
  const res = await fetch(`${API_BASE}/api/category/my-tool`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ value, option }),
  });
  if (!res.ok) throw new Error(await res.text() || `HTTP ${res.status}`);
  return res.json();
}
```

### API Helper Functions

The project uses reusable helpers in API files. Prefer these over raw fetch for consistency.

**postString** (for simple string tools):

```typescript
async function postString(path: string, value: string): Promise<StringResult> {
  const res = await fetch(`${API_BASE}${path}`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ value }),
  });
  if (!res.ok) throw new Error(await res.text() || `HTTP ${res.status}`);
  return res.json();
}
```

**postJSON** (typed, for complex payloads):

```typescript
async function postJSON<T>(path: string, body: object): Promise<T> {
  const res = await fetch(`${API_BASE}${path}`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(body),
  });
  if (!res.ok) throw new Error(await res.text() || `HTTP ${res.status}`);
  return res.json();
}
```

## 5. Sidebar Config

Add entry in `frontend/src/config/sidebarConfig.ts`:

```typescript
{ id: 'my-tool', label: 'My Tool', path: '/tools/category/my-tool', subGroup: 'GroupName' },
```

- `subGroup` creates collapsible sections within a category
- Single-item subGroups render as direct links

## 6. Tool Component Config

Add to the appropriate tool component's config array (e.g., `StringTools.tsx`):

```typescript
{
  id: 'my-tool',
  label: 'My Tool',
  description: 'What this tool does in one sentence.',
  example: { input: 'sample input', output: 'sample output' },
  placeholder: 'Enter text to process…',
  buttonLabel: 'Process',
  apiFn: myTool,
},
```

Don't forget to:
1. Import the API function
2. Add the ID to the `ToolId` type union

## 7. Add Route (If New Page)

For tools in existing categories, the component handles routing via props.

For new standalone pages, add in `frontend/src/App.tsx`:

```tsx
<Route path="/tools/category/my-tool" element={<MyToolPage />} />
```

## 8. Test End-to-End

1. Start backend: `cd backend && go run .`
2. Start frontend: `cd frontend && npm run dev`
3. Navigate to the new tool in the sidebar
4. Test with example input
5. Verify error handling (empty input, invalid data)
