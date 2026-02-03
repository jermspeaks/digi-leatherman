import { useState } from 'react';
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
} from '../api/stringTools';
import type { StringResult } from '../api/stringTools';

export type StringToolId =
  | 'url-encode'
  | 'url-decode'
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

  const tool = toolProp ?? tab;
  const showTabs = toolProp == null;
  const config = TOOL_MAP[tool];

  const run = async () => {
    setError(null);
    setOutput('');
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
            <label className="block mb-1">Output</label>
            <pre className="m-0 p-3 bg-bg-elevated rounded-lg overflow-x-auto whitespace-pre-wrap break-all">
              {output}
            </pre>
          </div>
        )}
      </div>
    </section>
  );
}
