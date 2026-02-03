import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import { ThemeProvider } from './context/ThemeContext';
import { Layout } from './components/Layout';
import { StringTools } from './components/StringTools';
import './App.css';

function App() {
  return (
    <ThemeProvider>
      <BrowserRouter>
        <Routes>
          <Route path="/" element={<Layout />}>
            <Route index element={<Navigate to="/tools/string/url-encode" replace />} />
            <Route path="tools/string/url-encode" element={<StringTools tool="url-encode" />} />
            <Route path="tools/string/url-decode" element={<StringTools tool="url-decode" />} />
            <Route path="*" element={<Navigate to="/tools/string/url-encode" replace />} />
          </Route>
        </Routes>
      </BrowserRouter>
    </ThemeProvider>
  );
}

export default App;
