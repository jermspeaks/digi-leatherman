import { useState, useEffect } from 'react';
import { NavLink, useLocation } from 'react-router-dom';
import { sidebarConfig } from '../config/sidebarConfig';
import './Sidebar.css';

type SidebarProps = {
  collapsed: boolean;
  onCollapsedChange: (collapsed: boolean) => void;
};

export function Sidebar({ collapsed, onCollapsedChange }: SidebarProps) {
  const location = useLocation();
  const pathname = location.pathname;

  const [expandedIds, setExpandedIds] = useState<Set<string>>(() => {
    const active = sidebarConfig.find((cat) =>
      cat.items.some((item) => item.path === pathname)
    );
    const ids = new Set<string>();
    if (active) ids.add(active.id);
    return ids;
  });

  useEffect(() => {
    const active = sidebarConfig.find((cat) =>
      cat.items.some((item) => item.path === pathname)
    );
    if (active) {
      setExpandedIds((prev) => new Set(prev).add(active.id));
    }
  }, [pathname]);

  const toggleCategory = (id: string) => {
    setExpandedIds((prev) => {
      const next = new Set(prev);
      if (next.has(id)) next.delete(id);
      else next.add(id);
      return next;
    });
  };

  return (
    <aside
      className={`sidebar ${collapsed ? 'sidebar--collapsed' : ''}`}
      aria-label="Tool navigation"
    >
      <nav className="sidebar-nav">
        {sidebarConfig.map((category) => {
          const isExpanded = expandedIds.has(category.id);
          const hasItems = category.items.length > 0;

          return (
            <div key={category.id} className="sidebar-category">
              <button
                type="button"
                className="sidebar-category-heading"
                onClick={() => toggleCategory(category.id)}
                aria-expanded={hasItems ? isExpanded : undefined}
                aria-controls={`sidebar-cat-${category.id}`}
                id={`sidebar-cat-heading-${category.id}`}
              >
                <span className="sidebar-category-chevron" aria-hidden>
                  {hasItems ? (isExpanded ? '▼' : '▶') : '○'}
                </span>
                <span className="sidebar-category-label">{category.label}</span>
              </button>
              {hasItems && (
                <ul
                  id={`sidebar-cat-${category.id}`}
                  className="sidebar-category-items"
                  role="group"
                  aria-labelledby={`sidebar-cat-heading-${category.id}`}
                  hidden={!isExpanded}
                >
                  {category.items.map((item) => (
                      <li key={item.id}>
                        <NavLink
                          to={item.path}
                          className={({ isActive: active }) =>
                            `sidebar-link ${active ? 'sidebar-link--active' : ''}`
                          }
                          end={false}
                        >
                          {item.label}
                        </NavLink>
                      </li>
                  ))}
                </ul>
              )}
            </div>
          );
        })}
      </nav>
    </aside>
  );
}
