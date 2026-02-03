export interface SidebarItem {
  id: string;
  label: string;
  path: string;
}

export interface SidebarCategory {
  id: string;
  label: string;
  items: SidebarItem[];
}

export const sidebarConfig: SidebarCategory[] = [
  {
    id: 'string',
    label: 'String tools',
    items: [
      { id: 'url-encode', label: 'URL encode', path: '/tools/string/url-encode' },
      { id: 'url-decode', label: 'URL decode', path: '/tools/string/url-decode' },
    ],
  },
  {
    id: 'other',
    label: 'Other tools',
    items: [],
  },
];

export interface BreadcrumbLabels {
  categoryLabel: string;
  itemLabel?: string;
  /** Path to link for category (e.g. first tool in category) */
  categoryPath?: string;
}

/** Resolve path to breadcrumb labels: path -> labels */
export function getBreadcrumbLabels(path: string): BreadcrumbLabels | null {
  for (const cat of sidebarConfig) {
    for (const item of cat.items) {
      if (item.path === path) {
        const categoryPath = cat.items[0]?.path;
        return {
          categoryLabel: cat.label,
          itemLabel: item.label,
          categoryPath,
        };
      }
    }
  }
  if (path === '/' || path === '') return { categoryLabel: 'Home' };
  return null;
}
