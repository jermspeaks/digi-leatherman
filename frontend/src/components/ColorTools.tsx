import { useRef, useState, type ChangeEvent } from 'react';
import { Copy, Upload, Check, X } from 'lucide-react';
import { Button } from '@/components/ui/button';
import {
  convertColor,
  checkContrast,
  simulateBlindness,
  extractPalette,
  generateShades,
  generateHarmonies,
  type ColorConvertResult,
  type ContrastResult,
  type BlindnessResult,
  type PaletteResult,
  type ShadesResult,
  type HarmoniesResult,
} from '../api/colorTools';

export type ColorToolId =
  | 'convert'
  | 'contrast'
  | 'blindness'
  | 'palette'
  | 'shades'
  | 'harmonies'
  | 'picker';

type ToolMeta = {
  id: ColorToolId;
  label: string;
  description: string;
  example: string;
  placeholder?: string;
  buttonLabel: string;
};

const TOOL_CONFIG: ToolMeta[] = [
  {
    id: 'convert',
    label: 'Convert',
    description: 'Convert any color format (hex, RGB, HSL, oklch, named colors) to all other formats.',
    example: '#ff6b6b, rgb(255, 107, 107), hsl(0, 100%, 71%), coral',
    placeholder: 'Enter a color in any format...',
    buttonLabel: 'Convert',
  },
  {
    id: 'picker',
    label: 'Picker',
    description: 'Pick a color visually and see it in all formats.',
    example: 'Click the color picker to select a color',
    buttonLabel: 'Pick',
  },
  {
    id: 'contrast',
    label: 'Contrast',
    description: 'Check WCAG contrast ratio between two colors for accessibility compliance.',
    example: '#000000 and #ffffff → 21:1 (AAA)',
    placeholder: 'Enter first color...',
    buttonLabel: 'Check Contrast',
  },
  {
    id: 'blindness',
    label: 'Color Blindness',
    description: 'Simulate how colors appear with different types of color vision deficiency.',
    example: '#ff0000 → protanopia simulation',
    placeholder: 'Enter a color to simulate...',
    buttonLabel: 'Simulate',
  },
  {
    id: 'palette',
    label: 'Palette',
    description: 'Extract dominant colors from an uploaded image.',
    example: 'Upload an image to extract its color palette',
    buttonLabel: 'Extract',
  },
  {
    id: 'shades',
    label: 'Shades & Tints',
    description: 'Generate darker (shades) and lighter (tints) variations of a color.',
    example: '#3b82f6 → 5 shades and 5 tints',
    placeholder: 'Enter a base color...',
    buttonLabel: 'Generate',
  },
  {
    id: 'harmonies',
    label: 'Harmonies',
    description: 'Generate color harmonies: complementary, triadic, analogous, and more.',
    example: '#3b82f6 → complementary, triadic, analogous colors',
    placeholder: 'Enter a base color...',
    buttonLabel: 'Generate',
  },
];

export const TOOL_DESCRIPTIONS: Record<string, string> = Object.fromEntries(
  TOOL_CONFIG.map((c) => [c.id, c.description])
);

const TOOL_MAP = Object.fromEntries(TOOL_CONFIG.map((c) => [c.id, c])) as Record<ColorToolId, ToolMeta>;

const BLINDNESS_TYPES = [
  { id: 'protanopia', label: 'Protanopia (Red-blind)' },
  { id: 'deuteranopia', label: 'Deuteranopia (Green-blind)' },
  { id: 'tritanopia', label: 'Tritanopia (Blue-blind)' },
  { id: 'achromatopsia', label: 'Achromatopsia (Total)' },
];

type ColorToolsProps = {
  tool?: ColorToolId;
};

