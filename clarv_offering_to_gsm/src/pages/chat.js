import React, { useState } from 'react';
import { encrypt, decrypt } from '../api/cryptoApi';

const CHAT_KEY = 'OnlyChars';

export default function Chat() {
  const [message, setMessage] = useState('');
  const [chatLog, setChatLog] = useState([]);
  const [cipherHex, setCipherHex] = useState('');

  const handleSend = async () => {
    if (!message.trim()) return;

    try {
      const encrypted = await encrypt(message, CHAT_KEY);
      setCipherHex(encrypted);

      setChatLog((prev) => [
        ...prev,
        { type: 'user', text: message },
        { type: 'encrypted', text: encrypted },
      ]);

      setMessage('');
    } catch (err) {
      console.error('Encryption failed:', err);
    }
  };

  const handleDecrypt = async () => {
    try {
      const decrypted = await decrypt(cipherHex, CHAT_KEY);
      setChatLog((prev) => [
        ...prev,
        { type: 'decrypted', text: decrypted },
      ]);
    } catch (err) {
      console.error('Decryption failed:', err);
    }
  };

  return (
    <div style={styles.container}>
      <h2 style={styles.header}>🔐 Encrypted Chat</h2>
      <div style={styles.chatBox}>
        {chatLog.map((entry, index) => (
          <div
            key={index}
            style={{
              ...styles.message,
              ...(entry.type === 'user'
                ? styles.userMsg
                : entry.type === 'decrypted'
                ? styles.decryptedMsg
                : styles.encryptedMsg),
            }}
          >
            <strong>{entry.type}:</strong> {entry.text}
          </div>
        ))}
      </div>

      <div style={styles.inputArea}>
        <input
          value={message}
          onChange={(e) => setMessage(e.target.value)}
          placeholder="Type your message..."
          style={styles.input}
        />
        <button onClick={handleSend} style={styles.sendBtn}>
          Send
        </button>
        <button onClick={handleDecrypt} style={styles.decryptBtn}>
          Decrypt Last
        </button>
      </div>
    </div>
  );
}

const styles = {
  container: {
    padding: 20,
    maxWidth: 600,
    margin: '0 auto',
    fontFamily: 'Arial, sans-serif',
  },
  header: {
    textAlign: 'center',
  },
  chatBox: {
    border: '1px solid #ccc',
    borderRadius: 8,
    padding: 10,
    height: 300,
    overflowY: 'auto',
    backgroundColor: '#f9f9f9',
    marginBottom: 10,
  },
  message: {
    margin: '8px 0',
    padding: 8,
    borderRadius: 6,
  },
  userMsg: {
    backgroundColor: '#e1f5fe',
    textAlign: 'right',
  },
  encryptedMsg: {
    backgroundColor: '#ffecb3',
  },
  decryptedMsg: {
    backgroundColor: '#c8e6c9',
    fontStyle: 'italic',
  },
  inputArea: {
    display: 'flex',
    gap: 8,
  },
  input: {
    flex: 1,
    padding: 8,
    fontSize: 16,
  },
  sendBtn: {
    padding: '8px 16px',
    backgroundColor: '#2196f3',
    color: '#fff',
    border: 'none',
    borderRadius: 4,
  },
  decryptBtn: {
    padding: '8px 16px',
    backgroundColor: '#4caf50',
    color: '#fff',
    border: 'none',
    borderRadius: 4,
  },
};
