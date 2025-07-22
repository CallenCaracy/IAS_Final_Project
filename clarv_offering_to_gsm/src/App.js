// src/App.js
import React from 'react';
import { BrowserRouter as Router, Routes, Route, useNavigate } from 'react-router-dom';
import Chat from './pages/chat';

function Home() {
  const navigate = useNavigate();

  return (
    <div style={{ padding: 20 }}>
      <h1>Welcome to Encryption App</h1>
      <button onClick={() => navigate('/chat')}>Go to Chat</button>
    </div>
  );
}

export default function App() {
  return (
    <Router>
      <Routes>
        <Route path="/" element={<Home />} />
        <Route path="/chat" element={<Chat />} />
      </Routes>
    </Router>
  );
}
