import React, { useEffect, useRef, useState } from "react";
import axios from "axios";
import "../css/chat.css"; 
import men from "../assets/boys.jpg";

const wsURL = "ws://192.168.137.97:8080/ws";
const apiBase = "http://192.168.137.97:8080";
// const wsURL = "ws://localhost:8080/ws";
// const apiBase = "http://localhost:8080";
const secretKey = "TESTTEST";

const Chat = () => {
  const [chat, setChat] = useState([]);
  const [msg, setMsg] = useState("");
  const ws = useRef(null);

    useEffect(() => {
    ws.current = new WebSocket(wsURL);

    ws.current.onmessage = async (event) => {
        try {
        // event.data is now a base64 string
        const base64 = event.data;

        // decode base64 to binary string
        const binaryString = atob(base64);

        // convert binary string to hex string
        const hex = [...binaryString]
            .map((c) => c.charCodeAt(0).toString(16).padStart(2, "0"))
            .join("");

        const res = await axios.post(`${apiBase}/decrypt`, {
            cipherHex: hex,
            key: secretKey,
        });

        setChat((prev) => [
            ...prev,
            {
            sender: "Peer",
            text: res.data.plainText,
            time: new Date().toLocaleTimeString([], { hour: "2-digit", minute: "2-digit" }),
            },
        ]);
        } catch (err) {
        console.error("Decryption failed:", err);
        setChat((prev) => [...prev, { sender: "System", text: "⚠️ Failed to decrypt message", time: "" }]);
        }
    };

    return () => ws.current?.close();
    }, []);

    const send = async () => {
    if (!msg) return;
    try {
        const res = await axios.post(`${apiBase}/encrypt`, {
        plainText: msg,
        key: secretKey,
        });
        const hex = res.data.cipherHex;
        // convert hex string to Uint8Array
        const byteArray = new Uint8Array(hex.match(/.{1,2}/g).map((byte) => parseInt(byte, 16)));
        // convert Uint8Array to base64 string
        const base64 = btoa(String.fromCharCode(...byteArray));
        ws.current.send(base64);

        setChat((prev) => [
        ...prev,
        {
            sender: "You",
            text: msg,
            time: new Date().toLocaleTimeString([], { hour: "2-digit", minute: "2-digit" }),
        },
        ]);
        setMsg("");
    } catch {
        setChat((prev) => [...prev, { sender: "System", text: "⚠️ Failed to encrypt message", time: "" }]);
    }
    };

  return (
    <>
    <div className="chat-container">
        <div className="chat-header">
            <img src={men} alt="Chat Banner" className="chat-banner" />
            <h2 className="chat-title">Sekreto sa panag-igsoonay</h2>
        </div>
      <div className="chat-box">
        {chat.map((msgObj, idx) => {
          const isYou = msgObj.sender === "You";
          return (
            <div key={idx} className={`chat-message ${isYou ? "you" : ""}`}>
              <span className={`chat-bubble ${isYou ? "you" : ""}`}>
                <strong>{msgObj.sender}:</strong> {msgObj.text}
                <div className="chat-timestamp">{msgObj.time}</div>
              </span>
            </div>
          );
        })}
      </div>

      <div className="chat-input-area">
        <input
          type="text"
          value={msg}
          placeholder="Type a message..."
          onChange={(e) => setMsg(e.target.value)}
          className="chat-input"
        />
        <button onClick={send} className="chat-send-btn">
          Send
        </button>
      </div>
    </div>
    </>
  );
};

export default Chat;
