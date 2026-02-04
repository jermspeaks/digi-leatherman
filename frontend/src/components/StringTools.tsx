import { useRef, useState } from 'react';
import { Copy } from 'lucide-react';
import { Button } from '@/components/ui/button';
import {
  urlEncode,
  urlDecode,
  base64Encode,
  base64Decode,
  trim,
  upperCase,
  lowerCase,
  capitalCase,
  snakeCase,
  kebabCase,
  camelCase,
  pascalCase,
  sentenceCase,
  parseUrlParams,
  createUrlWithParams,
} from '../api/stringTools';
import type { StringResult } from '../api/stringTools';

export type StringToolId =
  | 'url-encode'
  | 'url-decode'
  | 'url-parse-params'
  | 'url-param-creator'
  | 'base64-encode'
  | 'base64-decode'
  | 'trim'
  | 'upper-case'
  | 'lower-case'
  | 'capital-case'
  | 'snake-case'
  | 'kebab-case'
  | 'camel-case'
  | 'pascal-case'
  | 'sentence-case';

type ToolConfig = {
  id: StringToolId;
  label: string;
  description: string;
  example: { input: string; output: string };
  placeholder: string;
  buttonLabel: string;
  apiFn: (value: string) => Promise<StringResult>;
};

const TOOL_CONFIG: ToolConfig[] = [
  {
    id: 'url-encode',
    label: 'URL encode',
    description: 'Encodes text for safe use in URL query strings (spaces become +, special characters percent-encoded).',
    example: { input: 'hello world', output: 'hello+world' },
    placeholder: 'Text to encode…',
    buttonLabel: 'Encode',
    apiFn: urlEncode,
  },
  {
    id: 'url-decode',
    label: 'URL decode',
    description: 'Decodes URL-encoded text back to plain text.',
    example: { input: 'hello+world', output: 'hello world' },
    placeholder: 'Encoded string to decode…',
    buttonLabel: 'Decode',
    apiFn: urlDecode,
  },
  {
    id: 'url-parse-params',
    label: 'URL - Parser for params',
    description: 'Parses query parameters from a full URL or raw query string and returns them as JSON.',
    example: { input: 'https://example.com?foo=bar&baz=qux', output: '{"baz":["qux"],"foo":["bar"]}' },
    placeholder: 'URL or query string (e.g. ?foo=bar&baz=qux)…',
    buttonLabel: 'Parse',
    apiFn: parseUrlParams,
  },
  {
    id: 'url-param-creator',
    label: 'URL - param creator',
    description: 'Builds a URL from a base URL (first line) and key=value params (one per line).',
    example: { input: 'https://api.example.com/search\nq=hello\nlimit=10', output: 'https://api.example.com/search?limit=10&q=hello' },
    placeholder: 'First line: base URL. Following lines: key=value (one per line)…',
    buttonLabel: 'Create URL',
    apiFn: createUrlWithParams,
  },
  {
    id: 'base64-encode',
    label: 'Base64 encode',
    description: 'Encodes text to Base64.',
    example: { input: 'Hi', output: 'SGk=' },
    placeholder: 'Text to encode…',
    buttonLabel: 'Encode',
    apiFn: base64Encode,
  },
  {
    id: 'base64-decode',
    label: 'Base64 decode',
    description: 'Decodes Base64 back to plain text.',
    example: { input: 'SGk=', output: 'Hi' },
    placeholder: 'Base64 string to decode…',
    buttonLabel: 'Decode',
    apiFn: base64Decode,
  },
  {
    id: 'trim',
    label: 'Trim',
    description: 'Removes leading and trailing whitespace.',
    example: { input: '  hello world  ', output: 'hello world' },
    placeholder: 'Text to trim…',
    buttonLabel: 'Trim',
    apiFn: trim,
  },
  {
    id: 'upper-case',
    label: 'Upper Case',
    description: 'Converts all characters to uppercase.',
    example: { input: 'Hello World', output: 'HELLO WORLD' },
    placeholder: 'Text to convert…',
    buttonLabel: 'Convert',
    apiFn: upperCase,
  },
  {
    id: 'lower-case',
    label: 'Lower Case',
    description: 'Converts all characters to lowercase.',
    example: { input: 'Hello World', output: 'hello world' },
    placeholder: 'Text to convert…',
    buttonLabel: 'Convert',
    apiFn: lowerCase,
  },
  {
    id: 'capital-case',
    label: 'Capital Case',
    description: 'Converts to title case (first letter of each word uppercase).',
    example: { input: 'hello world', output: 'Hello World' },
    placeholder: 'Text to convert…',
    buttonLabel: 'Convert',
    apiFn: capitalCase,
  },
  {
    id: 'snake-case',
    label: 'Snake Case',
    description: 'Converts words to snake_case (lowercase with underscores).',
    example: { input: 'hello world', output: 'hello_world' },
    placeholder: 'Text to convert…',
    buttonLabel: 'Convert',
    apiFn: snakeCase,
  },
  {
    id: 'kebab-case',
    label: 'Kebab Case',
    description: 'Converts words to kebab-case (lowercase with hyphens).',
    example: { input: 'hello world', output: 'hello-world' },
    placeholder: 'Text to convert…',
    buttonLabel: 'Convert',
    apiFn: kebabCase,
  },
  {
    id: 'camel-case',
    label: 'Camel Case',
    description: 'Converts words to camelCase (first word lowercase, rest capitalized).',
    example: { input: 'hello world', output: 'helloWorld' },
    placeholder: 'Text to convert…',
    buttonLabel: 'Convert',
    apiFn: camelCase,
  },
  {
    id: 'pascal-case',
    label: 'Pascal Case',
    description: 'Converts words to PascalCase (each word capitalized).',
    example: { input: 'hello world', output: 'HelloWorld' },
    placeholder: 'Text to convert…',
    buttonLabel: 'Convert',
    apiFn: pascalCase,
  },
  {
    id: 'sentence-case',
    label: 'Sentence Case',
    description: 'Converts to sentence case (first character uppercase, rest lowercase).',
    example: { input: 'HELLO WORLD', output: 'Hello world' },
    placeholder: 'Text to convert…',
    buttonLabel: 'Convert',
    apiFn: sentenceCase,
  },
];

