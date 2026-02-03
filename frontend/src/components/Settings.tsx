import { useEffect, useState } from 'react';
import { useTheme } from '../context/ThemeContext';
import { Card, CardContent, CardHeader, CardTitle } from './ui/card';
import { Button } from './ui/button';

const COLOR_TOKENS: { name: string; variable: string }[] = [
  { name: 'Text', variable: '--color-text' },
  { name: 'Text muted', variable: '--color-text-muted' },
  { name: 'Background', variable: '--color-bg' },
  { name: 'Background elevated', variable: '--color-bg-elevated' },
  { name: 'Border', variable: '--color-border' },
  { name: 'Accent', variable: '--color-accent' },
  { name: 'Accent hover', variable: '--color-accent-hover' },
  { name: 'Sidebar bg', variable: '--color-sidebar-bg' },
  { name: 'Sidebar active', variable: '--color-sidebar-active' },
  { name: 'Navbar bg', variable: '--color-navbar-bg' },
  { name: 'Primary', variable: '--color-primary' },
  { name: 'Muted', variable: '--color-muted' },
  { name: 'Destructive', variable: '--color-destructive' },
];

const SPACING_TOKENS = [
  { name: '1', value: '0.25rem' },
  { name: '2', value: '0.5rem' },
  { name: '4', value: '1rem' },
  { name: '6', value: '1.5rem' },
];

function getComputedColorValue(cssVar: string): string {
  if (typeof document === 'undefined') return '';
  const raw = getComputedStyle(document.documentElement).getPropertyValue(cssVar).trim();
  return raw || '';
}

export function Settings() {
  const { theme, effectiveTheme, setTheme } = useTheme();
  const [colorValues, setColorValues] = useState<Record<string, string>>({});

  useEffect(() => {
    const next: Record<string, string> = {};
    COLOR_TOKENS.forEach(({ variable }) => {
      next[variable] = getComputedColorValue(variable);
    });
    setColorValues(next);
  }, [effectiveTheme]);

  return (
    <div className="space-y-6">
      <section aria-labelledby="settings-theme-heading">
        <Card>
          <CardHeader>
            <CardTitle id="settings-theme-heading">Theme</CardTitle>
          </CardHeader>
          <CardContent className="flex flex-wrap gap-2">
            <Button
              type="button"
              variant={theme === 'dark' ? 'default' : 'outline'}
              size="sm"
              onClick={() => setTheme('dark')}
              aria-pressed={theme === 'dark'}
            >
              Dark
            </Button>
            <Button
              type="button"
              variant={theme === 'light' ? 'default' : 'outline'}
              size="sm"
              onClick={() => setTheme('light')}
              aria-pressed={theme === 'light'}
            >
              Light
            </Button>
            <Button
              type="button"
              variant={theme === 'system' ? 'default' : 'outline'}
              size="sm"
              onClick={() => setTheme('system')}
              aria-pressed={theme === 'system'}
            >
              System
            </Button>
            {theme === 'system' && (
              <span className="inline-flex items-center text-sm text-text-muted ml-1">
                Following system ({effectiveTheme})
              </span>
            )}
          </CardContent>
        </Card>
      </section>

      <section aria-labelledby="colophon-heading">
        <Card>
          <CardHeader>
            <CardTitle id="colophon-heading">Colophon</CardTitle>
          </CardHeader>
          <CardContent className="space-y-8">
            <div>
              <h3 className="text-sm font-medium mb-3">Colors</h3>
              <div className="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 gap-3">
                {COLOR_TOKENS.map(({ name, variable }) => (
                  <div key={variable} className="flex flex-col gap-1.5">
                    <div
                      className="h-10 w-full rounded border border-border"
                      style={{ backgroundColor: `var(${variable})` }}
                      aria-hidden
                    />
                    <span className="text-xs font-medium">{name}</span>
                    {colorValues[variable] && (
                      <code className="text-xs text-text-muted truncate" title={colorValues[variable]}>
                        {colorValues[variable]}
                      </code>
                    )}
                  </div>
                ))}
              </div>
            </div>

            <div>
              <h3 className="text-sm font-medium mb-3">Typography</h3>
              <p className="text-sm text-text-muted mb-2">
                Font: system-ui, Avenir, Helvetica, Arial, sans-serif
              </p>
              <p className="text-sm text-text-muted mb-2">Body: 400, line-height 1.5</p>
              <p className="mb-3">The quick brown fox jumps over the lazy dog.</p>
              <p className="text-sm text-text-muted mb-2">Heading 1: 3.2em, line-height 1.1</p>
              <h1 className="mb-0">Sample heading</h1>
            </div>

            <div>
              <h3 className="text-sm font-medium mb-3">Border radius</h3>
              <p className="text-sm text-text-muted mb-2">0.625rem (var(--radius))</p>
              <div
                className="h-12 w-24 border-2 border-border bg-bg-elevated"
                style={{ borderRadius: 'var(--radius)' }}
                aria-hidden
              />
            </div>

            <div>
              <h3 className="text-sm font-medium mb-3">Spacing</h3>
              <div className="flex flex-wrap items-end gap-6">
                {SPACING_TOKENS.map(({ name, value }) => (
                  <div key={name} className="flex flex-col items-center gap-1.5">
                    <div
                      className="w-4 bg-accent"
                      style={{ height: value }}
                      aria-hidden
                    />
                    <span className="text-xs">
                      {name} ({value})
                    </span>
                  </div>
                ))}
              </div>
            </div>
          </CardContent>
        </Card>
      </section>

      <p className="text-sm text-text-muted">
        Digi Leatherman — Built with React, Vite, Tailwind CSS, shadcn/ui; backend in Go.
      </p>
    </div>
  );
}
