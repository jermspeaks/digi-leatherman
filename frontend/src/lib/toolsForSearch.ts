import { sidebarConfig } from '@/config/sidebarConfig';
import { LOREM_TOOL_DESCRIPTIONS } from '@/components/LoremTools';
import { TOOL_DESCRIPTIONS } from '@/components/StringTools';

export interface ToolForSearch {
  id: string;
  label: string;
  description: string;
  path: string;
  categoryLabel: string;
}

/**
 * Flattens sidebarConfig into a list of tools with descriptions for the command palette.
 * Search matches both label and description via a combined string for cmdk.
 */
export function getToolsForSearch(): ToolForSearch[] {
  const result: ToolForSearch[] = [];
  for (const category of sidebarConfig) {
    for (const item of category.items) {
      result.push({
        id: item.id,
        label: item.label,
        description: TOOL_DESCRIPTIONS[item.id] ?? LOREM_TOOL_DESCRIPTIONS[item.id] ?? '',
        path: item.path,
        categoryLabel: category.label,
      });
    }
  }
  return result;
}
