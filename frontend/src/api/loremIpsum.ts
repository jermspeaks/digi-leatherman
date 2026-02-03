const API_BASE = import.meta.env.VITE_API_URL ?? 'http://localhost:8100';

export interface StringResult {
  result: string;
}

export type LoremType = 'words' | 'sentences' | 'paragraphs';

export async function generate(
  type: LoremType,
  count: number
): Promise<StringResult> {
  const res = await fetch(`${API_BASE}/api/lorem-ipsum/generate`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ type, count }),
  });
  if (!res.ok) {
    const text = await res.text();
    throw new Error(text || `HTTP ${res.status}`);
  }
  return res.json();
}
