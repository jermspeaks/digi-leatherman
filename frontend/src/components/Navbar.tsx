import { Link, NavLink, useLocation } from "react-router-dom";

const navLinkBase =
  "px-3 py-2 text-sm font-medium no-underline rounded-none border-b-2 border-transparent hover:text-accent";
const navLinkActive = "text-accent border-accent";

export function Navbar() {
  const location = useLocation();
  const pathname = location.pathname;
  const isTools = pathname.startsWith("/tools");
  const isSettings = pathname === "/settings";

  return (
    <header
      className="sticky top-0 z-[100] h-14 bg-navbar-bg border-b border-navbar-border"
      role="banner"
    >
      <div className="flex items-center justify-between h-full max-w-full px-4">
        <div className="flex items-center gap-6">
          <Link
            to="/"
            className="text-xl font-semibold text-text no-underline hover:text-accent"
            aria-label="Digi Leatherman home"
          >
            Digi Leatherman
          </Link>
        </div>
        <nav className="flex items-center gap-1" aria-label="Primary">
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
