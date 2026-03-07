const API_BASE = import.meta.env.VITE_API_BASE ?? 'http://localhost:8100';

export interface DateParseResult {
  iso: string;
  unix: number;
  unixMs: number;
  rfc2822: string;
  relative: string;
}

export interface TimezoneResult {
  result: string;
  fromTime: string;
  toTime: string;
}

export interface CountdownResult {
  days: number;
  hours: number;
  minutes: number;
  seconds: number;
  isPast: boolean;
  text: string;
}

async function postJSON<T>(path: string, body: object): Promise<T> {
  const res = await fetch(`${API_BASE}${path}`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(body),
  });
  if (!res.ok) {
    const text = await res.text();
    throw new Error(text || `HTTP ${res.status}`);
  }
  return res.json();
}

export async function dateParse(value: string): Promise<DateParseResult> {
  return postJSON<DateParseResult>('/api/date/parse', { value });
}

export async function dateTimezone(
  value: string,
  fromZone: string,
  toZone: string
): Promise<TimezoneResult> {
  return postJSON<TimezoneResult>('/api/date/timezone', { value, fromZone, toZone });
}

export async function dateCountdown(value: string): Promise<CountdownResult> {
  return postJSON<CountdownResult>('/api/date/countdown', { value });
}
