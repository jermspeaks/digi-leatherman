export interface SidebarItem {
  id: string;
  label: string;
  path: string;
  /** Optional; items with same value are rendered as one collapsible group. Single-item groups show as a direct link. */
  subGroup?: string;
}

export interface SidebarCategory {
  id: string;
  label: string;
  items: SidebarItem[];
}

export const sidebarConfig: SidebarCategory[] = [
  {
    id: 'string',
    label: 'Strings',
    items: [
      { id: 'url-encode', label: 'URL encode', path: '/tools/string/url-encode', subGroup: 'URL' },
      { id: 'url-decode', label: 'URL decode', path: '/tools/string/url-decode', subGroup: 'URL' },
      { id: 'base64-encode', label: 'Base64 encode', path: '/tools/string/base64-encode', subGroup: 'Base64' },
      { id: 'base64-decode', label: 'Base64 decode', path: '/tools/string/base64-decode', subGroup: 'Base64' },
      { id: 'trim', label: 'Trim', path: '/tools/string/trim', subGroup: 'Trim' },
      { id: 'upper-case', label: 'Upper Case', path: '/tools/string/upper-case', subGroup: 'Case' },
      { id: 'lower-case', label: 'Lower Case', path: '/tools/string/lower-case', subGroup: 'Case' },
      { id: 'capital-case', label: 'Capital Case', path: '/tools/string/capital-case', subGroup: 'Case' },
      { id: 'snake-case', label: 'Snake Case', path: '/tools/string/snake-case', subGroup: 'Case' },
      { id: 'kebab-case', label: 'Kebab Case', path: '/tools/string/kebab-case', subGroup: 'Case' },
      { id: 'camel-case', label: 'Camel Case', path: '/tools/string/camel-case', subGroup: 'Case' },
      { id: 'pascal-case', label: 'Pascal Case', path: '/tools/string/pascal-case', subGroup: 'Case' },
      { id: 'sentence-case', label: 'Sentence Case', path: '/tools/string/sentence-case', subGroup: 'Case' },
    ],
  },
  {
    id: 'other',
    label: 'Other',
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
  if (path === '/settings') return { categoryLabel: 'Settings' };
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
