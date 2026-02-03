import { useTheme } from '../context/ThemeContext';
import { Card, CardContent, CardHeader, CardTitle } from './ui/card';
import { Button } from './ui/button';
export function Settings() {
  const { theme, effectiveTheme, setTheme } = useTheme();

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
    </div>
  );
}
