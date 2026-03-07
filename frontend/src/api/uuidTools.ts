const API_BASE = import.meta.env.VITE_API_URL ?? 'http://localhost:8100';

export interface UUIDResult {
  result: string;
}

export interface BulkGenerateResult {
  results: string[];
}

export interface UUIDValidateResult {
  valid: boolean;
  version?: number;
  variant?: string;
  error?: string;
}

export interface UUIDParseResult {
  valid: boolean;
  version?: number;
  variant?: string;
  timestamp?: string | null;
  clockSeq?: number | null;
  node?: string | null;
  error?: string;
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

export async function generateUUIDv4(): Promise<UUIDResult> {
  return postJSON<UUIDResult>('/api/uuid/generate-v4', {});
}

export async function generateUUIDv7(): Promise<UUIDResult> {
  return postJSON<UUIDResult>('/api/uuid/generate-v7', {});
}

export async function bulkGenerateUUID(version: number, count: number): Promise<BulkGenerateResult> {
  return postJSON<BulkGenerateResult>('/api/uuid/bulk-generate', { version, count });
}

export async function validateUUID(value: string): Promise<UUIDValidateResult> {
  return postJSON<UUIDValidateResult>('/api/uuid/validate', { value });
}

export async function parseUUID(value: string): Promise<UUIDParseResult> {
  return postJSON<UUIDParseResult>('/api/uuid/parse', { value });
}

export async function formatUUID(value: string, format: string): Promise<UUIDResult> {
  return postJSON<UUIDResult>('/api/uuid/format', { value, format });
}
