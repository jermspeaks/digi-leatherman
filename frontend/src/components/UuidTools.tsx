import { useRef, useState } from 'react';
import { Copy, RefreshCw } from 'lucide-react';
import { Button } from '@/components/ui/button';
import {
  generateUUIDv4,
  generateUUIDv7,
  bulkGenerateUUID,
  validateUUID,
  parseUUID,
  formatUUID,
} from '../api/uuidTools';

export type UuidToolId =
  | 'generate-v4'
  | 'generate-v7'
  | 'bulk-generate'
  | 'validate'
  | 'parse'
  | 'format';

type ToolConfig = {
  id: UuidToolId;
  label: string;
  description: string;
};

const TOOL_CONFIG: ToolConfig[] = [
  {
    id: 'generate-v4',
    label: 'Generate v4',
    description: 'Generates a random UUID v4 (RFC 4122). This is the most commonly used UUID version.',
  },
  {
    id: 'generate-v7',
    label: 'Generate v7',
    description: 'Generates a time-ordered UUID v7 (RFC 9562). Includes millisecond timestamp for natural sorting.',
  },
  {
    id: 'bulk-generate',
    label: 'Bulk Generate',
    description: 'Generate multiple UUIDs at once. Choose version and count (max 1000).',
  },
  {
    id: 'validate',
    label: 'Validate',
    description: 'Check if a string is a valid UUID and get its version and variant information.',
  },
  {
    id: 'parse',
    label: 'Parse',
    description: 'Extract detailed metadata from a UUID including version, variant, and timestamp (for v1/v7).',
  },
  {
    id: 'format',
    label: 'Format',
    description: 'Convert a UUID to different formats: standard, no-hyphens, uppercase, braces, or URN.',
  },
];

export const UUID_TOOL_DESCRIPTIONS: Record<string, string> = Object.fromEntries(
  TOOL_CONFIG.map((c) => [c.id, c.description])
);

const TOOL_MAP = Object.fromEntries(TOOL_CONFIG.map((c) => [c.id, c])) as Record<UuidToolId, ToolConfig>;

const FORMAT_OPTIONS = [
  { value: 'standard', label: 'Standard (with hyphens)' },
  { value: 'no-hyphens', label: 'No hyphens' },
  { value: 'uppercase', label: 'Uppercase' },
  { value: 'braces', label: 'With braces {uuid}' },
  { value: 'urn', label: 'URN (urn:uuid:...)' },
];

type UuidToolsProps = {
  tool?: UuidToolId;
};

