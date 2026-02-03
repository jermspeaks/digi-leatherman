const API_BASE = import.meta.env.VITE_API_URL ?? 'http://localhost:8080';

export interface StringResult {
  result: string;
}

export async function urlEncode(value: string): Promise<StringResult> {
  const res = await fetch(`${API_BASE}/api/string/url-encode`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ value }),
  });
  if (!res.ok) {
    const text = await res.text();
    throw new Error(text || `HTTP ${res.status}`);
  }
  return res.json();
}

export async function urlDecode(value: string): Promise<StringResult> {
  const res = await fetch(`${API_BASE}/api/string/url-decode`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ value }),
  });
  if (!res.ok) {
    const text = await res.text();
    throw new Error(text || `HTTP ${res.status}`);
  }
  return res.json();
}
