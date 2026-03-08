---
name: react-tool-component
description: React tool component patterns for Digital Leatherman frontend. Use when creating tool UI components, adding new tools to existing components, or working with tool state. Triggers on "tool component", "react component", "frontend tool", "tool ui".
---

# React Tool Component Patterns

Tool components follow a config-driven pattern with consistent state management.

## Component Files Reference

All tool components live in `frontend/src/components/*Tools.tsx`. Use an existing component when adding a tool in an existing category; create a new component only when adding a new category (new sidebar category and route group in `App.tsx`).

| Component | Structure |
|-----------|-----------|
| **StringTools** | Config array + `apiFn` per tool; generic run handler |
| **JsonTools** | Config array, no `apiFn`; switch in handler; optional `placeholderPath`, `placeholderB`, `exampleB` for path/diff |
| **UuidTools** | Config array + separate render logic per tool (dropdowns, number inputs) |
| **LoremTools** | Dynamic config + options builder; radios, checkboxes, type/vocabulary/list/format options |
| **DateTools** | `ToolMeta` config, switch-based handler, typed results; `OutputRow` and grid for countdown |
| **ColorTools** | `ToolMeta` config; color picker, contrast, palette, etc. |

## ToolConfig Structure

Each tool is defined by a config object. The shape varies by component; see **ToolConfig / ToolMeta variations** below.

```typescript
type ToolConfig = {
  id: StringToolId;
  label: string;
  description: string;
  example: { input: string; output: string };
  placeholder: string;
  buttonLabel: string;
  apiFn: (value: string) => Promise<StringResult>;
};
```

Example config entry:

```typescript
{
  id: 'url-encode',
  label: 'URL encode',
  description: 'Encodes text for safe use in URL query strings.',
  example: { input: 'hello world', output: 'hello+world' },
  placeholder: 'Text to encode…',
  buttonLabel: 'Encode',
  apiFn: urlEncode,
},
```

### ToolConfig / ToolMeta variations

- **StringTools:** `apiFn` in config; single input/output. Use the type above.
- **JsonTools:** No `apiFn`; handler uses a switch and calls API functions by tool id. Optional `placeholderPath`, `placeholderB`, `exampleB` for multi-input tools (path query, diff). See `frontend/src/components/JsonTools.tsx` for the exact type.
- **UuidTools / LoremTools:** Extra options in state (version, format, count, type, etc.); config has label/description; handler builds the request from state.
- **DateTools / ColorTools:** Use `ToolMeta` (label, description, etc.) without `apiFn`; switch in handler; result types vary (e.g. `DateParseResult`, `TimezoneResult`). See `frontend/src/components/DateTools.tsx` or `ColorTools.tsx` for the exact type.

## Tool ID Type

Define a union type for all tool IDs:

```typescript
export type StringToolId =
  | 'url-encode'
  | 'url-decode'
  | 'my-new-tool';  // Add new tools here
```

## State Management Pattern

Standard state for tool components:

```typescript
const [input, setInput] = useState('');
const [output, setOutput] = useState('');
const [error, setError] = useState<string | null>(null);
const [loading, setLoading] = useState(false);
const [copied, setCopied] = useState(false);
const copyTimeoutRef = useRef<ReturnType<typeof setTimeout> | null>(null);
```

## Run Handler Pattern

Execute the tool action:

```typescript
const run = async () => {
  setError(null);
  setOutput('');
  if (copyTimeoutRef.current) {
    clearTimeout(copyTimeoutRef.current);
    copyTimeoutRef.current = null;
  }
  setCopied(false);
  if (!input.trim()) return;
  
  setLoading(true);
  try {
    const { result } = await config.apiFn(input);
    setOutput(result);
  } catch (e) {
    setError(e instanceof Error ? e.message : 'Request failed');
  } finally {
    setLoading(false);
  }
};
```

## Copy to Clipboard Pattern

With timeout feedback:

```typescript
const handleCopy = async () => {
  if (!output) return;
  if (copyTimeoutRef.current) clearTimeout(copyTimeoutRef.current);
  try {
    await navigator.clipboard.writeText(output);
    setCopied(true);
    copyTimeoutRef.current = setTimeout(() => {
      copyTimeoutRef.current = null;
      setCopied(false);
    }, 1500);
  } catch {
    // Permission denied or unsupported
  }
};
```

## Component Structure