/** Map of tool id to description for command palette search. */
export const TOOL_DESCRIPTIONS: Record<string, string> = Object.fromEntries(
  TOOL_CONFIG.map((c) => [c.id, c.description])
);

const TOOL_MAP = Object.fromEntries(TOOL_CONFIG.map((c) => [c.id, c])) as Record<StringToolId, ToolConfig>;

type StringToolsProps = {
  /** When set, show only this tool (no tabs). Used for per-route views. */
  tool?: StringToolId;
};

export function StringTools({ tool: toolProp }: StringToolsProps) {
  const [tab, setTab] = useState<StringToolId>(toolProp ?? 'url-encode');
  const [input, setInput] = useState('');
  const [output, setOutput] = useState('');
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);
  const [copied, setCopied] = useState(false);
  const copyTimeoutRef = useRef<ReturnType<typeof setTimeout> | null>(null);

  const tool = toolProp ?? tab;
  const showTabs = toolProp == null;
  const config = TOOL_MAP[tool];

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

  const useExample = () => setInput(config.example.input);

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
      // Permission denied or unsupported; leave button state unchanged
    }
  };

  return (
    <section className="text-left mt-6">
      <h2 className="mb-3 text-2xl text-text">Strings</h2>
      {showTabs && (
        <nav className="flex gap-2 mb-4" aria-label="String tools">
          {TOOL_CONFIG.map((c) => (
            <button
              key={c.id}
              type="button"
              className={`py-2 px-4 ${tab === c.id ? 'border-accent bg-sidebar-active' : ''}`}
              onClick={() => setTab(c.id)}
            >
              {c.label}
            </button>
          ))}
        </nav>
      )}
      <p className="mb-2 text-text-secondary max-w-2xl">{config.description}</p>
      <p className="mb-4 text-sm text-text-secondary max-w-2xl">
        Example: <code className="px-1 rounded bg-bg-elevated">{config.example.input}</code> → <code className="px-1 rounded bg-bg-elevated">{config.example.output}</code>
        <button type="button" onClick={useExample} className="ml-2 text-accent hover:underline">
          Use example
        </button>
      </p>
      <div className="flex flex-col gap-3 max-w-2xl">
        <label htmlFor="string-input" className="font-medium">
          Input
        </label>
        <textarea
          id="string-input"
          className="w-full min-h-24 p-2.5 font-inherit text-base rounded-lg border border-border bg-bg text-text resize-y"
          value={input}
          onChange={(e) => setInput(e.target.value)}
          placeholder={config.placeholder}
          rows={4}
        />
        <button type="button" onClick={run} disabled={loading}>
          {loading ? '…' : config.buttonLabel}
        </button>
        {error && (
          <p className="text-red-500 m-0" role="alert">
            {error}
          </p>
        )}
        {output && (
          <div className="mt-2">
            <div className="flex items-center justify-between gap-2 mb-1">
              <label className="font-medium">Output</label>
              <Button
                type="button"
                variant="outline"
                size="sm"
                onClick={handleCopy}
                className="shrink-0"
                aria-label={copied ? 'Copied' : 'Copy to clipboard'}
              >
                <Copy className="size-3.5" aria-hidden />
                {copied ? 'Copied!' : 'Copy'}
              </Button>
            </div>
            <pre className="m-0 p-3 bg-bg-elevated rounded-lg overflow-x-auto whitespace-pre-wrap break-all">
              {output}
            </pre>
          </div>
        )}
      </div>
    </section>
  );
}
