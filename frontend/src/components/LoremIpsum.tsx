import { useRef, useState } from 'react';
import { Copy } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { generate, type LoremType } from '../api/loremIpsum';

const MAX: Record<LoremType, number> = {
  words: 1000,
  sentences: 100,
  paragraphs: 50,
};

const TYPE_OPTIONS: { value: LoremType; label: string }[] = [
  { value: 'words', label: 'Words' },
  { value: 'sentences', label: 'Sentences' },
  { value: 'paragraphs', label: 'Paragraphs' },
];

/** Description for command palette search. */
export const LOREM_IPSUM_DESCRIPTION: Record<string, string> = {
  'lorem-ipsum':
    'Generate placeholder text by words, sentences, or paragraphs.',
};

export function LoremIpsum() {
  const [type, setType] = useState<LoremType>('paragraphs');
  const [count, setCount] = useState(3);
  const [result, setResult] = useState('');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [copied, setCopied] = useState(false);
  const copyTimeoutRef = useRef<ReturnType<typeof setTimeout> | null>(null);

  const maxCount = MAX[type];
  const countClamped = Math.min(Math.max(1, count), maxCount);

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
      const { result: text } = await generate(type, countClamped);
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

  return (
    <section className="text-left mt-6">
      <h2 className="mb-3 text-2xl text-text">Lorem Ipsum</h2>
      <p className="mb-4 text-text-secondary max-w-2xl">
        Generate placeholder text. Choose the unit (words, sentences, or
        paragraphs) and how many you want.
      </p>
      <div className="flex flex-col gap-3 max-w-2xl">
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
              <label key={opt.value} className="flex items-center gap-2 cursor-pointer">
                <input
                  type="radio"
                  name="lorem-type"
                  value={opt.value}
                  checked={type === opt.value}
                  onChange={() => {
                    setType(opt.value);
                    setCount((c) => Math.min(c, MAX[opt.value]));
                  }}
                  className="rounded border-border"
                />
                <span>{opt.label}</span>
              </label>
            ))}
          </div>
        </div>
        <div>
          <label htmlFor="lorem-count" className="font-medium block mb-2">
            Count (1–{maxCount})
          </label>
          <input
            id="lorem-count"
            type="number"
            min={1}
            max={maxCount}
            value={count}
            onChange={(e) => setCount(e.target.valueAsNumber || 1)}
            className="w-full max-w-[8rem] p-2.5 font-inherit text-base rounded-lg border border-border bg-bg text-text"
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
