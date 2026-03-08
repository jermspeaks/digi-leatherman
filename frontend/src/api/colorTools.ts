const API_BASE = import.meta.env.VITE_API_URL ?? 'http://localhost:8100';

export interface ColorConvertResult {
  hex: string;
  rgb: string;
  hsl: string;
  oklch: string;
}

export interface ContrastResult {
  ratio: number;
  ratioText: string;
  aaNormal: boolean;
  aaLarge: boolean;
  aaaNormal: boolean;
  aaaLarge: boolean;
  color1Hex: string;
  color2Hex: string;
  color1Lighter: boolean;
}

export interface BlindnessResult {
  original: string;
  simulated: string;
  type: string;
}

export interface PaletteResult {
  colors: string[];
}

export interface ShadesResult {
  base: string;
  shades: string[];
  tints: string[];
}

export interface HarmoniesResult {
  base: string;
  complementary: string;
  triadic: string[];
  analogous: string[];
  splitComplementary: string[];
  tetradic: string[];
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

export async function convertColor(value: string): Promise<ColorConvertResult> {
  return postJSON('/api/color/convert', { value });
}

export async function checkContrast(color1: string, color2: string): Promise<ContrastResult> {
  return postJSON('/api/color/contrast', { color1, color2 });
}

export async function simulateBlindness(value: string, type: string): Promise<BlindnessResult> {
  return postJSON('/api/color/blindness', { value, type });
}

export async function extractPalette(image: string, count?: number): Promise<PaletteResult> {
  return postJSON('/api/color/palette', { image, count: count ?? 5 });
}

export async function generateShades(value: string, count?: number): Promise<ShadesResult> {
  return postJSON('/api/color/shades', { value, count: count ?? 5 });
}

export async function generateHarmonies(value: string): Promise<HarmoniesResult> {
  return postJSON('/api/color/harmonies', { value });
}
