import { useMemo } from 'react';
import { useNavigate } from 'react-router-dom';
import {
  CommandDialog,
  CommandEmpty,
  CommandGroup,
  CommandInput,
  CommandItem,
  CommandList,
} from '@/components/ui/command';
import { getToolsForSearch } from '@/lib/toolsForSearch';

type CommandPaletteProps = {
  open: boolean;
  onOpenChange: (open: boolean) => void;
};

export function CommandPalette({ open, onOpenChange }: CommandPaletteProps) {
  const navigate = useNavigate();
  const tools = useMemo(() => getToolsForSearch(), []);

  const handleSelect = (path: string) => {
    navigate(path);
    onOpenChange(false);
  };

  return (
    <CommandDialog
      open={open}
      onOpenChange={onOpenChange}
      title="Search tools"
      description="Search by tool name or description."
    >
      <CommandInput placeholder="Search tools…" />
      <CommandList>
        <CommandEmpty>No tools found.</CommandEmpty>
        <CommandGroup heading="Tools">
          {tools.map((tool) => (
            <CommandItem
              key={tool.id}
              value={`${tool.label} ${tool.categoryLabel} ${tool.description}`}
              onSelect={() => handleSelect(tool.path)}
            >
              <div className="flex flex-col gap-0.5 min-w-0">
                <div className="flex items-center gap-2">
                  <span>{tool.label}</span>
                  <span className="text-muted-foreground text-xs shrink-0">
                    {tool.categoryLabel}
                  </span>
                </div>
                {tool.description && (
                  <span className="text-muted-foreground text-xs truncate">
                    {tool.description}
                  </span>
                )}
              </div>
            </CommandItem>
          ))}
        </CommandGroup>
      </CommandList>
    </CommandDialog>
  );
}
