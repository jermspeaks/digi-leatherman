import { Link, useLocation } from 'react-router-dom';
import { getBreadcrumbLabels } from '../config/sidebarConfig';
import './Breadcrumbs.css';

export function Breadcrumbs() {
  const location = useLocation();
  const pathname = location.pathname;
  const labels = getBreadcrumbLabels(pathname);

  if (!labels) return null;

  const { categoryLabel, itemLabel } = labels;
  const isHome = pathname === '/' || pathname === '';

  if (isHome) {
    return (
      <nav className="breadcrumbs" aria-label="Breadcrumb">
        <ol className="breadcrumbs-list">
          <li aria-current="page">Home</li>
        </ol>
      </nav>
    );
  }

  return (
    <nav className="breadcrumbs" aria-label="Breadcrumb">
      <ol className="breadcrumbs-list">
        <li>
          <Link to="/">Home</Link>
        </li>
        <li className="breadcrumbs-sep" aria-hidden>
          /
        </li>
        <li>
          {itemLabel && labels.categoryPath ? (
            <Link to={labels.categoryPath}>{categoryLabel}</Link>
          ) : (
            <span>{categoryLabel}</span>
          )}
        </li>
        {itemLabel && (
          <>
            <li className="breadcrumbs-sep" aria-hidden>
              /
            </li>
            <li aria-current="page">{itemLabel}</li>
          </>
        )}
      </ol>
    </nav>
  );
}
