import { useRef, useState } from 'react';
import { Copy } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { generate } from '../api/loremIpsum';
import type {
  LoremToolId,
  LoremType,
  LoremOptions,
} from '../api/loremIpsum';

/** Sidebar/route tool id (camel-case) -> API tool name (camelcase). */
const TOOL_ID_TO_API: Record<string, string> = {
  'camel-case': 'camelcase',
};

function toApiTool(toolId: string): string {
  return TOOL_ID_TO_API[toolId] ?? toolId;
}

const TYPE_OPTIONS: { value: LoremType; label: string }[] = [
  { value: 'words', label: 'Words' },
  { value: 'sentences', label: 'Sentences' },
  { value: 'paragraphs', label: 'Paragraphs' },
];

const VOCABULARY_OPTIONS: { value: 'default' | 'bacon' | 'hipster'; label: string }[] = [
  { value: 'default', label: 'Default' },
  { value: 'bacon', label: 'Bacon' },
  { value: 'hipster', label: 'Hipster' },
];

const MAX: Record<string, number> = {
  generator: 50,
  characters: 10000,
  bytes: 10000,
  title: 20,
  slug: 20,
  'camel-case': 20,
  list: 50,
  headings: 20,
  html: 20,
  markdown: 20,
  json: 10,
};

/** Descriptions for command palette search. */
export const LOREM_TOOL_DESCRIPTIONS: Record<string, string> = {
  generator: 'Generate placeholder text by words, sentences, or paragraphs.',
  characters: 'Generate a string of exactly N characters (optionally whole words only).',
  bytes: 'Generate a string of exactly N bytes (optionally whole words only).',
  title: 'Generate one short title phrase of N words.',
  slug: 'Generate a URL-friendly slug of N hyphenated words.',
  'camel-case': 'Generate a camelCase variable name of N words.',
  list: 'Generate a bullet or numbered list of N lorem items.',
  headings: 'Generate N heading lines (plain or markdown style).',
  html: 'Generate lorem wrapped in HTML (paragraphs or list).',
  markdown: 'Generate lorem in Markdown (paragraphs, list, or headings).',
  json: 'Generate a JSON object with lorem values for API mocking.',
};

type LoremToolsProps = {
  tool: string;
};

