import { createContext, useContext, useEffect, useLayoutEffect, useState } from 'react';

export type Theme = 'light' | 'dark' | 'system';

export type EffectiveTheme = 'light' | 'dark';

const STORAGE_KEY = 'digi-leatherman-theme';

function getStoredTheme(): Theme | null {
  try {
    const v = localStorage.getItem(STORAGE_KEY);
    if (v === 'light' || v === 'dark' || v === 'system') return v;
  } catch {
    /* ignore */
  }
  return null;
}

function getEffectiveThemeFromSystem(): EffectiveTheme {
  if (typeof window === 'undefined') return 'dark';
  return window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light';
}

function resolveEffectiveTheme(preference: Theme): EffectiveTheme {
  if (preference === 'light' || preference === 'dark') return preference;
  return getEffectiveThemeFromSystem();
}

type ThemeContextValue = {
  theme: Theme;
  effectiveTheme: EffectiveTheme;
  setTheme: (theme: Theme) => void;
};

const ThemeContext = createContext<ThemeContextValue | null>(null);

export function ThemeProvider({ children }: { children: React.ReactNode }) {
  const [theme, setThemeState] = useState<Theme>(() => getStoredTheme() ?? 'system');
  const [effectiveTheme, setEffectiveTheme] = useState<EffectiveTheme>(() =>
    resolveEffectiveTheme(theme)
  );

  useLayoutEffect(() => {
    const next = resolveEffectiveTheme(theme);
    setEffectiveTheme(next);
    document.documentElement.setAttribute('data-theme', next);
  }, [theme]);

  useEffect(() => {
    if (theme !== 'system') return;
    const mq = window.matchMedia('(prefers-color-scheme: dark)');
    const handler = () => {
      const next = getEffectiveThemeFromSystem();
      setEffectiveTheme(next);
      document.documentElement.setAttribute('data-theme', next);
    };
    mq.addEventListener('change', handler);
    return () => mq.removeEventListener('change', handler);
  }, [theme]);

  useEffect(() => {
    try {
      localStorage.setItem(STORAGE_KEY, theme);
    } catch {
      /* ignore */
    }
  }, [theme]);

  useEffect(() => {
    const stored = getStoredTheme();
    if (stored) setThemeState(stored);
  }, []);

  const setTheme = (next: Theme) => setThemeState(next);

  return (
    <ThemeContext.Provider value={{ theme, effectiveTheme, setTheme }}>
      {children}
    </ThemeContext.Provider>
  );
}

export function useTheme() {
  const ctx = useContext(ThemeContext);
  if (!ctx) throw new Error('useTheme must be used within ThemeProvider');
  return ctx;
}
