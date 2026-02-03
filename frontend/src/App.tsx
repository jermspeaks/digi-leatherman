import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import { ThemeProvider } from './context/ThemeContext';
import { TooltipProvider } from '@/components/ui/tooltip';
import { Layout } from './components/Layout';
import { StringTools } from './components/StringTools';

function App() {
  return (
    <ThemeProvider>
      <TooltipProvider>
      <BrowserRouter>
        <Routes>
          <Route path="/" element={<Layout />}>
            <Route index element={<Navigate to="/tools/string/url-encode" replace />} />
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
            <Route path="*" element={<Navigate to="/tools/string/url-encode" replace />} />
          </Route>
        </Routes>
      </BrowserRouter>
      </TooltipProvider>
    </ThemeProvider>
  );
}

export default App;