export function ColorTools({ tool: toolProp }: ColorToolsProps) {
  const [tab, setTab] = useState<ColorToolId>(toolProp ?? 'convert');
  const [input, setInput] = useState('');
  const [input2, setInput2] = useState('');
  const [pickerColor, setPickerColor] = useState('#3b82f6');
  const [blindnessType, setBlindnessType] = useState('deuteranopia');
  const [shadesCount, setShadesCount] = useState(5);
  const [paletteCount, setPaletteCount] = useState(5);
  const [imagePreview, setImagePreview] = useState<string | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);
  const [copied, setCopied] = useState<string | null>(null);
  const copyTimeoutRef = useRef<ReturnType<typeof setTimeout> | null>(null);
  const fileInputRef = useRef<HTMLInputElement>(null);

  const [convertResult, setConvertResult] = useState<ColorConvertResult | null>(null);
  const [contrastResult, setContrastResult] = useState<ContrastResult | null>(null);
  const [blindnessResult, setBlindnessResult] = useState<BlindnessResult | null>(null);
  const [paletteResult, setPaletteResult] = useState<PaletteResult | null>(null);
  const [shadesResult, setShadesResult] = useState<ShadesResult | null>(null);
  const [harmoniesResult, setHarmoniesResult] = useState<HarmoniesResult | null>(null);

  const tool = toolProp ?? tab;
  const showTabs = toolProp == null;
  const config = TOOL_MAP[tool];

  const clearResults = () => {
    setConvertResult(null);
    setContrastResult(null);
    setBlindnessResult(null);
    setPaletteResult(null);
    setShadesResult(null);
    setHarmoniesResult(null);
  };

  const run = async () => {
    setError(null);
    clearResults();
    if (copyTimeoutRef.current) {
      clearTimeout(copyTimeoutRef.current);
      copyTimeoutRef.current = null;
    }
    setCopied(null);

    setLoading(true);
    try {
      switch (tool) {
        case 'convert': {
          if (!input.trim()) return;
          const result = await convertColor(input);
          setConvertResult(result);
          break;
        }
        case 'picker': {
          const result = await convertColor(pickerColor);
          setConvertResult(result);
          break;
        }
        case 'contrast': {
          if (!input.trim() || !input2.trim()) {
            setError('Both colors are required');
            return;
          }
          const result = await checkContrast(input, input2);
          setContrastResult(result);
          break;
        }
        case 'blindness': {
          if (!input.trim()) return;
          const result = await simulateBlindness(input, blindnessType);
          setBlindnessResult(result);
          break;
        }
        case 'palette': {
          if (!imagePreview) {
            setError('Please upload an image');
            return;
          }
          const result = await extractPalette(imagePreview, paletteCount);
          setPaletteResult(result);
          break;
        }
        case 'shades': {
          if (!input.trim()) return;
          const result = await generateShades(input, shadesCount);
          setShadesResult(result);
          break;
        }
        case 'harmonies': {
          if (!input.trim()) return;
          const result = await generateHarmonies(input);
          setHarmoniesResult(result);
          break;
        }
      }
    } catch (e) {
      setError(e instanceof Error ? e.message : 'Request failed');
    } finally {
      setLoading(false);
    }
  };

  const useExample = () => {
    switch (tool) {
      case 'convert':
        setInput('#ff6b6b');
        break;
      case 'picker':
        setPickerColor('#3b82f6');
        break;
      case 'contrast':
        setInput('#000000');
        setInput2('#ffffff');
        break;
      case 'blindness':
        setInput('#ff0000');
        setBlindnessType('protanopia');
        break;
      case 'shades':
        setInput('#3b82f6');
        setShadesCount(5);
        break;
      case 'harmonies':
        setInput('#3b82f6');
        break;
    }
  };

  const handleCopy = async (value: string, label: string) => {
    if (!value) return;
    if (copyTimeoutRef.current) clearTimeout(copyTimeoutRef.current);
    try {
      await navigator.clipboard.writeText(value);
      setCopied(label);
      copyTimeoutRef.current = setTimeout(() => {
        copyTimeoutRef.current = null;
        setCopied(null);
      }, 1500);
    } catch {
      // Permission denied
    }
  };

  const handleImageUpload = (e: ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (!file) return;
    if (!file.type.startsWith('image/')) {
      setError('Please upload an image file');
      return;
    }
    const reader = new FileReader();
    reader.onload = (event) => {
      setImagePreview(event.target?.result as string);
      setError(null);
    };
    reader.readAsDataURL(file);
  };

  const ColorSwatch = ({ color, label, showCopy = true }: { color: string; label?: string; showCopy?: boolean }) => (
    <div className="flex items-center gap-2">
      <div
        className="size-8 rounded border border-border shrink-0"
        style={{ backgroundColor: color }}
        title={color}
      />
      {label && <span className="text-sm text-text-secondary">{label}</span>}
      <code className="text-sm text-text">{color}</code>
      {showCopy && (
        <Button
          type="button"
          variant="ghost"
          size="sm"
          onClick={() => handleCopy(color, color)}
          className="shrink-0 size-7 p-0"
          aria-label={copied === color ? 'Copied' : `Copy ${color}`}
        >
          {copied === color ? <Check className="size-3.5" /> : <Copy className="size-3.5" />}
        </Button>
      )}
    </div>
  );

  const OutputRow = ({ label, value }: { label: string; value: string }) => (
    <div className="flex items-center justify-between gap-2 py-2 border-b border-border last:border-b-0">
      <div className="flex flex-col gap-0.5 min-w-0 flex-1">
        <span className="text-sm text-text-secondary">{label}</span>
        <code className="text-text break-all">{value}</code>
      </div>
      <Button
        type="button"
        variant="outline"
        size="sm"
        onClick={() => handleCopy(value, label)}
        className="shrink-0"
        aria-label={copied === label ? 'Copied' : `Copy ${label}`}
      >
        <Copy className="size-3.5" aria-hidden />
        {copied === label ? 'Copied!' : 'Copy'}
      </Button>
    </div>
  );

  const WCAGBadge = ({ pass, label }: { pass: boolean; label: string }) => (
    <div
      className={`flex items-center gap-1.5 px-2 py-1 rounded text-sm font-medium ${
        pass ? 'bg-green-500/20 text-green-600 dark:text-green-400' : 'bg-red-500/20 text-red-600 dark:text-red-400'
      }`}
    >
      {pass ? <Check className="size-3.5" /> : <X className="size-3.5" />}
      {label}
    </div>
  );

  return (
    <section className="text-left mt-6">
      <h2 className="mb-3 text-2xl text-text">Colors</h2>
      {showTabs && (
        <nav className="flex flex-wrap gap-2 mb-4" aria-label="Color tools">
          {TOOL_CONFIG.map((c) => (
            <button
              key={c.id}
              type="button"
              className={`py-2 px-4 ${tab === c.id ? 'border-accent bg-sidebar-active' : ''}`}
              onClick={() => {
                setTab(c.id);
                clearResults();
                setError(null);
              }}
            >
              {c.label}
            </button>
          ))}
        </nav>
      )}
      <p className="mb-2 text-text-secondary max-w-2xl">{config.description}</p>
      <p className="mb-4 text-sm text-text-secondary max-w-2xl">
        Example: <code className="px-1 rounded bg-bg-elevated">{config.example}</code>
        {tool !== 'palette' && (
          <button type="button" onClick={useExample} className="ml-2 text-accent hover:underline">
            Use example
          </button>
        )}
      </p>
      <div className="flex flex-col gap-3 max-w-2xl">
        {/* Convert tool */}
        {tool === 'convert' && (
          <>
            <label htmlFor="color-input" className="font-medium">
              Color
            </label>
            <input
              id="color-input"
              type="text"
              className="w-full p-2.5 font-inherit text-base rounded-lg border border-border bg-bg text-text"
              value={input}
              onChange={(e) => setInput(e.target.value)}
              placeholder={config.placeholder}
            />
          </>
        )}

        {/* Picker tool */}
        {tool === 'picker' && (
          <>
            <label htmlFor="color-picker" className="font-medium">
              Pick a Color
            </label>
            <div className="flex items-center gap-4">
              <input
                id="color-picker"
                type="color"
                value={pickerColor}
                onChange={(e) => setPickerColor(e.target.value)}
                className="size-16 p-0 border-2 border-border rounded-lg cursor-pointer"
              />
              <div className="flex flex-col gap-1">
                <code className="text-lg font-medium text-text">{pickerColor}</code>
                <span className="text-sm text-text-secondary">Click to pick a color</span>
              </div>
            </div>
          </>
        )}

        {/* Contrast tool */}
        {tool === 'contrast' && (
          <>
            <div className="flex flex-wrap gap-4">
              <div className="flex flex-col gap-1.5 flex-1 min-w-[140px]">
                <label htmlFor="color1-input" className="font-medium">
                  Foreground Color
                </label>
                <input
                  id="color1-input"
                  type="text"
                  className="w-full p-2.5 font-inherit text-base rounded-lg border border-border bg-bg text-text"
                  value={input}
                  onChange={(e) => setInput(e.target.value)}
                  placeholder="#000000"
                />
              </div>
              <div className="flex flex-col gap-1.5 flex-1 min-w-[140px]">
                <label htmlFor="color2-input" className="font-medium">
                  Background Color
                </label>
                <input
                  id="color2-input"
                  type="text"
                  className="w-full p-2.5 font-inherit text-base rounded-lg border border-border bg-bg text-text"
                  value={input2}
                  onChange={(e) => setInput2(e.target.value)}
                  placeholder="#ffffff"
                />
              </div>
            </div>
          </>
        )}

        {/* Blindness tool */}
        {tool === 'blindness' && (
          <>
            <label htmlFor="blindness-input" className="font-medium">
              Color
            </label>
            <input
              id="blindness-input"
              type="text"
              className="w-full p-2.5 font-inherit text-base rounded-lg border border-border bg-bg text-text"
              value={input}
              onChange={(e) => setInput(e.target.value)}
              placeholder={config.placeholder}
            />
            <label htmlFor="blindness-type" className="font-medium">
              Simulation Type
            </label>
            <select
              id="blindness-type"
              value={blindnessType}
              onChange={(e) => setBlindnessType(e.target.value)}
              className="p-2.5 rounded-lg border border-border bg-bg text-text"
            >
              {BLINDNESS_TYPES.map((t) => (
                <option key={t.id} value={t.id}>
                  {t.label}
                </option>
              ))}
            </select>
          </>
        )}

        {/* Palette tool */}
        {tool === 'palette' && (
          <>
            <label className="font-medium">Upload Image</label>
            <div
              className="flex flex-col items-center justify-center gap-3 p-6 border-2 border-dashed border-border rounded-lg bg-bg-elevated cursor-pointer hover:border-accent transition-colors"
              onClick={() => fileInputRef.current?.click()}
            >
              {imagePreview ? (
                <img
                  src={imagePreview}
                  alt="Uploaded preview"
                  className="max-h-48 rounded-lg object-contain"
                />
              ) : (
                <>
                  <Upload className="size-8 text-text-secondary" />
                  <span className="text-text-secondary">Click to upload or drag and drop</span>
                  <span className="text-sm text-text-secondary">PNG, JPG, GIF up to 10MB</span>
                </>
              )}
            </div>
            <input
              ref={fileInputRef}
              type="file"
              accept="image/*"
              onChange={handleImageUpload}
              className="hidden"
            />
            <div className="flex items-center gap-3">
              <label htmlFor="palette-count" className="text-sm font-medium">
                Number of colors:
              </label>
              <input
                id="palette-count"
                type="number"
                min={2}
                max={10}
                value={paletteCount}
                onChange={(e) => setPaletteCount(Math.min(10, Math.max(2, parseInt(e.target.value) || 5)))}
                className="w-16 p-2 rounded-lg border border-border bg-bg text-text text-center"
              />
            </div>
          </>
        )}

        {/* Shades tool */}
        {tool === 'shades' && (
          <>
            <label htmlFor="shades-input" className="font-medium">
              Base Color
            </label>
            <input
              id="shades-input"
              type="text"
              className="w-full p-2.5 font-inherit text-base rounded-lg border border-border bg-bg text-text"
              value={input}
              onChange={(e) => setInput(e.target.value)}
              placeholder={config.placeholder}
            />
            <div className="flex items-center gap-3">
              <label htmlFor="shades-count" className="text-sm font-medium">
                Number of shades/tints:
              </label>
              <input
                id="shades-count"
                type="number"
                min={1}
                max={10}
                value={shadesCount}
                onChange={(e) => setShadesCount(Math.min(10, Math.max(1, parseInt(e.target.value) || 5)))}
                className="w-16 p-2 rounded-lg border border-border bg-bg text-text text-center"
              />
            </div>
          </>
        )}

        {/* Harmonies tool */}
        {tool === 'harmonies' && (
          <>
            <label htmlFor="harmonies-input" className="font-medium">
              Base Color
            </label>
            <input
              id="harmonies-input"
              type="text"
              className="w-full p-2.5 font-inherit text-base rounded-lg border border-border bg-bg text-text"
              value={input}
              onChange={(e) => setInput(e.target.value)}
              placeholder={config.placeholder}
            />
          </>
        )}

        <button type="button" onClick={run} disabled={loading}>
          {loading ? '…' : config.buttonLabel}
        </button>

        {error && (
          <p className="text-red-500 m-0" role="alert">
            {error}
          </p>
        )}

        {/* Convert/Picker results */}
        {convertResult && (
          <div className="mt-2">
            <label className="font-medium mb-2 block">Output</label>
            <div className="bg-bg-elevated rounded-lg p-3">
              <div className="flex items-center gap-3 mb-3 pb-3 border-b border-border">
                <div
                  className="size-16 rounded-lg border border-border"
                  style={{ backgroundColor: convertResult.hex }}
                />
                <div className="flex flex-col">
                  <span className="text-lg font-medium text-text">{convertResult.hex}</span>
                  <span className="text-sm text-text-secondary">Preview</span>
                </div>
              </div>
              <OutputRow label="Hex" value={convertResult.hex} />
              <OutputRow label="RGB" value={convertResult.rgb} />
              <OutputRow label="HSL" value={convertResult.hsl} />
              <OutputRow label="OKLCH" value={convertResult.oklch} />
            </div>
          </div>
        )}

        {/* Contrast results */}
        {contrastResult && (
          <div className="mt-2">
            <label className="font-medium mb-2 block">Output</label>
            <div className="bg-bg-elevated rounded-lg p-3">
              <div className="flex items-center gap-4 mb-4 pb-4 border-b border-border">
                <div className="flex-1 flex items-center justify-center p-4 rounded-lg"
                  style={{ backgroundColor: contrastResult.color2Hex }}
                >
                  <span
                    className="text-lg font-medium"
                    style={{ color: contrastResult.color1Hex }}
                  >
                    Sample Text
                  </span>
                </div>
                <div className="flex flex-col items-center">
                  <span className="text-3xl font-bold text-text">{contrastResult.ratioText}</span>
                  <span className="text-sm text-text-secondary">Contrast Ratio</span>
                </div>
              </div>
              <div className="flex items-center gap-4 mb-3">
                <ColorSwatch color={contrastResult.color1Hex} label="Foreground" showCopy={false} />
                <ColorSwatch color={contrastResult.color2Hex} label="Background" showCopy={false} />
              </div>
              <div className="grid grid-cols-2 gap-2 mt-4">
                <WCAGBadge pass={contrastResult.aaNormal} label="AA Normal" />
                <WCAGBadge pass={contrastResult.aaLarge} label="AA Large" />
                <WCAGBadge pass={contrastResult.aaaNormal} label="AAA Normal" />
                <WCAGBadge pass={contrastResult.aaaLarge} label="AAA Large" />
              </div>
              <p className="text-sm text-text-secondary mt-3">
                Normal text: 4.5:1 for AA, 7:1 for AAA. Large text (18pt+ or 14pt+ bold): 3:1 for AA, 4.5:1 for AAA.
              </p>
            </div>
          </div>
        )}

        {/* Blindness results */}
        {blindnessResult && (
          <div className="mt-2">
            <label className="font-medium mb-2 block">Output</label>
            <div className="bg-bg-elevated rounded-lg p-3">
              <div className="flex items-center gap-4 mb-3">
                <div className="flex flex-col items-center gap-1">
                  <div
                    className="size-16 rounded-lg border border-border"
                    style={{ backgroundColor: blindnessResult.original }}
                  />
                  <span className="text-sm text-text-secondary">Original</span>
                  <code className="text-sm text-text">{blindnessResult.original}</code>
                </div>
                <span className="text-2xl text-text-secondary">→</span>
                <div className="flex flex-col items-center gap-1">
                  <div
                    className="size-16 rounded-lg border border-border"
                    style={{ backgroundColor: blindnessResult.simulated }}
                  />
                  <span className="text-sm text-text-secondary">{BLINDNESS_TYPES.find(t => t.id === blindnessResult.type)?.label}</span>
                  <code className="text-sm text-text">{blindnessResult.simulated}</code>
                </div>
              </div>
            </div>
          </div>
        )}

        {/* Palette results */}
        {paletteResult && (
          <div className="mt-2">
            <label className="font-medium mb-2 block">Extracted Colors</label>
            <div className="bg-bg-elevated rounded-lg p-3">
              <div className="flex gap-2 mb-4">
                {paletteResult.colors.map((color, i) => (
                  <div
                    key={i}
                    className="flex-1 h-16 rounded border border-border cursor-pointer hover:scale-105 transition-transform"
                    style={{ backgroundColor: color }}
                    onClick={() => handleCopy(color, color)}
                    title={`Click to copy: ${color}`}
                  />
                ))}
              </div>
              <div className="flex flex-col gap-2">
                {paletteResult.colors.map((color, i) => (
                  <ColorSwatch key={i} color={color} />
                ))}
              </div>
            </div>
          </div>
        )}

        {/* Shades results */}
        {shadesResult && (
          <div className="mt-2">
            <label className="font-medium mb-2 block">Output</label>
            <div className="bg-bg-elevated rounded-lg p-3">
              <div className="mb-4">
                <span className="text-sm font-medium text-text-secondary block mb-2">Base Color</span>
                <ColorSwatch color={shadesResult.base} />
              </div>
              <div className="mb-4">
                <span className="text-sm font-medium text-text-secondary block mb-2">Shades (Darker)</span>
                <div className="flex gap-1 mb-2">
                  {shadesResult.shades.map((color, i) => (
                    <div
                      key={i}
                      className="flex-1 h-10 rounded cursor-pointer hover:scale-105 transition-transform"
                      style={{ backgroundColor: color }}
                      onClick={() => handleCopy(color, color)}
                      title={`Click to copy: ${color}`}
                    />
                  ))}
                </div>
                <div className="flex flex-wrap gap-2">
                  {shadesResult.shades.map((color, i) => (
                    <code key={i} className="text-xs px-1.5 py-0.5 bg-bg rounded">{color}</code>
                  ))}
                </div>
              </div>
              <div>
                <span className="text-sm font-medium text-text-secondary block mb-2">Tints (Lighter)</span>
                <div className="flex gap-1 mb-2">
                  {shadesResult.tints.map((color, i) => (
                    <div
                      key={i}
                      className="flex-1 h-10 rounded border border-border cursor-pointer hover:scale-105 transition-transform"
                      style={{ backgroundColor: color }}
                      onClick={() => handleCopy(color, color)}
                      title={`Click to copy: ${color}`}
                    />
                  ))}
                </div>
                <div className="flex flex-wrap gap-2">
                  {shadesResult.tints.map((color, i) => (
                    <code key={i} className="text-xs px-1.5 py-0.5 bg-bg rounded">{color}</code>
                  ))}
                </div>
              </div>
            </div>
          </div>
        )}

        {/* Harmonies results */}
        {harmoniesResult && (
          <div className="mt-2">
            <label className="font-medium mb-2 block">Output</label>
            <div className="bg-bg-elevated rounded-lg p-3 space-y-4">
              <div>
                <span className="text-sm font-medium text-text-secondary block mb-2">Base Color</span>
                <ColorSwatch color={harmoniesResult.base} />
              </div>
              <div>
                <span className="text-sm font-medium text-text-secondary block mb-2">Complementary</span>
                <div className="flex gap-2">
                  <div
                    className="size-10 rounded border border-border"
                    style={{ backgroundColor: harmoniesResult.base }}
                    title={harmoniesResult.base}
                  />
                  <div
                    className="size-10 rounded border border-border cursor-pointer hover:scale-105 transition-transform"
                    style={{ backgroundColor: harmoniesResult.complementary }}
                    onClick={() => handleCopy(harmoniesResult.complementary, 'complementary')}
                    title={harmoniesResult.complementary}
                  />
                </div>
                <code className="text-xs mt-1 block">{harmoniesResult.complementary}</code>
              </div>
              <div>
                <span className="text-sm font-medium text-text-secondary block mb-2">Triadic</span>
                <div className="flex gap-2">
                  <div className="size-10 rounded border border-border" style={{ backgroundColor: harmoniesResult.base }} title={harmoniesResult.base} />
                  {harmoniesResult.triadic.map((color, i) => (
                    <div
                      key={i}
                      className="size-10 rounded border border-border cursor-pointer hover:scale-105 transition-transform"
                      style={{ backgroundColor: color }}
                      onClick={() => handleCopy(color, `triadic-${i}`)}
                      title={color}
                    />
                  ))}
                </div>
                <div className="flex gap-2 mt-1">
                  {harmoniesResult.triadic.map((color, i) => (
                    <code key={i} className="text-xs">{color}</code>
                  ))}
                </div>
              </div>
              <div>
                <span className="text-sm font-medium text-text-secondary block mb-2">Analogous</span>
                <div className="flex gap-2">
                  {harmoniesResult.analogous.slice(0, 1).map((color, i) => (
                    <div
                      key={i}
                      className="size-10 rounded border border-border cursor-pointer hover:scale-105 transition-transform"
                      style={{ backgroundColor: color }}
                      onClick={() => handleCopy(color, `analogous-${i}`)}
                      title={color}
                    />
                  ))}
                  <div className="size-10 rounded border border-border" style={{ backgroundColor: harmoniesResult.base }} title={harmoniesResult.base} />
                  {harmoniesResult.analogous.slice(1).map((color, i) => (
                    <div
                      key={i + 1}
                      className="size-10 rounded border border-border cursor-pointer hover:scale-105 transition-transform"
                      style={{ backgroundColor: color }}
                      onClick={() => handleCopy(color, `analogous-${i + 1}`)}
                      title={color}
                    />
                  ))}
                </div>
                <div className="flex gap-2 mt-1">
                  {harmoniesResult.analogous.map((color, i) => (
                    <code key={i} className="text-xs">{color}</code>
                  ))}
                </div>
              </div>
              <div>
                <span className="text-sm font-medium text-text-secondary block mb-2">Split Complementary</span>
                <div className="flex gap-2">
                  <div className="size-10 rounded border border-border" style={{ backgroundColor: harmoniesResult.base }} title={harmoniesResult.base} />
                  {harmoniesResult.splitComplementary.map((color, i) => (
                    <div
                      key={i}
                      className="size-10 rounded border border-border cursor-pointer hover:scale-105 transition-transform"
                      style={{ backgroundColor: color }}
                      onClick={() => handleCopy(color, `split-${i}`)}
                      title={color}
                    />
                  ))}
                </div>
                <div className="flex gap-2 mt-1">
                  {harmoniesResult.splitComplementary.map((color, i) => (
                    <code key={i} className="text-xs">{color}</code>
                  ))}
                </div>
              </div>
              <div>
                <span className="text-sm font-medium text-text-secondary block mb-2">Tetradic (Square)</span>
                <div className="flex gap-2">
                  <div className="size-10 rounded border border-border" style={{ backgroundColor: harmoniesResult.base }} title={harmoniesResult.base} />
                  {harmoniesResult.tetradic.map((color, i) => (
                    <div
                      key={i}
                      className="size-10 rounded border border-border cursor-pointer hover:scale-105 transition-transform"
                      style={{ backgroundColor: color }}
                      onClick={() => handleCopy(color, `tetradic-${i}`)}
                      title={color}
                    />
                  ))}
                </div>
                <div className="flex gap-2 mt-1">
                  {harmoniesResult.tetradic.map((color, i) => (
                    <code key={i} className="text-xs">{color}</code>
                  ))}
                </div>
              </div>
            </div>
          </div>
        )}
      </div>
    </section>
  );
}