export function LoremTools({ tool }: LoremToolsProps) {
  const toolId = tool;
  const maxCount = MAX[toolId] ?? 100;
  const [count, setCount] = useState(toolId === 'generator' ? 3 : 5);
  const [type, setType] = useState<LoremType>('paragraphs');
  const [startWithClassic, setStartWithClassic] = useState(false);
  const [vocabulary, setVocabulary] = useState<'default' | 'bacon' | 'hipster'>('default');
  const [wholeWordsOnly, setWholeWordsOnly] = useState(false);
  const [listStyle, setListStyle] = useState<'bullet' | 'numbered'>('bullet');
  const [headingStyle, setHeadingStyle] = useState<'plain' | 'markdown'>('plain');
  const [format, setFormat] = useState<'paragraphs' | 'list' | 'headings'>('paragraphs');
  const [jsonKeys, setJsonKeys] = useState('title, body, summary');
  const [result, setResult] = useState('');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [copied, setCopied] = useState(false);
  const copyTimeoutRef = useRef<ReturnType<typeof setTimeout> | null>(null);

  const countClamped = Math.min(Math.max(1, count), maxCount);

  const buildOptions = (): LoremOptions => {
    const opts: LoremOptions = {};
    if (toolId === 'generator') {
      opts.type = type;
      opts.startWithClassic = startWithClassic;
      opts.vocabulary = vocabulary;
    }
    if (toolId === 'characters' || toolId === 'bytes') {
      opts.wholeWordsOnly = wholeWordsOnly;
    }
    if (toolId === 'list') {
      opts.listStyle = listStyle;
    }
    if (toolId === 'headings') {
      opts.headingStyle = headingStyle;
    }
    if (toolId === 'html' || toolId === 'markdown') {
      opts.format = format;
    }
    if (toolId === 'json') {
      const keys = jsonKeys
        .split(/[,;\s]+/)
        .map((k) => k.trim())
        .filter(Boolean);
      if (keys.length) opts.keys = keys;
    }
    return opts;
  };

  const run = async () => {
    setError(null);
    setResult('');
    if (copyTimeoutRef.current) {
      clearTimeout(copyTimeoutRef.current);
      copyTimeoutRef.current = null;
    }
    setCopied(false);
    setLoading(true);
    try {
      const { result: text } = await generate({
        tool: toApiTool(toolId) as LoremToolId,
        count: countClamped,
        options: buildOptions(),
      });
      setResult(text);
    } catch (e) {
      setError(e instanceof Error ? e.message : 'Request failed');
    } finally {
      setLoading(false);
    }
  };

  const handleCopy = async () => {
    if (!result) return;
    if (copyTimeoutRef.current) clearTimeout(copyTimeoutRef.current);
    try {
      await navigator.clipboard.writeText(result);
      setCopied(true);
      copyTimeoutRef.current = setTimeout(() => {
        copyTimeoutRef.current = null;
        setCopied(false);
      }, 1500);
    } catch {
      // Permission denied or unsupported
    }
  };

  const generatorMax =
    type === 'words' ? 1000 : type === 'sentences' ? 100 : 50;

  return (
    <section className="text-left mt-6">
      <h2 className="mb-3 text-2xl text-text">
        {toolId === 'generator'
          ? 'Generator'
          : toolId === 'camel-case'
            ? 'Camel Case'
            : toolId.charAt(0).toUpperCase() + toolId.slice(1)}
      </h2>
      <p className="mb-4 text-text-secondary max-w-2xl">
        {LOREM_TOOL_DESCRIPTIONS[toolId] ?? 'Generate lorem placeholder content.'}
      </p>
      <div className="flex flex-col gap-3 max-w-2xl">
        {/* Generator-specific */}
        {toolId === 'generator' && (
          <>
            <div>
              <label htmlFor="lorem-type" className="font-medium block mb-2">
                Type
              </label>
              <div
                id="lorem-type"
                className="flex gap-4"
                role="radiogroup"
                aria-label="Output type"
              >
                {TYPE_OPTIONS.map((opt) => (
                  <label
                    key={opt.value}
                    className="flex items-center gap-2 cursor-pointer"
                  >
                    <input
                      type="radio"
                      name="lorem-type"
                      value={opt.value}
                      checked={type === opt.value}
                      onChange={() => {
                        setType(opt.value);
                        const maxForType =
                          opt.value === 'words' ? 1000 : opt.value === 'sentences' ? 100 : 50;
                        setCount((c) => Math.min(c, maxForType));
                      }}
                      className="rounded border-border"
                    />
                    <span>{opt.label}</span>
                  </label>
                ))}
              </div>
            </div>
            <div>
              <label className="font-medium block mb-2">Options</label>
              <label className="flex items-center gap-2 cursor-pointer mb-2">
                <input
                  type="checkbox"
                  checked={startWithClassic}
                  onChange={(e) => setStartWithClassic(e.target.checked)}
                  className="rounded border-border"
                />
                <span>Start with classic Latin</span>
              </label>
              <div className="flex items-center gap-4">
                <span className="text-text-secondary">Vocabulary:</span>
                {VOCABULARY_OPTIONS.map((opt) => (
                  <label
                    key={opt.value}
                    className="flex items-center gap-2 cursor-pointer"
                  >
                    <input
                      type="radio"
                      name="vocabulary"
                      value={opt.value}
                      checked={vocabulary === opt.value}
                      onChange={() => setVocabulary(opt.value)}
                      className="rounded border-border"
                    />
                    <span>{opt.label}</span>
                  </label>
                ))}
              </div>
            </div>
          </>
        )}

        {/* Characters / Bytes */}
        {(toolId === 'characters' || toolId === 'bytes') && (
          <label className="flex items-center gap-2 cursor-pointer">
            <input
              type="checkbox"
              checked={wholeWordsOnly}
              onChange={(e) => setWholeWordsOnly(e.target.checked)}
              className="rounded border-border"
            />
            <span>Whole words only</span>
          </label>
        )}

        {/* List */}
        {toolId === 'list' && (
          <div>
            <label className="font-medium block mb-2">List style</label>
            <div className="flex gap-4">
              <label className="flex items-center gap-2 cursor-pointer">
                <input
                  type="radio"
                  name="list-style"
                  checked={listStyle === 'bullet'}
                  onChange={() => setListStyle('bullet')}
                  className="rounded border-border"
                />
                <span>Bullet</span>
              </label>
              <label className="flex items-center gap-2 cursor-pointer">
                <input
                  type="radio"
                  name="list-style"
                  checked={listStyle === 'numbered'}
                  onChange={() => setListStyle('numbered')}
                  className="rounded border-border"
                />
                <span>Numbered</span>
              </label>
            </div>
          </div>
        )}

        {/* Headings */}
        {toolId === 'headings' && (
          <div>
            <label className="font-medium block mb-2">Heading style</label>
            <div className="flex gap-4">
              <label className="flex items-center gap-2 cursor-pointer">
                <input
                  type="radio"
                  name="heading-style"
                  checked={headingStyle === 'plain'}
                  onChange={() => setHeadingStyle('plain')}
                  className="rounded border-border"
                />
                <span>Plain</span>
              </label>
              <label className="flex items-center gap-2 cursor-pointer">
                <input
                  type="radio"
                  name="heading-style"
                  checked={headingStyle === 'markdown'}
                  onChange={() => setHeadingStyle('markdown')}
                  className="rounded border-border"
                />
                <span>Markdown</span>
              </label>
            </div>
          </div>
        )}

        {/* HTML / Markdown */}
        {(toolId === 'html' || toolId === 'markdown') && (
          <div>
            <label className="font-medium block mb-2">Format</label>
            <div className="flex gap-4 flex-wrap">
              <label className="flex items-center gap-2 cursor-pointer">
                <input
                  type="radio"
                  name="format"
                  checked={format === 'paragraphs'}
                  onChange={() => setFormat('paragraphs')}
                  className="rounded border-border"
                />
                <span>Paragraphs</span>
              </label>
              <label className="flex items-center gap-2 cursor-pointer">
                <input
                  type="radio"
                  name="format"
                  checked={format === 'list'}
                  onChange={() => setFormat('list')}
                  className="rounded border-border"
                />
                <span>List</span>
              </label>
              {toolId === 'markdown' && (
                <label className="flex items-center gap-2 cursor-pointer">
                  <input
                    type="radio"
                    name="format"
                    checked={format === 'headings'}
                    onChange={() => setFormat('headings')}
                    className="rounded border-border"
                  />
                  <span>Headings</span>
                </label>
              )}
            </div>
          </div>
        )}

        {/* JSON */}
        {toolId === 'json' && (
          <div>
            <label htmlFor="json-keys" className="font-medium block mb-2">
              Keys (comma-separated)
            </label>
            <input
              id="json-keys"
              type="text"
              value={jsonKeys}
              onChange={(e) => setJsonKeys(e.target.value)}
              placeholder="title, body, summary"
              className="w-full max-w-md p-2.5 font-inherit text-base rounded-lg border border-border bg-bg text-text"
            />
          </div>
        )}

        {/* Count (all tools) */}
        <div>
          <label htmlFor="lorem-count" className="font-medium block mb-2">
            {toolId === 'generator'
              ? `Count (1–${generatorMax})`
              : `Count (1–${maxCount})`}
          </label>
          <input
            id="lorem-count"
            type="number"
            min={1}
            max={toolId === 'generator' ? generatorMax : maxCount}
            value={count}
            onChange={(e) =>
              setCount(
                Math.min(
                  toolId === 'generator' ? generatorMax : maxCount,
                  Math.max(1, e.target.valueAsNumber || 1)
                )
              )
            }
            className="w-full max-w-32 p-2.5 font-inherit text-base rounded-lg border border-border bg-bg text-text"
          />
        </div>

        <button
          type="button"
          onClick={run}
          disabled={loading}
          className="self-start"
        >
          {loading ? '…' : 'Generate'}
        </button>
        {error && (
          <p className="text-red-500 m-0" role="alert">
            {error}
          </p>
        )}
        {result && (
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
              {result}
            </pre>
          </div>
        )}
      </div>
    </section>
  );
}
