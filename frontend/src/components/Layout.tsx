import { useState } from 'react';
import { Outlet } from 'react-router-dom';
import { Navbar } from './Navbar';
import { Sidebar } from './Sidebar';
import { Breadcrumbs } from './Breadcrumbs';

export function Layout() {
  const [sidebarCollapsed, setSidebarCollapsed] = useState(false);

  return (
    <div className="flex flex-col min-h-screen">
      <Navbar />
      <div className="flex flex-1 min-h-0">
        {sidebarCollapsed && (
          <div className="w-10 min-w-10 shrink-0 flex items-center justify-center bg-sidebar-bg border-r border-border">
            <button
              type="button"
              className="p-2 text-sm text-text-muted bg-transparent border border-border rounded-none cursor-pointer hover:text-text hover:bg-sidebar-active"
              onClick={() => setSidebarCollapsed(!sidebarCollapsed)}
              aria-label="Expand sidebar"
              title="Expand sidebar"
            >
              <span className="inline-block" aria-hidden>
                »
              </span>
            </button>
          </div>
        )}
        <Sidebar
          collapsed={sidebarCollapsed}
          onCollapsedChange={setSidebarCollapsed}
        />
        <main className="flex-1 min-w-0 overflow-y-auto bg-bg">
          <div className="max-w-3xl mx-auto px-6 py-6 pb-8 text-left">
            <Breadcrumbs />
            <Outlet />
          </div>
        </main>
      </div>
    </div>
  );
}
