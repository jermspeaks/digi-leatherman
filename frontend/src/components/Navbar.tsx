import { Link } from 'react-router-dom';
import { useTheme } from '../context/ThemeContext';

export function Navbar() {
  const { theme, setTheme } = useTheme();

  const toggleTheme = () => {
    setTheme(theme === 'dark' ? 'light' : 'dark');
  };

  return (
    <header
      className="sticky top-0 z-[100] h-14 bg-navbar-bg border-b border-navbar-border"
      role="banner"
    >
      <div className="flex items-center justify-between h-full max-w-full px-4">
        <Link
          to="/"
          className="text-xl font-semibold text-text no-underline hover:text-accent"
          aria-label="Digi Leatherman home"
        >
          Digi Leatherman
        </Link>
        <div className="flex items-center gap-2">
          <button
            type="button"
            className="p-2 border-none bg-transparent text-2xl leading-none hover:bg-sidebar-active hover:border-transparent"
            onClick={toggleTheme}
            aria-label={theme === 'dark' ? 'Switch to light mode' : 'Switch to dark mode'}
            title={theme === 'dark' ? 'Light mode' : 'Dark mode'}
          >
            {theme === 'dark' ? (
              <span className="inline-block" aria-hidden>☀</span>
            ) : (
              <span className="inline-block" aria-hidden>☾</span>
            )}
          </button>
        </div>
      </div>
    </header>
  );
}
