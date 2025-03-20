import React, { useState, useEffect, useRef } from 'react';

// Custom hook for grid calculations (with backend dimensions)
const useGridCalculation = (gridContainerRef, tileSize, backendData) => {
  const [gridSize, setGridSize] = useState(8); // Default size until backend data arrives

  useEffect(() => {
    if (backendData && backendData.rows && backendData.columns) {
      // Use grid dimensions from backend
      const maxDimension = Math.max(backendData.rows, backendData.columns);
      setGridSize(maxDimension);
    }
  }, [backendData]);

  // Calculate optimal tile size based on screen dimensions and grid size
  const calculateOptimalTileSize = () => {
    if (!gridContainerRef.current || !backendData) return tileSize;
    
    const containerWidth = gridContainerRef.current.clientWidth;
    const containerHeight = window.innerHeight - 300; // Account for header, controls, and input bar
    
    const maxWidth = Math.floor(containerWidth / gridSize);
    const maxHeight = Math.floor(containerHeight / gridSize);
    
    // Return the smaller of the two dimensions to ensure the grid fits
    return Math.min(maxWidth, maxHeight, tileSize);
  };

  // Handle window resize
  useEffect(() => {
    const handleResize = () => {
      calculateOptimalTileSize();
    };
    
    // Add event listener
    window.addEventListener('resize', handleResize);
    
    // Cleanup
    return () => {
      window.removeEventListener('resize', handleResize);
    };
  }, [gridSize, tileSize]);

  return { gridSize, calculateOptimalTileSize };
};

// GridHeader Component
const GridHeader = () => {
  return (
    <div style={{ 
      fontSize: '1.5rem', 
      fontWeight: 'bold', 
      marginBottom: '1rem', 
      textAlign: 'center', 
      color: '#2563eb' 
    }}>
      Ezy Teach - Grid Display
    </div>
  );
};

// GridControls Component (simplified for backend data)
const GridControls = ({ 
  tileSize, 
  setTileSize,
  showCoordinates, 
  setShowCoordinates, 
  gridSize,
  matrixId 
}) => {
  return (
    <div style={{ 
      backgroundColor: 'white', 
      padding: '1rem', 
      borderRadius: '0.5rem', 
      boxShadow: '0 1px 3px rgba(0, 0, 0, 0.1)', 
      marginBottom: '1rem' 
    }}>
      <div style={{ 
        display: 'flex', 
        flexWrap: 'wrap', 
        alignItems: 'center', 
        gap: '1rem' 
      }}>
        <div>
          <label style={{ 
            display: 'block', 
            fontSize: '0.875rem', 
            fontWeight: '500', 
            color: '#374151', 
            marginBottom: '0.25rem' 
          }}>
            Tile Size
          </label>
          <input
            type="range"
            min="30"
            max="120"
            value={tileSize}
            onChange={(e) => setTileSize(Number(e.target.value))}
            style={{ width: '8rem' }}
          />
          <span style={{ 
            marginLeft: '0.5rem', 
            fontSize: '0.875rem', 
            color: '#4b5563' 
          }}>
            {tileSize}px
          </span>
        </div>
        
        <div>
          <label style={{ 
            display: 'block', 
            fontSize: '0.875rem', 
            fontWeight: '500', 
            color: '#374151', 
            marginBottom: '0.25rem' 
          }}>
            Show Coordinates
          </label>
          <div style={{ display: 'flex', alignItems: 'center' }}>
            <input
              type="checkbox"
              checked={showCoordinates}
              onChange={(e) => setShowCoordinates(e.target.checked)}
              style={{ marginRight: '0.5rem' }}
            />
            <span style={{ fontSize: '0.875rem', color: '#4b5563' }}>
              {showCoordinates ? 'Visible' : 'Hidden'}
            </span>
          </div>
        </div>
        
        <div style={{ marginLeft: 'auto' }}>
          <div style={{ 
            fontSize: '0.875rem', 
            fontWeight: '500', 
            color: '#374151' 
          }}>
            Current Matrix: <span style={{ fontWeight: 'bold' }}>{matrixId || 'None'}</span>
          </div>
          <div style={{ fontSize: '0.875rem', fontWeight: '500', color: '#374151' }}>
            Grid Size: <span style={{ fontWeight: 'bold' }}>{gridSize}x{gridSize}</span>
          </div>
          <div style={{ fontSize: '0.75rem', color: '#6b7280' }}>
            Total Cells: {gridSize * gridSize}
          </div>
        </div>
      </div>
    </div>
  );
};

