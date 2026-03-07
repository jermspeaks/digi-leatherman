import { useRef, useState } from 'react';
import { Copy } from 'lucide-react';
import { Button } from '@/components/ui/button';
import {
  dateParse,
  dateTimezone,
  dateCountdown,
  type DateParseResult,
  type TimezoneResult,
  type CountdownResult,
} from '../api/dateTools';

export type DateToolId = 'parse' | 'timezone' | 'countdown';

const COMMON_TIMEZONES = [
  'UTC',
  'America/New_York',
  'America/Chicago',
  'America/Denver',
  'America/Los_Angeles',
  'America/Anchorage',
  'Pacific/Honolulu',
  'Europe/London',
  'Europe/Paris',
  'Europe/Berlin',
  'Europe/Moscow',
  'Asia/Dubai',
  'Asia/Kolkata',
  'Asia/Singapore',
  'Asia/Shanghai',
  'Asia/Tokyo',
  'Australia/Sydney',
  'Pacific/Auckland',
];

type ToolMeta = {
  id: DateToolId;
  label: string;
  description: string;
  example: string;
  placeholder: string;
  buttonLabel: string;
};

const TOOL_CONFIG: ToolMeta[] = [
  {
    id: 'parse',
    label: 'Parse & Format',
    description: 'Parse any date format and convert to ISO 8601, Unix timestamp, RFC 2822, and relative time.',
    example: '2024-03-15T10:30:00Z, 1710498600, Mar 15, 2024',
    placeholder: 'Enter a date in any format...',
    buttonLabel: 'Parse',
  },
  {
    id: 'timezone',
    label: 'Timezone Convert',
    description: 'Convert a date/time from one timezone to another.',
    example: '2024-03-15 12:00:00 from UTC to America/New_York',
    placeholder: 'Enter a date/time to convert...',
    buttonLabel: 'Convert',
  },
  {
    id: 'countdown',
    label: 'Countdown',
    description: 'Calculate the time remaining until (or elapsed since) a target date.',
    example: '2025-01-01T00:00:00Z',
    placeholder: 'Enter a target date...',
    buttonLabel: 'Calculate',
  },
];

export const TOOL_DESCRIPTIONS: Record<string, string> = Object.fromEntries(
  TOOL_CONFIG.map((c) => [c.id, c.description])
);

const TOOL_MAP = Object.fromEntries(TOOL_CONFIG.map((c) => [c.id, c])) as Record<DateToolId, ToolMeta>;

type DateToolsProps = {
  tool?: DateToolId;
};

