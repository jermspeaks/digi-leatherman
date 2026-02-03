import { Link, NavLink, useLocation } from "react-router-dom";
import { SearchIcon } from "lucide-react";

const navLinkBase =
  "px-3 py-2 text-sm font-medium no-underline rounded-none border-b-2 border-transparent hover:text-accent";
const navLinkActive = "text-accent border-accent";

type NavbarProps = {
  onOpenCommandPalette?: () => void;
};

export function Navbar({ onOpenCommandPalette }: NavbarProps) {
  const location = useLocation();
  const pathname = location.pathname;
  const isTools = pathname.startsWith("/tools");
  const isSettings = pathname === "/settings";

  return (
    <header
      className="sticky top-0 z-[100] h-14 bg-navbar-bg border-b border-navbar-border"
      role="banner"
    >
      <div className="flex items-center justify-between h-full max-w-full px-4 gap-4">
        <div className="flex items-center gap-6 min-w-0">
          <Link
            to="/"
            className="text-xl font-semibold text-text no-underline hover:text-accent shrink-0"
            aria-label="Digi Leatherman home"
          >
            Digi Leatherman
          </Link>
          {onOpenCommandPalette && (
            <button
              type="button"
              onClick={onOpenCommandPalette}
              className="flex items-center gap-2 px-3 py-1.5 w-full max-w-[200px] text-left text-sm text-muted-foreground bg-muted/50 border border-input rounded-md hover:bg-muted hover:text-foreground transition-colors"
              aria-label="Open command palette"
            >
              <SearchIcon className="size-4 shrink-0 opacity-50" />
              <span className="truncate">Search tools…</span>
              <kbd className="ml-auto hidden sm:inline-flex h-5 items-center gap-0.5 rounded border border-border bg-muted px-1 font-mono text-[10px] font-medium text-muted-foreground">
                ⌘K
              </kbd>
            </button>
          )}
        </div>
        <nav className="flex items-center gap-1 shrink-0" aria-label="Primary">
          <NavLink
            to="/tools/string/url-encode"
            className={({ isActive }) =>
              `${navLinkBase} ${isActive || isTools ? navLinkActive : "text-text-muted"}`
            }
            end={false}
          >
            Tools
          </NavLink>
          <NavLink
            to="/settings"
            className={({ isActive }) =>
              `${navLinkBase} ${isActive || isSettings ? navLinkActive : "text-text-muted"}`
            }
            end
          >
            Settings
          </NavLink>
        </nav>
      </div>
    </header>
  );
}
