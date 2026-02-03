import { StringTools } from './components/StringTools';
import './App.css';

function App() {
  return (
    <>
      <header>
        <h1>Digi Leatherman</h1>
        <nav className="main-nav" aria-label="Tool categories">
          <span className="nav-item active">String tools</span>
          <span className="nav-item disabled">Other tools</span>
        </nav>
      </header>
      <main>
        <StringTools />
      </main>
    </>
  );
}

export default App;
