---
name: react-tool-component
description: React tool component patterns for Digital Leatherman frontend. Use when creating tool UI components, adding new tools to existing components, or working with tool state. Triggers on "tool component", "react component", "frontend tool", "tool ui".
---

# React Tool Component Patterns

Tool components follow a config-driven pattern with consistent state management.

## ToolConfig Structure

Each tool is defined by a config object:

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

Export descriptions for command palette:

```typescript
export const TOOL_DESCRIPTIONS: Record<string, string> = Object.fromEntries(
  TOOL_CONFIG.map((c) => [c.id, c.description])
);
```
