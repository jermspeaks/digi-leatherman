const API_BASE = import.meta.env.VITE_API_URL ?? 'http://localhost:8100';

export interface StringResult {
  result: string;
}

async function postString(path: string, value: string): Promise<StringResult> {
  const res = await fetch(`${API_BASE}${path}`, {
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

export async function urlEncode(value: string): Promise<StringResult> {
  return postString('/api/string/url-encode', value);
}

export async function urlDecode(value: string): Promise<StringResult> {
  return postString('/api/string/url-decode', value);
}

export async function base64Encode(value: string): Promise<StringResult> {
  return postString('/api/string/base64-encode', value);
}

export async function base64Decode(value: string): Promise<StringResult> {
  return postString('/api/string/base64-decode', value);
}

export async function trim(value: string): Promise<StringResult> {
  return postString('/api/string/trim', value);
}

export async function upperCase(value: string): Promise<StringResult> {
  return postString('/api/string/upper-case', value);
}

export async function lowerCase(value: string): Promise<StringResult> {
  return postString('/api/string/lower-case', value);
}

export async function capitalCase(value: string): Promise<StringResult> {
  return postString('/api/string/capital-case', value);
}

export async function snakeCase(value: string): Promise<StringResult> {
  return postString('/api/string/snake-case', value);
}

export async function kebabCase(value: string): Promise<StringResult> {
  return postString('/api/string/kebab-case', value);
}

export async function camelCase(value: string): Promise<StringResult> {
  return postString('/api/string/camel-case', value);
}

export async function pascalCase(value: string): Promise<StringResult> {
  return postString('/api/string/pascal-case', value);
}

export async function sentenceCase(value: string): Promise<StringResult> {
  return postString('/api/string/sentence-case', value);
}

export async function parseUrlParams(value: string): Promise<StringResult> {
  return postString('/api/string/url-parse-params', value);
}

export async function createUrlWithParams(value: string): Promise<StringResult> {
  return postString('/api/string/url-param-creator', value);
}

export async function spellOut(value: string, alphabet?: string): Promise<StringResult> {
  const res = await fetch(`${API_BASE}/api/string/spell-out`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ value, alphabet: alphabet ?? 'nato' }),
  });
  if (!res.ok) {
    const text = await res.text();
    throw new Error(text || `HTTP ${res.status}`);
  }
  return res.json();
}
