---
name: shadcn-ui-usage
description: Guide to shadcn/ui components in Digital Leatherman. Use when working with UI components, styling, or building interfaces. Triggers on "shadcn", "ui component", "button", "dialog", "card", "tabs", "input", "textarea".
---

# shadcn/ui Components

This project uses shadcn/ui components built on Radix UI primitives with Tailwind CSS styling.

## Available Components

Located in `frontend/src/components/ui/`:

| Component | File | Description |
|-----------|------|-------------|
| Alert | `alert.tsx` | Informational alerts |
| Breadcrumb | `breadcrumb.tsx` | Navigation breadcrumbs |
| Button | `button.tsx` | Buttons with variants |
| Card | `card.tsx` | Container cards |
| Command | `command.tsx` | Command palette (cmdk) |
| Dialog | `dialog.tsx` | Modal dialogs |
| Input | `input.tsx` | Text inputs |
| Label | `label.tsx` | Form labels |
| Tabs | `tabs.tsx` | Tabbed interfaces |
| Textarea | `textarea.tsx` | Multi-line text input |
| Tooltip | `tooltip.tsx` | Hover tooltips |

## Import Pattern

```typescript
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Dialog, DialogContent, DialogHeader } from '@/components/ui/dialog';
```

## Button

Variants: `default`, `destructive`, `outline`, `secondary`, `ghost`, `link`
Sizes: `default`, `sm`, `lg`, `icon`

```tsx
<Button variant="outline" size="sm" onClick={handleClick}>
  <Copy className="size-3.5" />
  Copy
</Button>

<Button variant="destructive">Delete</Button>

<Button variant="ghost" size="icon">
  <Settings />
</Button>
```

## Input

```tsx
<Input
  type="text"
  placeholder="Enter value..."
  value={value}
  onChange={(e) => setValue(e.target.value)}
/>
```

## Textarea

```tsx
<Textarea
  value={input}
  onChange={(e) => setInput(e.target.value)}
  placeholder="Enter text..."
  rows={4}
/>
```

## Dialog

```tsx
<Dialog open={open} onOpenChange={setOpen}>
  <DialogTrigger asChild>
    <Button>Open</Button>
  </DialogTrigger>
  <DialogContent>
    <DialogHeader>
      <DialogTitle>Title</DialogTitle>
      <DialogDescription>Description text</DialogDescription>
    </DialogHeader>
    {/* Content */}
  </DialogContent>
</Dialog>
```

## Tabs

```tsx
<Tabs defaultValue="tab1">
  <TabsList>
    <TabsTrigger value="tab1">Tab 1</TabsTrigger>
    <TabsTrigger value="tab2">Tab 2</TabsTrigger>
  </TabsList>
  <TabsContent value="tab1">Content 1</TabsContent>
  <TabsContent value="tab2">Content 2</TabsContent>
</Tabs>
```

## Card

```tsx
<Card>
  <CardHeader>
    <CardTitle>Title</CardTitle>
    <CardDescription>Description</CardDescription>
  </CardHeader>
  <CardContent>
    Content here
  </CardContent>
</Card>
```

## Command (Command Palette)

```tsx
<Command>
  <CommandInput placeholder="Search..." />
  <CommandList>
    <CommandEmpty>No results</CommandEmpty>
    <CommandGroup heading="Actions">
      <CommandItem onSelect={() => action()}>
        <Icon className="mr-2 size-4" />
        Action Name
      </CommandItem>
    </CommandGroup>
  </CommandList>
</Command>
```

## cn() Utility

Merge Tailwind classes with `clsx` + `tailwind-merge`:

```typescript
import { cn } from '@/lib/utils';

<div className={cn(
  "base-classes",
  isActive && "active-classes",
  className
)} />
```

Located in `frontend/src/lib/utils.ts`:

```typescript
import { clsx, type ClassValue } from 'clsx';
import { twMerge } from 'tailwind-merge';

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}
```

## Theme CSS Variables

Defined in `frontend/src/index.css`:

- `--text` / `--text-secondary` - Text colors
- `--bg` / `--bg-elevated` - Background colors
- `--border` - Border color
- `--accent` - Accent/link color
- `--sidebar-*` - Sidebar specific colors

Use with Tailwind: `text-text`, `bg-bg-elevated`, `border-border`, etc.

## Icons

Using `lucide-react`:

```typescript
import { Copy, Settings, ChevronRight } from 'lucide-react';

<Copy className="size-4" />
<Settings className="size-5 text-text-secondary" />
```

## Adding New Components

If you need a component not in `ui/`, add it from shadcn:

```bash
cd frontend
npx shadcn@latest add [component-name]
```
