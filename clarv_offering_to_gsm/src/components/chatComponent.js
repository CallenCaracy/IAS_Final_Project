import React, { useState } from 'react';
import { encrypt, decrypt } from './api/cryptoApi';

function ChatComponent() {
  const [inputText, setInputText] = useState('');
  const [key, setKey] = useState('');
  const [encryptedText, setEncryptedText] = useState('');
  const [decryptedText, setDecryptedText] = useState('');

  const handleEncrypt = async () => {
    const cipherHex = await encrypt(inputText, key);
    setEncryptedText(cipherHex);
  };

  const handleDecrypt = async () => {
    const plainText = await decrypt(encryptedText, key);
    setDecryptedText(plainText);
  };

  return (
    <div>
      <input value={inputText} onChange={e => setInputText(e.target.value)} placeholder="Enter text" />
      <input value={key} onChange={e => setKey(e.target.value)} placeholder="Enter key" />
      <button onClick={handleEncrypt}>Encrypt</button>
      <button onClick={handleDecrypt}>Decrypt</button>
      <div>Encrypted: {encryptedText}</div>
      <div>Decrypted: {decryptedText}</div>
    </div>
  );
}

export default ChatComponent;
