import { useState, useEffect } from 'react';
import { NavLink, useLocation } from 'react-router-dom';
import { sidebarConfig, type SidebarItem } from '../config/sidebarConfig';

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

const sidebarLinkBase =
  'block py-[0.4rem] px-3 text-sm font-medium text-text-muted no-underline border-l-[3px] border-transparent -ml-6 pl-[2.0625rem] hover:text-text';
const sidebarLinkSubgroup =
  'block py-[0.4rem] px-3 text-sm font-medium text-text-muted no-underline border-l-[3px] border-transparent ml-0 pl-3 hover:text-text';
const sidebarLinkActive = 'text-accent bg-sidebar-active border-transparent';
const sidebarLinkActiveCollapsed = 'border-l-sidebar-active-bar';

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
      className={`w-64 min-w-64 shrink-0 flex flex-col bg-sidebar-bg border-r border-border transition-[width,min-width] duration-200 ${
        collapsed ? 'w-0 min-w-0 overflow-hidden border-r-0' : ''
      }`}
      aria-label="Tool navigation"
    >
      <nav className="flex-1 overflow-y-auto py-3">
        {sidebarConfig.map((category) => {
          const isExpanded = expandedIds.has(category.id);
          const hasItems = category.items.length > 0;
          const groups = getItemGroups(category.items);

          return (
            <div key={category.id} className="mb-1">
              <button
                type="button"
                className="flex items-center gap-2 w-full py-2 px-4 text-left text-[0.9rem] font-semibold text-text bg-transparent border-none cursor-pointer rounded-none hover:bg-sidebar-active"
                onClick={() => toggleExpanded(category.id)}
                aria-expanded={hasItems ? isExpanded : undefined}
                aria-controls={`sidebar-cat-${category.id}`}
                id={`sidebar-cat-heading-${category.id}`}
              >
                <span className="text-[0.65rem] text-text-muted shrink-0" aria-hidden>
                  {hasItems ? (isExpanded ? '▼' : '▶') : '○'}
                </span>
                <span className="flex-1">{category.label}</span>
              </button>
              {hasItems && (
                <ul
                  id={`sidebar-cat-${category.id}`}
                  className="list-none m-0 pl-6"
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
                              `${sidebarLinkBase} ${active ? sidebarLinkActive : ''} ${active && collapsed ? sidebarLinkActiveCollapsed : ''}`
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
                      <li key={group.key} className="list-none m-0">
                        <button
                          type="button"
                          className="flex items-center gap-[0.35rem] w-full py-[0.35rem] px-3 pl-6 text-left text-[0.8rem] font-medium text-text-muted bg-transparent border-none cursor-pointer rounded-none hover:text-text hover:bg-sidebar-active"
                          onClick={() => toggleExpanded(subId)}
                          aria-expanded={isSubExpanded}
                          aria-controls={`sidebar-sub-${subId}`}
                          id={`sidebar-sub-heading-${subId}`}
                        >
                          <span className="text-[0.6rem] text-text-muted shrink-0" aria-hidden>
                            {isSubExpanded ? '▼' : '▶'}
                          </span>
                          <span className="flex-1">{group.label}</span>
                        </button>
                        <ul
                          id={`sidebar-sub-${subId}`}
                          className="list-none m-0 pl-8"
                          role="group"
                          aria-labelledby={`sidebar-sub-heading-${subId}`}
                          hidden={!isSubExpanded}
                        >
                          {group.items.map((item) => (
                            <li key={item.id}>
                              <NavLink
                                to={item.path}
                                className={({ isActive: active }) =>
                                  `${sidebarLinkSubgroup} ${active ? sidebarLinkActive : ''} ${active && collapsed ? sidebarLinkActiveCollapsed : ''}`
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
          className="shrink-0 py-2 px-4 m-2 text-sm text-text-muted bg-transparent border border-border rounded-none cursor-pointer text-left hover:text-text hover:bg-sidebar-active"
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