export function UuidTools({ tool: toolProp }: UuidToolsProps) {
  const [input, setInput] = useState('');
  const [output, setOutput] = useState('');
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);
  const [copied, setCopied] = useState(false);
  const copyTimeoutRef = useRef<ReturnType<typeof setTimeout> | null>(null);

  const [bulkCount, setBulkCount] = useState(10);
  const [bulkVersion, setBulkVersion] = useState<4 | 7>(4);
  const [bulkResults, setBulkResults] = useState<string[]>([]);

  const [validateResult, setValidateResult] = useState<{ valid: boolean; version?: number; variant?: string; error?: string } | null>(null);
  const [parseResult, setParseResult] = useState<{
    valid: boolean;
    version?: number;
    variant?: string;
    timestamp?: string | null;
    clockSeq?: number | null;
    node?: string | null;
    error?: string;
  } | null>(null);

  const [formatType, setFormatType] = useState('standard');

  const tool = toolProp ?? 'generate-v4';
  const config = TOOL_MAP[tool];

  const clearState = () => {
    setOutput('');
    setError(null);
    setBulkResults([]);
    setValidateResult(null);
    setParseResult(null);
    if (copyTimeoutRef.current) {
      clearTimeout(copyTimeoutRef.current);
      copyTimeoutRef.current = null;
    }
    setCopied(false);
  };

  const handleGenerateV4 = async () => {
    clearState();
    setLoading(true);
    try {
      const { result } = await generateUUIDv4();
      setOutput(result);
    } catch (e) {
      setError(e instanceof Error ? e.message : 'Request failed');
    } finally {
      setLoading(false);
    }
  };

  const handleGenerateV7 = async () => {
    clearState();
    setLoading(true);
    try {
      const { result } = await generateUUIDv7();
      setOutput(result);
    } catch (e) {
      setError(e instanceof Error ? e.message : 'Request failed');
    } finally {
      setLoading(false);
    }
  };

  const handleBulkGenerate = async () => {
    clearState();
    setLoading(true);
    try {
      const { results } = await bulkGenerateUUID(bulkVersion, bulkCount);
      setBulkResults(results);
    } catch (e) {
      setError(e instanceof Error ? e.message : 'Request failed');
    } finally {
      setLoading(false);
    }
  };

  const handleValidate = async () => {
    clearState();
    if (!input.trim()) return;
    setLoading(true);
    try {
      const result = await validateUUID(input);
      setValidateResult(result);
    } catch (e) {
      setError(e instanceof Error ? e.message : 'Request failed');
    } finally {
      setLoading(false);
    }
  };

  const handleParse = async () => {
    clearState();
    if (!input.trim()) return;
    setLoading(true);
    try {
      const result = await parseUUID(input);
      setParseResult(result);
    } catch (e) {
      setError(e instanceof Error ? e.message : 'Request failed');
    } finally {
      setLoading(false);
    }
  };

  const handleFormat = async () => {
    clearState();
    if (!input.trim()) return;
    setLoading(true);
    try {
      const { result } = await formatUUID(input, formatType);
      setOutput(result);
    } catch (e) {
      setError(e instanceof Error ? e.message : 'Request failed');
    } finally {
      setLoading(false);
    }
  };

  const handleCopy = async (text: string) => {
    if (copyTimeoutRef.current) clearTimeout(copyTimeoutRef.current);
    try {
      await navigator.clipboard.writeText(text);
      setCopied(true);
      copyTimeoutRef.current = setTimeout(() => {
        copyTimeoutRef.current = null;
        setCopied(false);
      }, 1500);
    } catch {
      // Permission denied or unsupported
    }
  };

  const renderGenerateTool = (version: 4 | 7) => {
    const handler = version === 4 ? handleGenerateV4 : handleGenerateV7;
    const exampleOutput = version === 4
      ? '550e8400-e29b-41d4-a716-446655440000'
      : '01936b3c-7d4a-7123-8456-426614174000';

    return (
      <div className="flex flex-col gap-4 max-w-2xl">
        <p className="text-sm text-text-secondary">
          Example output: <code className="px-1 rounded bg-bg-elevated">{exampleOutput}</code>
        </p>
        <Button onClick={handler} disabled={loading} className="w-fit">
          <RefreshCw className={`size-4 mr-2 ${loading ? 'animate-spin' : ''}`} />
          {loading ? 'Generating...' : `Generate UUID v${version}`}
        </Button>
        {output && (
          <div className="mt-2">
            <div className="flex items-center justify-between gap-2 mb-1">
              <label className="font-medium">Result</label>
              <Button
                type="button"
                variant="outline"
                size="sm"
                onClick={() => handleCopy(output)}
                className="shrink-0"
              >
                <Copy className="size-3.5" aria-hidden />
                {copied ? 'Copied!' : 'Copy'}
              </Button>
            </div>
            <pre className="m-0 p-3 bg-bg-elevated rounded-lg overflow-x-auto font-mono text-lg">
              {output}
            </pre>
          </div>
        )}
      </div>
    );
  };

  const renderBulkGenerate = () => (
    <div className="flex flex-col gap-4 max-w-2xl">
      <div className="flex gap-4 items-end">
        <div className="flex flex-col gap-1">
          <label htmlFor="bulk-version" className="text-sm font-medium">Version</label>
          <select
            id="bulk-version"
            value={bulkVersion}
            onChange={(e) => setBulkVersion(Number(e.target.value) as 4 | 7)}
            className="p-2 rounded-lg border border-border bg-bg text-text"
          >
            <option value={4}>v4 (random)</option>
            <option value={7}>v7 (time-ordered)</option>
          </select>
        </div>
        <div className="flex flex-col gap-1">
          <label htmlFor="bulk-count" className="text-sm font-medium">Count</label>
          <input
            id="bulk-count"
            type="number"
            min={1}
            max={1000}
            value={bulkCount}
            onChange={(e) => setBulkCount(Math.min(1000, Math.max(1, Number(e.target.value))))}
            className="p-2 w-24 rounded-lg border border-border bg-bg text-text"
          />
        </div>
        <Button onClick={handleBulkGenerate} disabled={loading}>
          <RefreshCw className={`size-4 mr-2 ${loading ? 'animate-spin' : ''}`} />
          {loading ? 'Generating...' : 'Generate'}
        </Button>
      </div>
      {bulkResults.length > 0 && (
        <div className="mt-2">
          <div className="flex items-center justify-between gap-2 mb-1">
            <label className="font-medium">Results ({bulkResults.length})</label>
            <Button
              type="button"
              variant="outline"
              size="sm"
              onClick={() => handleCopy(bulkResults.join('\n'))}
              className="shrink-0"
            >
              <Copy className="size-3.5" aria-hidden />
              {copied ? 'Copied!' : 'Copy all'}
            </Button>
          </div>
          <pre className="m-0 p-3 bg-bg-elevated rounded-lg overflow-auto max-h-80 font-mono text-sm">
            {bulkResults.join('\n')}
          </pre>
        </div>
      )}
    </div>
  );

  const renderValidate = () => (
    <div className="flex flex-col gap-4 max-w-2xl">
      <p className="text-sm text-text-secondary">
        Example: <code className="px-1 rounded bg-bg-elevated">550e8400-e29b-41d4-a716-446655440000</code>
      </p>
      <div className="flex flex-col gap-1">
        <label htmlFor="uuid-input" className="font-medium">UUID to validate</label>
        <input
          id="uuid-input"
          type="text"
          value={input}
          onChange={(e) => setInput(e.target.value)}
          placeholder="Enter a UUID..."
          className="p-2.5 rounded-lg border border-border bg-bg text-text font-mono"
        />
      </div>
      <Button onClick={handleValidate} disabled={loading} className="w-fit">
        {loading ? 'Validating...' : 'Validate'}
      </Button>
      {validateResult && (
        <div className="mt-2 p-4 bg-bg-elevated rounded-lg">
          {validateResult.valid ? (
            <div className="space-y-2">
              <p className="text-green-600 dark:text-green-400 font-medium">Valid UUID</p>
              <p><span className="text-text-secondary">Version:</span> {validateResult.version}</p>
              <p><span className="text-text-secondary">Variant:</span> {validateResult.variant}</p>
            </div>
          ) : (
            <p className="text-red-500">Invalid UUID{validateResult.error && `: ${validateResult.error}`}</p>
          )}
        </div>
      )}
    </div>
  );

  const renderParse = () => (
    <div className="flex flex-col gap-4 max-w-2xl">
      <p className="text-sm text-text-secondary">
        Example: <code className="px-1 rounded bg-bg-elevated">550e8400-e29b-41d4-a716-446655440000</code>
      </p>
      <div className="flex flex-col gap-1">
        <label htmlFor="uuid-parse-input" className="font-medium">UUID to parse</label>
        <input
          id="uuid-parse-input"
          type="text"
          value={input}
          onChange={(e) => setInput(e.target.value)}
          placeholder="Enter a UUID..."
          className="p-2.5 rounded-lg border border-border bg-bg text-text font-mono"
        />
      </div>
      <Button onClick={handleParse} disabled={loading} className="w-fit">
        {loading ? 'Parsing...' : 'Parse'}
      </Button>
      {parseResult && (
        <div className="mt-2 p-4 bg-bg-elevated rounded-lg">
          {parseResult.valid ? (
            <div className="space-y-2">
              <p className="text-green-600 dark:text-green-400 font-medium">Parsed UUID</p>
              <p><span className="text-text-secondary">Version:</span> {parseResult.version}</p>
              <p><span className="text-text-secondary">Variant:</span> {parseResult.variant}</p>
              {parseResult.timestamp && (
                <p><span className="text-text-secondary">Timestamp:</span> {parseResult.timestamp}</p>
              )}
              {parseResult.clockSeq != null && (
                <p><span className="text-text-secondary">Clock Sequence:</span> {parseResult.clockSeq}</p>
              )}
              {parseResult.node && (
                <p><span className="text-text-secondary">Node:</span> <code className="font-mono">{parseResult.node}</code></p>
              )}
            </div>
          ) : (
            <p className="text-red-500">Invalid UUID{parseResult.error && `: ${parseResult.error}`}</p>
          )}
        </div>
      )}
    </div>
  );

  const renderFormat = () => (
    <div className="flex flex-col gap-4 max-w-2xl">
      <p className="text-sm text-text-secondary">
        Example: <code className="px-1 rounded bg-bg-elevated">550e8400-e29b-41d4-a716-446655440000</code> → <code className="px-1 rounded bg-bg-elevated">550e8400e29b41d4a716446655440000</code> (no-hyphens)
      </p>
      <div className="flex flex-col gap-1">
        <label htmlFor="uuid-format-input" className="font-medium">UUID to format</label>
        <input
          id="uuid-format-input"
          type="text"
          value={input}
          onChange={(e) => setInput(e.target.value)}
          placeholder="Enter a UUID..."
          className="p-2.5 rounded-lg border border-border bg-bg text-text font-mono"
        />
      </div>
      <div className="flex gap-4 items-end">
        <div className="flex flex-col gap-1">
          <label htmlFor="format-type" className="text-sm font-medium">Format</label>
          <select
            id="format-type"
            value={formatType}
            onChange={(e) => setFormatType(e.target.value)}
            className="p-2 rounded-lg border border-border bg-bg text-text"
          >
            {FORMAT_OPTIONS.map((opt) => (
              <option key={opt.value} value={opt.value}>{opt.label}</option>
            ))}
          </select>
        </div>
        <Button onClick={handleFormat} disabled={loading}>
          {loading ? 'Formatting...' : 'Format'}
        </Button>
      </div>
      {output && (
        <div className="mt-2">
          <div className="flex items-center justify-between gap-2 mb-1">
            <label className="font-medium">Result</label>
            <Button
              type="button"
              variant="outline"
              size="sm"
              onClick={() => handleCopy(output)}
              className="shrink-0"
            >
              <Copy className="size-3.5" aria-hidden />
              {copied ? 'Copied!' : 'Copy'}
            </Button>
          </div>
          <pre className="m-0 p-3 bg-bg-elevated rounded-lg overflow-x-auto font-mono text-lg">
            {output}
          </pre>
        </div>
      )}
    </div>
  );

  const renderToolContent = () => {
    switch (tool) {
      case 'generate-v4':
        return renderGenerateTool(4);
      case 'generate-v7':
        return renderGenerateTool(7);
      case 'bulk-generate':
        return renderBulkGenerate();
      case 'validate':
        return renderValidate();
      case 'parse':
        return renderParse();
      case 'format':
        return renderFormat();
      default:
        return null;
    }
  };

  return (
    <section className="text-left mt-6">
      <h2 className="mb-3 text-2xl text-text">UUID</h2>
      <p className="mb-4 text-text-secondary max-w-2xl">{config.description}</p>
      {error && (
        <p className="text-red-500 mb-4" role="alert">
          {error}
        </p>
      )}
      {renderToolContent()}
    </section>
  );
}
