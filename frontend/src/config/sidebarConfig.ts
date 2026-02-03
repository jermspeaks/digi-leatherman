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
    label: 'Strings',
    items: [
      { id: 'url-encode', label: 'URL encode', path: '/tools/string/url-encode' },
      { id: 'url-decode', label: 'URL decode', path: '/tools/string/url-decode' },
      { id: 'base64-encode', label: 'Base64 encode', path: '/tools/string/base64-encode' },
      { id: 'base64-decode', label: 'Base64 decode', path: '/tools/string/base64-decode' },
      { id: 'trim', label: 'Trim', path: '/tools/string/trim' },
      { id: 'upper-case', label: 'Upper Case', path: '/tools/string/upper-case' },
      { id: 'lower-case', label: 'Lower Case', path: '/tools/string/lower-case' },
      { id: 'capital-case', label: 'Capital Case', path: '/tools/string/capital-case' },
      { id: 'snake-case', label: 'Snake Case', path: '/tools/string/snake-case' },
      { id: 'kebab-case', label: 'Kebab Case', path: '/tools/string/kebab-case' },
      { id: 'camel-case', label: 'Camel Case', path: '/tools/string/camel-case' },
      { id: 'pascal-case', label: 'Pascal Case', path: '/tools/string/pascal-case' },
      { id: 'sentence-case', label: 'Sentence Case', path: '/tools/string/sentence-case' },
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
