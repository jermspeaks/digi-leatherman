import { useState, useEffect } from 'react';
import { NavLink, useLocation } from 'react-router-dom';
import { sidebarConfig, type SidebarItem } from '../config/sidebarConfig';
import './Sidebar.css';

type SidebarProps = {
  collapsed: boolean;
  onCollapsedChange: (collapsed: boolean) => void;
};

type ItemGroup = { key: string; label: string; items: SidebarItem[] };

function getItemGroups(items: SidebarItem[]): ItemGroup[] {
  const order: string[] = [];
  const groupMap = new Map<string, SidebarItem[]>();
  for (const item of items) {
    const key = item.subGroup ?? '_';
    if (!groupMap.has(key)) {
      order.push(key);
      groupMap.set(key, []);
    }
    groupMap.get(key)!.push(item);
  }
  return order.map((key) => ({
    key,
    label: key === '_' ? 'Other' : key,
    items: groupMap.get(key)!,
  }));
}

export function Sidebar({ collapsed, onCollapsedChange }: SidebarProps) {
  const location = useLocation();
  const pathname = location.pathname;

  const [expandedIds, setExpandedIds] = useState<Set<string>>(() => {
    const active = sidebarConfig.find((cat) =>
      cat.items.some((item) => item.path === pathname)
    );
    const ids = new Set<string>();
    if (active) {
      ids.add(active.id);
      const groups = getItemGroups(active.items);
      const activeGroup = groups.find((g) =>
        g.items.some((item) => item.path === pathname)
      );
      if (activeGroup && activeGroup.items.length > 1) {
        ids.add(`${active.id}--${activeGroup.key}`);
      }
    }
    return ids;
  });

  useEffect(() => {
    const active = sidebarConfig.find((cat) =>
      cat.items.some((item) => item.path === pathname)
    );
    if (active) {
      setExpandedIds((prev) => {
        const next = new Set(prev).add(active.id);
        const groups = getItemGroups(active.items);
        const activeGroup = groups.find((g) =>
          g.items.some((item) => item.path === pathname)
        );
        if (activeGroup && activeGroup.items.length > 1) {
          next.add(`${active.id}--${activeGroup.key}`);
        }
        return next;
      });
    }
  }, [pathname]);

  const toggleExpanded = (id: string) => {
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
          const groups = getItemGroups(category.items);

          return (
            <div key={category.id} className="sidebar-category">
              <button
                type="button"
                className="sidebar-category-heading"
                onClick={() => toggleExpanded(category.id)}
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
                  {groups.map((group) => {
                    if (group.items.length === 1) {
                      const item = group.items[0];
                      return (
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
                      );
                    }
                    const subId = `${category.id}--${group.key}`;
                    const isSubExpanded = expandedIds.has(subId);
                    return (
                      <li key={group.key} className="sidebar-subgroup">
                        <button
                          type="button"
                          className="sidebar-subgroup-heading"
                          onClick={() => toggleExpanded(subId)}
                          aria-expanded={isSubExpanded}
                          aria-controls={`sidebar-sub-${subId}`}
                          id={`sidebar-sub-heading-${subId}`}
                        >
                          <span className="sidebar-subgroup-chevron" aria-hidden>
                            {isSubExpanded ? '▼' : '▶'}
                          </span>
                          <span className="sidebar-subgroup-label">
                            {group.label}
                          </span>
                        </button>
                        <ul
                          id={`sidebar-sub-${subId}`}
                          className="sidebar-subgroup-items"
                          role="group"
                          aria-labelledby={`sidebar-sub-heading-${subId}`}
                          hidden={!isSubExpanded}
                        >
                          {group.items.map((item) => (
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
                      </li>
                    );
                  })}
                </ul>
              )}
            </div>
          );
        })}
      </nav>
      {!collapsed && (
        <button
          type="button"
          className="sidebar-collapse-toggle"
          onClick={() => onCollapsedChange(true)}
          aria-label="Collapse sidebar"
          title="Collapse sidebar"
        >
          «
        </button>
      )}
    </aside>
  );
}
