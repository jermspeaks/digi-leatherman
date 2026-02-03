import { Link } from 'react-router-dom';
import { useTheme } from '../context/ThemeContext';
import './Navbar.css';

export function Navbar() {
  const { theme, setTheme } = useTheme();

  const toggleTheme = () => {
    setTheme(theme === 'dark' ? 'light' : 'dark');
  };

  return (
    <header className="navbar" role="banner">
      <div className="navbar-inner">
        <Link to="/" className="navbar-brand" aria-label="Digi Leatherman home">
          Digi Leatherman
        </Link>
        <div className="navbar-actions">
          <button
            type="button"
            className="navbar-theme-toggle"
            onClick={toggleTheme}
            aria-label={theme === 'dark' ? 'Switch to light mode' : 'Switch to dark mode'}
            title={theme === 'dark' ? 'Light mode' : 'Dark mode'}
          >
            {theme === 'dark' ? (
              <span className="theme-icon" aria-hidden>☀</span>
            ) : (
              <span className="theme-icon" aria-hidden>☾</span>
            )}
          </button>
        </div>
      </div>
    </header>
  );
}