```tsx
export function StringTools({ tool: toolProp }: StringToolsProps) {
  // State
  const [input, setInput] = useState('');
  const [output, setOutput] = useState('');
  // ... other state

  const config = TOOL_MAP[tool];

  return (
    <section className="text-left mt-6">
      <h2 className="mb-3 text-2xl text-text">{category}</h2>
      
      {/* Description */}
      <p className="mb-2 text-text-secondary max-w-2xl">{config.description}</p>
      
      {/* Example with "Use example" button */}
      <p className="mb-4 text-sm text-text-secondary max-w-2xl">
        Example: <code>{config.example.input}</code> → <code>{config.example.output}</code>
        <button onClick={useExample}>Use example</button>
      </p>
      
      {/* Input */}
      <textarea
        value={input}
        onChange={(e) => setInput(e.target.value)}
        placeholder={config.placeholder}
      />
      
      {/* Action button */}
      <button onClick={run} disabled={loading}>
        {loading ? '…' : config.buttonLabel}
      </button>
      
      {/* Error display */}
      {error && <p className="text-red-500">{error}</p>}
      
      {/* Output with copy button */}
      {output && (
        <div>
          <Button onClick={handleCopy}>
            <Copy /> {copied ? 'Copied!' : 'Copy'}
          </Button>
          <pre>{output}</pre>
        </div>
      )}
    </section>
  );
}
```

## Sidebar and Routing

**Sidebar:** `frontend/src/config/sidebarConfig.ts`. Structure: `SidebarCategory[]` with `id`, `label`, `items`; each item is `SidebarItem` with `id`, `label`, `path`, optional `subGroup`. Items sharing the same `subGroup` render as one collapsible group; a single-item subGroup renders as a direct link. Path convention: `/tools/{category}/{tool-id}`. Use `getBreadcrumbLabels(path)` for breadcrumbs.

**Routing:** `frontend/src/App.tsx`. Nested routes under `<Layout />`; each tool route renders the category component with `tool="tool-id"` (e.g. `<StringTools tool="url-encode" />`). Visiting a category path (e.g. `/tools/json`) redirects to the first tool (e.g. `/tools/json/format`). Default redirect goes to `/tools/string/url-encode`.

## API Client Patterns

API modules live in `frontend/src/api/*.ts`.

- **Base URL:** `API_BASE` from `import.meta.env.VITE_API_URL` (fallback `http://localhost:8100`).
- **Helpers:** Use generic `postJSON<T>(path, body)` for typed responses; some files use a `postString`-style helper for simple `{ value }` payloads.
- **Result types:** Vary by category—e.g. `StringResult`, `JsonResult`, `ValidateResult`, UUID/Date/Color-specific types. See the relevant `frontend/src/api/*.ts` file for exact shapes.

## Special UI patterns

- **Copy to clipboard:** Use a ref and 1500 ms timeout to reset “Copied!” (see Copy to Clipboard Pattern above).
- **Multi-input:** JsonTools path tool uses an extra input for the path; diff tool uses a second textarea. Use state such as `pathInput`, `valueB`.
- **Dropdowns/select:** UuidTools (version, format), DateTools (timezone lists), LoremTools (type, vocabulary, list style, format).
- **Number inputs:** UuidTools bulk count, LoremTools; clamp with `Math.min(max, Math.max(min, value))`.
- **Radio groups:** LoremTools (type, vocabulary, list style, format).
- **Checkboxes:** LoremTools (e.g. `startWithClassic`, `wholeWordsOnly`).
- **Structured output:** DateTools uses an `OutputRow` component with per-field copy; countdown uses a grid (days/hours/minutes/seconds).
- **Loading:** Set `loading` state, disable the action button, use `animate-spin` on an icon for a spinner.
- **Errors:** Set `error` state in the catch block; render in red (e.g. `text-red-500`).

## Adding a New Tool

1. Import the API function from `@/api/`
2. Add ID to the `ToolId` type union
3. Add config entry to `TOOL_CONFIG` array
4. The component automatically handles the new tool

## CSS Classes Reference

Common Tailwind classes used:

- `text-text` / `text-text-secondary` - Theme text colors
- `bg-bg` / `bg-bg-elevated` - Theme backgrounds
- `border-border` - Theme border
- `text-accent` - Accent color (links, active states)
- `max-w-2xl` - Constrain content width

## Exports for Search

Export descriptions for the command palette. Naming varies by component:

- **StringTools, DateTools, ColorTools:** `TOOL_DESCRIPTIONS`
- **JsonTools:** `JSON_TOOL_DESCRIPTIONS`
- **UuidTools:** `UUID_TOOL_DESCRIPTIONS`
- **LoremTools:** `LOREM_TOOL_DESCRIPTIONS` (static object, not derived from config)

For new components, use `TOOL_DESCRIPTIONS` unless the category already uses a prefixed name. When derived from config:

```typescript
export const TOOL_DESCRIPTIONS: Record<string, string> = Object.fromEntries(
  TOOL_CONFIG.map((c) => [c.id, c.description])
);
```
