const API_BASE = import.meta.env.VITE_API_URL ?? 'http://localhost:8100';

export interface JsonResult {
  result: string;
}

export interface ValidateResult {
  valid: boolean;
  error?: string;
}

async function postJson(path: string, body: object): Promise<Response> {
  return fetch(`${API_BASE}${path}`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(body),
  });
}

export async function formatJson(value: string): Promise<JsonResult> {
  const res = await postJson('/api/json/format', { value });
  if (!res.ok) {
    const text = await res.text();
    throw new Error(text || `HTTP ${res.status}`);
  }
  return res.json();
}

export async function minifyJson(value: string): Promise<JsonResult> {
  const res = await postJson('/api/json/minify', { value });
  if (!res.ok) {
    const text = await res.text();
    throw new Error(text || `HTTP ${res.status}`);
  }
  return res.json();
}

export async function validateJson(value: string): Promise<ValidateResult> {
  const res = await postJson('/api/json/validate', { value });
  if (!res.ok) {
    const text = await res.text();
    throw new Error(text || `HTTP ${res.status}`);
  }
  return res.json();
}

export async function pathQueryJson(value: string, path: string): Promise<JsonResult> {
  const res = await postJson('/api/json/path', { value, path });
  if (!res.ok) {
    const text = await res.text();
    throw new Error(text || `HTTP ${res.status}`);
  }
  return res.json();
}

export async function diffJson(valueA: string, valueB: string): Promise<JsonResult> {
  const res = await postJson('/api/json/diff', { valueA, valueB });
  if (!res.ok) {
    const text = await res.text();
    throw new Error(text || `HTTP ${res.status}`);
  }
  return res.json();
}
