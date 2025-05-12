import { Route, BrowserRouter as Router, Routes } from 'react-router-dom';
import './App.css';
import ERC1155Balance from './components/ERC1155Balance';
import ERC20Balance from './components/ERC20Balance';
import ERC721Balance from './components/ERC721Balance';
import HomePage from './components/HomePage';

function App() {
  return (
    <Router>
      <div className="App">
        <Routes>
          <Route path="/" element={<HomePage />} />
          <Route path="/erc20-balance" element={<ERC20Balance />} />
          <Route path="/nft-balance" element={<ERC721Balance />} />
          <Route path="/erc1155-balance" element={<ERC1155Balance />} />
        </Routes>
      </div>
    </Router>
  );
}

export default App;
