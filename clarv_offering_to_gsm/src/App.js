import React from 'react';
import { BrowserRouter as Router, Routes, Route, useNavigate } from 'react-router-dom';
import Chat from './pages/chat';
import './App.css'; 

function Home() {
  const navigate = useNavigate();

  return (
    <div className="landing-container">
      <h1 className="landing-title">🔐 Secure P2P Chat</h1>
      <p className="landing-subtext">
        Welcome! Start a secure, end-to-end encrypted peer-to-peer chat using our custom encryption system.
      </p>
      <button
        className="landing-button"
        onClick={() => navigate('/chat')}
      >
        🚀 Enter Chat
      </button>
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