export function DateTools({ tool: toolProp }: DateToolsProps) {
  const [tab, setTab] = useState<DateToolId>(toolProp ?? 'parse');
  const [input, setInput] = useState('');
  const [fromZone, setFromZone] = useState('UTC');
  const [toZone, setToZone] = useState('America/New_York');
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);
  const [copied, setCopied] = useState<string | null>(null);
  const copyTimeoutRef = useRef<ReturnType<typeof setTimeout> | null>(null);

  const [parseResult, setParseResult] = useState<DateParseResult | null>(null);
  const [timezoneResult, setTimezoneResult] = useState<TimezoneResult | null>(null);
  const [countdownResult, setCountdownResult] = useState<CountdownResult | null>(null);

  const tool = toolProp ?? tab;
  const showTabs = toolProp == null;
  const config = TOOL_MAP[tool];

  const clearResults = () => {
    setParseResult(null);
    setTimezoneResult(null);
    setCountdownResult(null);
  };

  const run = async () => {
    setError(null);
    clearResults();
    if (copyTimeoutRef.current) {
      clearTimeout(copyTimeoutRef.current);
      copyTimeoutRef.current = null;
    }
    setCopied(null);
    if (!input.trim()) return;
    setLoading(true);
    try {
      switch (tool) {
        case 'parse': {
          const result = await dateParse(input);
          setParseResult(result);
          break;
        }
        case 'timezone': {
          const result = await dateTimezone(input, fromZone, toZone);
          setTimezoneResult(result);
          break;
        }
        case 'countdown': {
          const result = await dateCountdown(input);
          setCountdownResult(result);
          break;
        }
      }
    } catch (e) {
      setError(e instanceof Error ? e.message : 'Request failed');
    } finally {
      setLoading(false);
    }
  };

  const useExample = () => {
    switch (tool) {
      case 'parse':
        setInput('2024-03-15T10:30:00Z');
        break;
      case 'timezone':
        setInput('2024-03-15 12:00:00');
        setFromZone('UTC');
        setToZone('America/New_York');
        break;
      case 'countdown':
        setInput('2025-12-31T23:59:59Z');
        break;
    }
  };

  const handleCopy = async (value: string, label: string) => {
    if (!value) return;
    if (copyTimeoutRef.current) clearTimeout(copyTimeoutRef.current);
    try {
      await navigator.clipboard.writeText(value);
      setCopied(label);
      copyTimeoutRef.current = setTimeout(() => {
        copyTimeoutRef.current = null;
        setCopied(null);
      }, 1500);
    } catch {
      // Permission denied
    }
  };

  const OutputRow = ({ label, value }: { label: string; value: string }) => (
    <div className="flex items-center justify-between gap-2 py-2 border-b border-border last:border-b-0">
      <div className="flex flex-col gap-0.5 min-w-0 flex-1">
        <span className="text-sm text-text-secondary">{label}</span>
        <code className="text-text break-all">{value}</code>
      </div>
      <Button
        type="button"
        variant="outline"
        size="sm"
        onClick={() => handleCopy(value, label)}
        className="shrink-0"
        aria-label={copied === label ? 'Copied' : `Copy ${label}`}
      >
        <Copy className="size-3.5" aria-hidden />
        {copied === label ? 'Copied!' : 'Copy'}
      </Button>
    </div>
  );

  return (
    <section className="text-left mt-6">
      <h2 className="mb-3 text-2xl text-text">Dates</h2>
      {showTabs && (
        <nav className="flex gap-2 mb-4" aria-label="Date tools">
          {TOOL_CONFIG.map((c) => (
            <button
              key={c.id}
              type="button"
              className={`py-2 px-4 ${tab === c.id ? 'border-accent bg-sidebar-active' : ''}`}
              onClick={() => {
                setTab(c.id);
                clearResults();
                setError(null);
              }}
            >
              {c.label}
            </button>
          ))}
        </nav>
      )}
      <p className="mb-2 text-text-secondary max-w-2xl">{config.description}</p>
      <p className="mb-4 text-sm text-text-secondary max-w-2xl">
        Example: <code className="px-1 rounded bg-bg-elevated">{config.example}</code>
        <button type="button" onClick={useExample} className="ml-2 text-accent hover:underline">
          Use example
        </button>
      </p>
      <div className="flex flex-col gap-3 max-w-2xl">
        <label htmlFor="date-input" className="font-medium">
          Input
        </label>
        <textarea
          id="date-input"
          className="w-full min-h-16 p-2.5 font-inherit text-base rounded-lg border border-border bg-bg text-text resize-y"
          value={input}
          onChange={(e) => setInput(e.target.value)}
          placeholder={config.placeholder}
          rows={2}
        />

        {tool === 'timezone' && (
          <div className="flex flex-wrap gap-4">
            <div className="flex flex-col gap-1.5 flex-1 min-w-[140px]">
              <label htmlFor="from-zone" className="text-sm font-medium">
                From Timezone
              </label>
              <select
                id="from-zone"
                value={fromZone}
                onChange={(e) => setFromZone(e.target.value)}
                className="p-2 rounded-lg border border-border bg-bg text-text"
              >
                {COMMON_TIMEZONES.map((tz) => (
                  <option key={tz} value={tz}>
                    {tz}
                  </option>
                ))}
              </select>
            </div>
            <div className="flex flex-col gap-1.5 flex-1 min-w-[140px]">
              <label htmlFor="to-zone" className="text-sm font-medium">
                To Timezone
              </label>
              <select
                id="to-zone"
                value={toZone}
                onChange={(e) => setToZone(e.target.value)}
                className="p-2 rounded-lg border border-border bg-bg text-text"
              >
                {COMMON_TIMEZONES.map((tz) => (
                  <option key={tz} value={tz}>
                    {tz}
                  </option>
                ))}
              </select>
            </div>
          </div>
        )}

        <button type="button" onClick={run} disabled={loading}>
          {loading ? '…' : config.buttonLabel}
        </button>

        {error && (
          <p className="text-red-500 m-0" role="alert">
            {error}
          </p>
        )}

        {parseResult && (
          <div className="mt-2">
            <label className="font-medium mb-2 block">Output</label>
            <div className="bg-bg-elevated rounded-lg p-3">
              <OutputRow label="ISO 8601" value={parseResult.iso} />
              <OutputRow label="Unix (seconds)" value={String(parseResult.unix)} />
              <OutputRow label="Unix (milliseconds)" value={String(parseResult.unixMs)} />
              <OutputRow label="RFC 2822" value={parseResult.rfc2822} />
              <OutputRow label="Relative" value={parseResult.relative} />
            </div>
          </div>
        )}

        {timezoneResult && (
          <div className="mt-2">
            <label className="font-medium mb-2 block">Output</label>
            <div className="bg-bg-elevated rounded-lg p-3">
              <OutputRow label={`From (${fromZone})`} value={timezoneResult.fromTime} />
              <OutputRow label={`To (${toZone})`} value={timezoneResult.toTime} />
            </div>
          </div>
        )}

        {countdownResult && (
          <div className="mt-2">
            <label className="font-medium mb-2 block">Output</label>
            <div className="bg-bg-elevated rounded-lg p-3">
              <div className="grid grid-cols-4 gap-4 text-center mb-4">
                <div className="flex flex-col">
                  <span className="text-3xl font-bold text-text">{countdownResult.days}</span>
                  <span className="text-sm text-text-secondary">days</span>
                </div>
                <div className="flex flex-col">
                  <span className="text-3xl font-bold text-text">{countdownResult.hours}</span>
                  <span className="text-sm text-text-secondary">hours</span>
                </div>
                <div className="flex flex-col">
                  <span className="text-3xl font-bold text-text">{countdownResult.minutes}</span>
                  <span className="text-sm text-text-secondary">minutes</span>
                </div>
                <div className="flex flex-col">
                  <span className="text-3xl font-bold text-text">{countdownResult.seconds}</span>
                  <span className="text-sm text-text-secondary">seconds</span>
                </div>
              </div>
              <div className="text-center pt-2 border-t border-border">
                <span className={`text-sm font-medium ${countdownResult.isPast ? 'text-text-secondary' : 'text-accent'}`}>
                  {countdownResult.isPast ? 'Time elapsed' : 'Time remaining'}
                </span>
              </div>
              <OutputRow label="Summary" value={countdownResult.text} />
            </div>
          </div>
        )}
      </div>
    </section>
  );
}
