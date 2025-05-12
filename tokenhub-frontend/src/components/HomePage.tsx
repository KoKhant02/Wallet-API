import { useNavigate } from 'react-router-dom';

function HomePage() {
  const navigate = useNavigate();

  const handleButtonClick = (route: string) => {
    navigate(route);
  };

  return (
    <div className="container mt-5">
      <h2 className="mb-4 text-start">Welcome to TokenHub</h2>
      
      <h3 className="text-start mb-4" style={{ fontSize: '1.5rem', fontWeight: '500' }}>Balance Check</h3>
      <div className="row justify-content-center">
        <div className="col-md-4 mb-3">
          <div 
            className="card border-0 shadow-lg hover-shadow"
            style={{ cursor: 'pointer', borderRadius: '10px' }} 
            onClick={() => handleButtonClick('/nft-balance')}
          >
            <div className="card-body text-center">
              <i className="bi bi-file-earmark-zip-fill" style={{ fontSize: '2rem', color: '#007bff' }}></i>
              <h5 className="card-title mt-3">NFT (ERC721)</h5>
              <button className="btn btn-primary w-100 mt-3">Check Now!</button>
            </div>
          </div>
        </div>
        <div className="col-md-4 mb-3">
          <div 
            className="card border-0 shadow-lg hover-shadow"
            style={{ cursor: 'pointer', borderRadius: '10px' }} 
            onClick={() => handleButtonClick('/erc20-balance')}
          >
            <div className="card-body text-center">
              <i className="bi bi-currency-dollar" style={{ fontSize: '2rem', color: '#28a745' }}></i>
              <h5 className="card-title mt-3">Token (ERC20)</h5>
              <button className="btn btn-primary w-100 mt-3">Check Now!</button>
            </div>
          </div>
        </div>
        <div className="col-md-4 mb-3">
          <div 
            className="card border-0 shadow-lg hover-shadow"
            style={{ cursor: 'pointer', borderRadius: '10px' }} 
            onClick={() => handleButtonClick('/erc1155-balance')}
          >
            <div className="card-body text-center">
              <i className="bi bi-stack" style={{ fontSize: '2rem', color: '#ff7f50' }}></i>
              <h5 className="card-title mt-3">ERC1155</h5>
              <button className="btn btn-primary w-100 mt-3">Check Now!</button>
            </div>
          </div>
        </div>
      </div>

    </div>
  );
}

export default HomePage;
