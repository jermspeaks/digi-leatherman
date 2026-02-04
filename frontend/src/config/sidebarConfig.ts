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
      { id: 'url-parse-params', label: 'URL - Parser for params', path: '/tools/string/url-parse-params', subGroup: 'URL' },
      { id: 'url-param-creator', label: 'URL - param creator', path: '/tools/string/url-param-creator', subGroup: 'URL' },
      { id: 'base64-encode', label: 'Base64 encode', path: '/tools/string/base64-encode', subGroup: 'Base64' },
      { id: 'base64-decode', label: 'Base64 decode', path: '/tools/string/base64-decode', subGroup: 'Base64' },
      { id: 'upper-case', label: 'Upper Case', path: '/tools/string/upper-case', subGroup: 'Case' },
      { id: 'lower-case', label: 'Lower Case', path: '/tools/string/lower-case', subGroup: 'Case' },
      { id: 'capital-case', label: 'Capital Case', path: '/tools/string/capital-case', subGroup: 'Case' },
      { id: 'snake-case', label: 'Snake Case', path: '/tools/string/snake-case', subGroup: 'Case' },
      { id: 'kebab-case', label: 'Kebab Case', path: '/tools/string/kebab-case', subGroup: 'Case' },
      { id: 'camel-case', label: 'Camel Case', path: '/tools/string/camel-case', subGroup: 'Case' },
      { id: 'pascal-case', label: 'Pascal Case', path: '/tools/string/pascal-case', subGroup: 'Case' },
      { id: 'sentence-case', label: 'Sentence Case', path: '/tools/string/sentence-case', subGroup: 'Case' },
      { id: 'trim', label: 'Trim', path: '/tools/string/trim', subGroup: 'Trim' },
    ],
  },
  {
    id: 'lorem-ipsum',
    label: 'Lorem Ipsum',
    items: [
      { id: 'generator', label: 'Generator', path: '/tools/lorem-ipsum/generator', subGroup: 'generator' },
      { id: 'characters', label: 'Characters', path: '/tools/lorem-ipsum/characters', subGroup: 'characters' },
      { id: 'bytes', label: 'Bytes', path: '/tools/lorem-ipsum/bytes', subGroup: 'bytes' },
      { id: 'title', label: 'Title', path: '/tools/lorem-ipsum/title', subGroup: 'title' },
      { id: 'slug', label: 'Slug', path: '/tools/lorem-ipsum/slug', subGroup: 'slug' },
      { id: 'camel-case', label: 'Camel Case', path: '/tools/lorem-ipsum/camel-case', subGroup: 'camel-case' },
      { id: 'list', label: 'Lists', path: '/tools/lorem-ipsum/list', subGroup: 'list' },
      { id: 'headings', label: 'Headings', path: '/tools/lorem-ipsum/headings', subGroup: 'headings' },
      { id: 'html', label: 'HTML', path: '/tools/lorem-ipsum/html', subGroup: 'html' },
      { id: 'markdown', label: 'Markdown', path: '/tools/lorem-ipsum/markdown', subGroup: 'markdown' },
      { id: 'json', label: 'JSON', path: '/tools/lorem-ipsum/json', subGroup: 'json' },
    ],
  },
  {
    id: 'json',
    label: 'JSON',
    items: [
      { id: 'format', label: 'Format', path: '/tools/json/format', subGroup: 'Format' },
      { id: 'minify', label: 'Minify', path: '/tools/json/minify', subGroup: 'Format' },
      { id: 'validate', label: 'Validate', path: '/tools/json/validate', subGroup: 'Validate' },
      { id: 'path', label: 'Path query', path: '/tools/json/path', subGroup: 'Query' },
      { id: 'diff', label: 'Diff', path: '/tools/json/diff', subGroup: 'Compare' },
    ],
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
