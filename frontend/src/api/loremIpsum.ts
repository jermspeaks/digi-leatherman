const API_BASE = import.meta.env.VITE_API_URL ?? 'http://localhost:8100';

export interface StringResult {
  result: string;
}

export type LoremType = 'words' | 'sentences' | 'paragraphs';

export type LoremToolId =
  | 'generator'
  | 'characters'
  | 'bytes'
  | 'title'
  | 'slug'
  | 'camelCase'
  | 'list'
  | 'headings'
  | 'html'
  | 'markdown'
  | 'json';

export interface LoremOptions {
  startWithClassic?: boolean;
  type?: LoremType;
  vocabulary?: 'default' | 'bacon' | 'hipster';
  wholeWordsOnly?: boolean;
  listStyle?: 'bullet' | 'numbered';
  headingStyle?: 'plain' | 'markdown';
  format?: 'paragraphs' | 'list' | 'headings';
  keys?: string[];
}

export interface LoremRequest {
  tool: LoremToolId;
  count: number;
  options?: LoremOptions;
}

export async function generate(request: LoremRequest): Promise<StringResult> {
  const res = await fetch(`${API_BASE}/api/lorem-ipsum/generate`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(request),
  });
  if (!res.ok) {
    const text = await res.text();
    throw new Error(text || `HTTP ${res.status}`);
  }
  return res.json();
}

/** Legacy: generate by type + count (maps to generator tool). */
export async function generateLegacy(
  type: LoremType,
  count: number
): Promise<StringResult> {
  return generate({
    tool: 'generator',
    count,
    options: { type },
  });
}
