import { useState } from 'react';
import { urlEncode, urlDecode } from '../api/stringTools';
import './StringTools.css';

export type StringToolId = 'url-encode' | 'url-decode';

type StringToolsProps = {
  /** When set, show only this tool (no tabs). Used for per-route views. */
  tool?: StringToolId;
};

export function StringTools({ tool: toolProp }: StringToolsProps) {
  const [tab, setTab] = useState<StringToolId>(toolProp ?? 'url-encode');
  const [input, setInput] = useState('');
  const [output, setOutput] = useState('');
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);

  const tool = toolProp ?? tab;
  const showTabs = toolProp == null;

  const run = async () => {
    setError(null);
    setOutput('');
    if (!input.trim()) return;
    setLoading(true);
    try {
      if (tool === 'url-encode') {
        const { result } = await urlEncode(input);
        setOutput(result);
      } else {
        const { result } = await urlDecode(input);
        setOutput(result);
      }
    } catch (e) {
      setError(e instanceof Error ? e.message : 'Request failed');
    } finally {
      setLoading(false);
    }
  };

  return (
    <section className="string-tools">
      <h2>String tools</h2>
      {showTabs && (
        <nav className="tool-tabs" aria-label="String tools">
          <button
            type="button"
            className={tab === 'url-encode' ? 'active' : ''}
            onClick={() => setTab('url-encode')}
          >
            URL encode
          </button>
          <button
            type="button"
            className={tab === 'url-decode' ? 'active' : ''}
            onClick={() => setTab('url-decode')}
          >
            URL decode
          </button>
        </nav>
      )}
      <div className="tool-card">
        <label htmlFor="string-input">Input</label>
        <textarea
          id="string-input"
          value={input}
          onChange={(e) => setInput(e.target.value)}
          placeholder={
            tool === 'url-encode'
              ? 'Text to encode…'
              : 'Encoded string to decode…'
          }
          rows={4}
        />
        <button type="button" onClick={run} disabled={loading}>
          {loading ? '…' : tool === 'url-encode' ? 'Encode' : 'Decode'}
        </button>
        {error && (
          <p className="tool-error" role="alert">
            {error}
          </p>
        )}
        {output && (
          <div className="tool-output">
            <label>Output</label>
            <pre>{output}</pre>
          </div>
        )}
      </div>
    </section>
  );
}
