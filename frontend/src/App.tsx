import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import { ThemeProvider } from './context/ThemeContext';
import { TooltipProvider } from '@/components/ui/tooltip';
import { Layout } from './components/Layout';
import { JsonTools } from './components/JsonTools';
import { LoremTools } from './components/LoremTools';
import { Settings } from './components/Settings';
import { StringTools } from './components/StringTools';

function App() {
  return (
    <ThemeProvider>
      <TooltipProvider>
      <BrowserRouter>
        <Routes>
          <Route path="/" element={<Layout />}>
            <Route index element={<Navigate to="/tools/string/url-encode" replace />} />
            <Route path="settings" element={<Settings />} />
            <Route path="tools/string/url-encode" element={<StringTools tool="url-encode" />} />
            <Route path="tools/string/url-decode" element={<StringTools tool="url-decode" />} />
            <Route path="tools/string/base64-encode" element={<StringTools tool="base64-encode" />} />
            <Route path="tools/string/base64-decode" element={<StringTools tool="base64-decode" />} />
            <Route path="tools/string/trim" element={<StringTools tool="trim" />} />
            <Route path="tools/string/upper-case" element={<StringTools tool="upper-case" />} />
            <Route path="tools/string/lower-case" element={<StringTools tool="lower-case" />} />
            <Route path="tools/string/capital-case" element={<StringTools tool="capital-case" />} />
            <Route path="tools/string/snake-case" element={<StringTools tool="snake-case" />} />
            <Route path="tools/string/kebab-case" element={<StringTools tool="kebab-case" />} />
            <Route path="tools/string/camel-case" element={<StringTools tool="camel-case" />} />
            <Route path="tools/string/pascal-case" element={<StringTools tool="pascal-case" />} />
            <Route path="tools/string/sentence-case" element={<StringTools tool="sentence-case" />} />
            <Route path="tools/lorem-ipsum" element={<Navigate to="/tools/lorem-ipsum/generator" replace />} />
            <Route path="tools/lorem-ipsum/generator" element={<LoremTools tool="generator" />} />
            <Route path="tools/lorem-ipsum/characters" element={<LoremTools tool="characters" />} />
            <Route path="tools/lorem-ipsum/bytes" element={<LoremTools tool="bytes" />} />
            <Route path="tools/lorem-ipsum/title" element={<LoremTools tool="title" />} />
            <Route path="tools/lorem-ipsum/slug" element={<LoremTools tool="slug" />} />
            <Route path="tools/lorem-ipsum/camel-case" element={<LoremTools tool="camel-case" />} />
            <Route path="tools/lorem-ipsum/list" element={<LoremTools tool="list" />} />
            <Route path="tools/lorem-ipsum/headings" element={<LoremTools tool="headings" />} />
            <Route path="tools/lorem-ipsum/html" element={<LoremTools tool="html" />} />
            <Route path="tools/lorem-ipsum/markdown" element={<LoremTools tool="markdown" />} />
            <Route path="tools/lorem-ipsum/json" element={<LoremTools tool="json" />} />
            <Route path="tools/json" element={<Navigate to="/tools/json/format" replace />} />
            <Route path="tools/json/format" element={<JsonTools tool="format" />} />
            <Route path="tools/json/minify" element={<JsonTools tool="minify" />} />
            <Route path="tools/json/validate" element={<JsonTools tool="validate" />} />
            <Route path="tools/json/path" element={<JsonTools tool="path" />} />
            <Route path="tools/json/diff" element={<JsonTools tool="diff" />} />
            <Route path="*" element={<Navigate to="/tools/string/url-encode" replace />} />
          </Route>
        </Routes>
      </BrowserRouter>
      </TooltipProvider>
    </ThemeProvider>
  );
}

export default App;
