import { useState } from 'react';
import { Outlet } from 'react-router-dom';
import { Navbar } from './Navbar';
import { Sidebar } from './Sidebar';
import { Breadcrumbs } from './Breadcrumbs';
import './Layout.css';

export function Layout() {
  const [sidebarCollapsed, setSidebarCollapsed] = useState(false);

  return (
    <div className="layout">
      <Navbar />
      <div className="layout-body">
        <div className="sidebar-toggle-strip">
          {sidebarCollapsed && (
            <button
              type="button"
              className="sidebar-toggle-strip-btn"
              onClick={() => setSidebarCollapsed(!sidebarCollapsed)}
              aria-label="Expand sidebar"
              title="Expand sidebar"
            >
              <span className="sidebar-toggle-strip-icon" aria-hidden>
                »
              </span>
            </button>
          )}
        </div>
        <Sidebar
          collapsed={sidebarCollapsed}
          onCollapsedChange={setSidebarCollapsed}
        />
        <main className="layout-main">
          <div className="layout-content">
            <Breadcrumbs />
            <Outlet />
          </div>
        </main>
      </div>
    </div>
  );
}
