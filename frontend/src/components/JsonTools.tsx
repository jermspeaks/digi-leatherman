import { useRef, useState } from 'react';
import { Copy } from 'lucide-react';
import { Button } from '@/components/ui/button';
import {
  formatJson,
  minifyJson,
  validateJson,
  pathQueryJson,
  diffJson,
} from '../api/jsonTools';
import type { JsonResult, ValidateResult } from '../api/jsonTools';

export type JsonToolId = 'format' | 'minify' | 'validate' | 'path' | 'diff';

type ToolConfig = {
  id: JsonToolId;
  label: string;
  description: string;
  example: { input: string; output?: string; path?: string };
  exampleB?: string;
  placeholder: string;
  placeholderPath?: string;
  placeholderB?: string;
  buttonLabel: string;
};

const TOOL_CONFIG: ToolConfig[] = [
  {
    id: 'format',
    label: 'Format',
    description: 'Pretty-print JSON with 2-space indentation.',
    example: { input: '{"a":1,"b":2}', output: '{\n  "a": 1,\n  "b": 2\n}' },
    placeholder: 'Paste JSON…',
    buttonLabel: 'Format',
  },
  {
    id: 'minify',
    label: 'Minify',
    description: 'Remove unnecessary whitespace from JSON.',
    example: { input: '{\n  "a": 1\n}', output: '{"a":1}' },
    placeholder: 'Paste JSON…',
    buttonLabel: 'Minify',
  },
  {
    id: 'validate',
    label: 'Validate',
    description: 'Check whether the input is valid JSON.',
    example: { input: '{"valid": true}', output: 'Valid' },
    placeholder: 'Paste JSON…',
    buttonLabel: 'Validate',
  },
  {
    id: 'path',
    label: 'Path query',
    description: 'Extract a value at a dot-separated path (e.g. "a.b" or "items.0.name").',
    example: { input: '{"a":{"b":42}}', path: 'a.b', output: '42' },
    placeholder: 'Paste JSON…',
    placeholderPath: 'e.g. a.b or items.0',
    buttonLabel: 'Query',
  },
  {
    id: 'diff',
    label: 'Diff',
    description: 'Compare two JSON values and show structural differences.',
    example: { input: '{"a":1}', output: '(no differences)' },
    exampleB: '{"a":1}',
    placeholder: 'First JSON…',
    placeholderB: 'Second JSON…',
    buttonLabel: 'Diff',
  },
];

/** Map of tool id to description for command palette search. */
export const JSON_TOOL_DESCRIPTIONS: Record<string, string> = Object.fromEntries(
  TOOL_CONFIG.map((c) => [c.id, c.description])
);

const TOOL_MAP = Object.fromEntries(TOOL_CONFIG.map((c) => [c.id, c])) as Record<JsonToolId, ToolConfig>;

type JsonToolsProps = {
  tool?: JsonToolId;
};

export function JsonTools({ tool: toolProp }: JsonToolsProps) {
  const [tab, setTab] = useState<JsonToolId>(toolProp ?? 'format');
  const [input, setInput] = useState('');
  const [pathInput, setPathInput] = useState('');
  const [valueB, setValueB] = useState('');
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
    setLoading(true);
    try {
      if (tool === 'format') {
        const res: JsonResult = await formatJson(input);
        setOutput(res.result);
      } else if (tool === 'minify') {
        const res: JsonResult = await minifyJson(input);
        setOutput(res.result);
      } else if (tool === 'validate') {
        const res: ValidateResult = await validateJson(input);
        setOutput(res.valid ? 'Valid' : `Invalid: ${res.error ?? 'syntax error'}`);
      } else if (tool === 'path') {
        if (!pathInput.trim()) {
          setError('Enter a path');
          return;
        }
        const res: JsonResult = await pathQueryJson(input, pathInput);
        setOutput(res.result);
      } else if (tool === 'diff') {
        const res: JsonResult = await diffJson(input, valueB);
        setOutput(res.result);
      }
    } catch (e) {
      setError(e instanceof Error ? e.message : 'Request failed');
    } finally {
      setLoading(false);
    }
  };

  const useExample = () => {
    setInput(config.example.input);
    if (config.placeholderPath && config.example.path) setPathInput(config.example.path);
    if (config.exampleB) setValueB(config.exampleB);
  };

  const canRun =
    tool === 'path' ? input.trim() && pathInput.trim() :
    tool === 'diff' ? input.trim() && valueB.trim() :
    input.trim().length > 0;

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

  const textareaClass = 'w-full min-h-24 p-2.5 font-mono text-sm rounded-lg border border-border bg-bg text-text resize-y';

  return (
    <section className="text-left mt-6">
      <h2 className="mb-3 text-2xl text-text">JSON</h2>
      {showTabs && (
        <nav className="flex gap-2 mb-4 flex-wrap" aria-label="JSON tools">
          {TOOL_CONFIG.map((c) => (
            <button
              key={c.id}
              type="button"
              className={`py-2 px-4 rounded ${tab === c.id ? 'border border-accent bg-sidebar-active' : ''}`}
              onClick={() => setTab(c.id)}
            >
              {c.label}
            </button>
          ))}
        </nav>
      )}
      <p className="mb-2 text-text-secondary max-w-2xl">{config.description}</p>
      <p className="mb-4 text-sm text-text-secondary max-w-2xl">
        Example: <code className="px-1 rounded bg-bg-elevated font-mono text-xs">{config.example.input}</code>
        {config.example.path != null && (
          <> path <code className="px-1 rounded bg-bg-elevated font-mono text-xs">{config.example.path}</code></>
        )}
        {config.example.output != null && (
          <> → <code className="px-1 rounded bg-bg-elevated font-mono text-xs">{config.example.output}</code></>
        )}
        <button type="button" onClick={useExample} className="ml-2 text-accent hover:underline">
          Use example
        </button>
      </p>
      <div className="flex flex-col gap-3 max-w-2xl">
        <label htmlFor="json-input" className="font-medium">
          {tool === 'diff' ? 'JSON A' : 'Input'}
        </label>
        <textarea
          id="json-input"
          className={textareaClass}
          value={input}
          onChange={(e) => setInput(e.target.value)}
          placeholder={config.placeholder}
          rows={4}
        />
        {tool === 'path' && (
          <>
            <label htmlFor="json-path" className="font-medium">Path</label>
            <input
              id="json-path"
              type="text"
              className="w-full p-2.5 font-mono text-sm rounded-lg border border-border bg-bg text-text"
              value={pathInput}
              onChange={(e) => setPathInput(e.target.value)}
              placeholder={config.placeholderPath}
            />
          </>
        )}
        {tool === 'diff' && (
          <>
            <label htmlFor="json-input-b" className="font-medium">JSON B</label>
            <textarea
              id="json-input-b"
              className={textareaClass}
              value={valueB}
              onChange={(e) => setValueB(e.target.value)}
              placeholder={config.placeholderB}
              rows={4}
            />
          </>
        )}
        <button type="button" onClick={run} disabled={loading || !canRun}>
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
            <pre className="m-0 p-3 bg-bg-elevated rounded-lg overflow-x-auto whitespace-pre-wrap break-all font-mono text-sm">
              {output}
            </pre>
          </div>
        )}
      </div>
    </section>
  );
}
