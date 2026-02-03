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
  placeholder: string;
  buttonLabel: string;
  apiFn: (value: string) => Promise<StringResult>;
};

const TOOL_CONFIG: ToolConfig[] = [
  { id: 'url-encode', label: 'URL encode', placeholder: 'Text to encode…', buttonLabel: 'Encode', apiFn: urlEncode },
  { id: 'url-decode', label: 'URL decode', placeholder: 'Encoded string to decode…', buttonLabel: 'Decode', apiFn: urlDecode },
  { id: 'base64-encode', label: 'Base64 encode', placeholder: 'Text to encode…', buttonLabel: 'Encode', apiFn: base64Encode },
  { id: 'base64-decode', label: 'Base64 decode', placeholder: 'Base64 string to decode…', buttonLabel: 'Decode', apiFn: base64Decode },
  { id: 'trim', label: 'Trim', placeholder: 'Text to trim…', buttonLabel: 'Trim', apiFn: trim },
  { id: 'upper-case', label: 'Upper Case', placeholder: 'Text to convert…', buttonLabel: 'Convert', apiFn: upperCase },
  { id: 'lower-case', label: 'Lower Case', placeholder: 'Text to convert…', buttonLabel: 'Convert', apiFn: lowerCase },
  { id: 'capital-case', label: 'Capital Case', placeholder: 'Text to convert…', buttonLabel: 'Convert', apiFn: capitalCase },
  { id: 'snake-case', label: 'Snake Case', placeholder: 'Text to convert…', buttonLabel: 'Convert', apiFn: snakeCase },
  { id: 'kebab-case', label: 'Kebab Case', placeholder: 'Text to convert…', buttonLabel: 'Convert', apiFn: kebabCase },
  { id: 'camel-case', label: 'Camel Case', placeholder: 'Text to convert…', buttonLabel: 'Convert', apiFn: camelCase },
  { id: 'pascal-case', label: 'Pascal Case', placeholder: 'Text to convert…', buttonLabel: 'Convert', apiFn: pascalCase },
  { id: 'sentence-case', label: 'Sentence Case', placeholder: 'Text to convert…', buttonLabel: 'Convert', apiFn: sentenceCase },
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