// InputBar Component the User input bar / Equation input bar
const InputBar = ({ userInput, setUserInput, handleSendData }) => {
  return (
    <div style={{
      display: 'flex',
      justifyContent: 'center',
      alignItems: 'center',
      marginBottom: '1rem',
      padding: '0.5rem',
      backgroundColor: 'white',
      borderRadius: '0.5rem',
      boxShadow: '0 1px 3px rgba(0, 0, 0, 0.1)'
    }}>
      <input
        type="text"
        value={userInput}
        onChange={(e) => setUserInput(e.target.value)}
        placeholder="Type your message here..."
        style={{
          flex: '1',
          maxWidth: '500px',
          padding: '0.5rem',
          border: '1px solid #d1d5db',
          borderRadius: '0.25rem',
          marginRight: '0.5rem'
        }}
      />
      <button
        onClick={handleSendData}
        style={{
          backgroundColor: '#2563eb',
          color: 'white',
          border: 'none',
          borderRadius: '0.25rem',
          padding: '0.5rem 1rem',
          cursor: 'pointer',
          fontWeight: '500'
        }}
      >
        Send
      </button>
    </div>
  );
};

// GridItem Component
const GridItem = ({ id, content, tileSize, showCoordinates }) => {
  return (
    <div
      key={id}
      id={id}
      style={{
        position: 'relative',
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center',
        backgroundColor: '#eff6ff',
        border: '1px solid #bfdbfe',
        borderRadius: '0.375rem',
        transition: 'background-color 0.2s',
        width: `${tileSize}px`,
        height: `${tileSize}px`,
        fontSize: `${Math.max(16, tileSize * 0.5)}px`,
      }}
      onMouseOver={(e) => e.currentTarget.style.backgroundColor = '#dbeafe'}
      onMouseOut={(e) => e.currentTarget.style.backgroundColor = '#eff6ff'}
    >
      {showCoordinates && (
        <span style={{
          position: 'absolute',
          top: '0.25rem',
          left: '0.25rem',
          fontSize: '0.75rem',
          color: '#6b7280',
          fontFamily: 'monospace'
        }}>
          {id}
        </span>
      )}
      {content}
    </div>
  );
};

// GridContainer Component
const GridContainer = ({ gridContainerRef, gridSize, tileSize, backendData, showCoordinates }) => {
  // Create grid items from backend data
  const createGridItems = () => {
    const items = [];
    
    if (!backendData || !backendData.cells) {
      // If no backend data, create empty grid
      for (let row = 1; row <= gridSize; row++) {
        for (let col = 1; col <= gridSize; col++) {
          items.push({
            id: `${row}x${col}`,
            content: '',
          });
        }
      }
      return items;
    }
    
    // Create grid with data from backend
    const cells = backendData.cells;
    const rows = backendData.rows || gridSize;
    const columns = backendData.columns || gridSize;
    
    for (let row = 1; row <= rows; row++) {
      for (let col = 1; col <= columns; col++) {
        const id = `${row}x${col}`;
        const content = cells[id] || '';
        
        items.push({
          id: id,
          content: content,
        });
      }
    }
    
    return items;
  };
  
  const gridItems = createGridItems();
  
  return (
    <div 
      ref={gridContainerRef} 
      style={{ 
        flex: 1,
        backgroundColor: 'white',
        borderRadius: '0.5rem',
        boxShadow: '0 1px 3px rgba(0, 0, 0, 0.1)',
        overflow: 'auto',
        padding: '1rem',
        display: 'flex',
        justifyContent: 'center'
      }}
    >
      <div 
        style={{
          display: 'grid',
          gap: '0.25rem',
          gridTemplateColumns: `repeat(${backendData?.columns || gridSize}, ${tileSize}px)`,
          gridTemplateRows: `repeat(${backendData?.rows || gridSize}, ${tileSize}px)`,
        }}
      >
        {gridItems.map((item) => (
          <GridItem 
            key={item.id}
            id={item.id}
            content={item.content}
            tileSize={tileSize}
            showCoordinates={showCoordinates}
          />
        ))}
      </div>
    </div>
  );
};

