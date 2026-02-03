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
import './StringTools.css';

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
    <section className="string-tools">
      <h2>Strings</h2>
      {showTabs && (
        <nav className="tool-tabs" aria-label="String tools">
          {TOOL_CONFIG.map((c) => (
            <button
              key={c.id}
              type="button"
              className={tab === c.id ? 'active' : ''}
              onClick={() => setTab(c.id)}
            >
              {c.label}
            </button>
          ))}
        </nav>
      )}
      <div className="tool-card">
        <label htmlFor="string-input">Input</label>
        <textarea
          id="string-input"
          value={input}
          onChange={(e) => setInput(e.target.value)}
          placeholder={config.placeholder}
          rows={4}
        />
        <button type="button" onClick={run} disabled={loading}>
          {loading ? '…' : config.buttonLabel}
        </button>
        {error && (
          <p className="tool-error" role="alert">
            {error}
          </p>
        )}
        {output && (
          <div className="tool-output">
            <label>Output</label>
            <pre>{output}</pre>
          </div>
        )}
      </div>
    </section>
  );
}
