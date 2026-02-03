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
              value={`${tool.label} ${tool.description}`}
              onSelect={() => handleSelect(tool.path)}
            >
              <span>{tool.label}</span>
              {tool.description && (
                <span className="ml-2 text-muted-foreground truncate max-w-[200px]">
                  {tool.description}
                </span>
              )}
            </CommandItem>
          ))}
        </CommandGroup>
      </CommandList>
    </CommandDialog>
  );
}