// Main GridDisplay Component
const GridDisplay = () => {
  // State for grid configuration
  const [tileSize, setTileSize] = useState(60); // Default tile size in pixels
  const [showCoordinates, setShowCoordinates] = useState(true); // Toggle for showing coordinates
  const [userInput, setUserInput] = useState(''); // For the input bar
  const [backendData, setBackendData] = useState(null); // State to store backend data
  
  // Reference to the grid container
  const gridContainerRef = useRef(null);
  
  // Use custom hook for grid calculations with backend data
  const { gridSize } = useGridCalculation(
    gridContainerRef, 
    tileSize, 
    backendData
  );
  // Function to handle sending data to backend
const handleSendData = () => {
  if (!userInput.trim()) return;
  
  // Prepare data to send to backend
  const dataToSend = {
    message: userInput,
    timestamp: new Date().toISOString()
  };
  
  console.log('Sending data to backend:', dataToSend);
  
  // Send data to Go backend
  fetch('http://localhost:8080/api/grid-data', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(dataToSend),
  })
    .then(response => {
      if (!response.ok) {
        throw new Error(`HTTP error! Status: ${response.status}`);
      }
      return response.json();
    })
    .then(data => {
      console.log('Success:', data);
      setBackendData(data.data); // Access the matrix data from the nested data property
      setUserInput('');
    })
    .catch((error) => {
      console.error('Error:', error);
      alert(`Failed to connect to backend: ${error.message}`);
    });
};
  // // Function to handle sending data to backend
  // const handleSendData = () => {
  //   if (!userInput.trim()) return;
    
  //   // Prepare data to send to backend
  //   const dataToSend = {
  //     message: userInput,
  //     timestamp: new Date().toISOString()
  //   };
    
  //   console.log('Sending data to backend:', dataToSend);
    
  //   // Example of how to send data to Go backend
  //   fetch('http://localhost:8080/api/grid-data', {
  //     method: 'POST',
  //     headers: {
  //       'Content-Type': 'application/json',
  //     },
  //     body: JSON.stringify(dataToSend),
  //   })
  //     .then(response => response.json())
  //     .then(data => {
  //       console.log('Success:', data);
  //       setBackendData(data); // Store the received data from backend
  //       setUserInput(''); // Clear input after successful send
  //     })
  //     .catch((error) => {
  //       console.error('Error:', error);
  //       alert(`Failed to connect to backend: ${error.message}`);
  //       // For development/testing without backend
  //       // Simulate response using the format mentioned
  //       const mockResponse = {
  //         cells: {
  //           "1x1": "5", "1x2": "+", "1x3": "3", "1x4": "=",
  //           "2x1": "8", "2x2": "-", "2x3": "2", "2x4": "=", 
  //           "3x1": "6", "3x2": "*", "3x3": "4", "3x4": "=", 
  //           "4x1": "12", "4x2": "/", "4x3": "3", "4x4": "="
  //         },
  //         columns: 4,
  //         rows: 4,
  //         matrixId: "A"
  //       };
  //       setBackendData(mockResponse);
  //       setUserInput(''); // Clear input after send
  //     });
  // };
  
  return (
    <div style={{ 
      display: 'flex', 
      flexDirection: 'column', 
      height: '100vh', 
      backgroundColor: '#f3f4f6', 
      padding: '1rem' 
    }}>
      {/* Header */}
      <GridHeader />
      
      {/* Controls */}
      <GridControls
        tileSize={tileSize}
        setTileSize={setTileSize}
        showCoordinates={showCoordinates}
        setShowCoordinates={setShowCoordinates}
        gridSize={backendData?.rows || gridSize}
        matrixId={backendData?.matrixId}
      />
      
      {/* Input Bar */}
      <InputBar 
        userInput={userInput}
        setUserInput={setUserInput}
        handleSendData={handleSendData}
      />
      
      {/* Grid Container */}
      <GridContainer
        gridContainerRef={gridContainerRef}
        gridSize={gridSize}
        tileSize={tileSize}
        backendData={backendData}
        showCoordinates={showCoordinates}
      />
    </div>
  );
};

// App Component
const App = () => {
  return (
    <div style={{ fontFamily: 'sans-serif' }}>
      <GridDisplay />
    </div>
  );
};

export default App;