import { Link, useLocation } from 'react-router-dom';
import { getBreadcrumbLabels } from '../config/sidebarConfig';

export function Breadcrumbs() {
  const location = useLocation();
  const pathname = location.pathname;
  const labels = getBreadcrumbLabels(pathname);

  if (!labels) return null;

  const { categoryLabel, itemLabel } = labels;
  const isHome = pathname === '/' || pathname === '';

  if (isHome) {
    return (
      <nav className="mb-3" aria-label="Breadcrumb">
        <ol className="flex flex-wrap items-center gap-1 m-0 p-0 list-none text-sm text-text-muted">
          <li className="inline-flex items-center text-text font-medium" aria-current="page">
            Home
          </li>
        </ol>
      </nav>
    );
  }

  return (
    <nav className="mb-3" aria-label="Breadcrumb">
      <ol className="flex flex-wrap items-center gap-1 m-0 p-0 list-none text-sm text-text-muted">
        <li className="inline-flex items-center">
          <Link to="/" className="text-text-muted hover:text-accent">
            Home
          </Link>
        </li>
        <li className="inline-flex items-center text-text-muted opacity-70" aria-hidden>
          /
        </li>
        <li className="inline-flex items-center">
          {itemLabel && labels.categoryPath ? (
            <Link to={labels.categoryPath} className="text-text-muted hover:text-accent">
              {categoryLabel}
            </Link>
          ) : (
            <span>{categoryLabel}</span>
          )}
        </li>
        {itemLabel && (
          <>
            <li className="inline-flex items-center text-text-muted opacity-70" aria-hidden>
              /
            </li>
            <li className="inline-flex items-center text-text font-medium" aria-current="page">
              {itemLabel}
            </li>
          </>
        )}
      </ol>
    </nav>
  );
}
