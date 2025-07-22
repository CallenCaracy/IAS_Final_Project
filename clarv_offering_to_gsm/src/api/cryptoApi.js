export async function encrypt(text, key) {
  const response = await fetch('http://localhost:8080/encrypt', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ plainText: text, key: key }),
  });
  const data = await response.json();
  return data.cipherHex;
}

export async function decrypt(cipherHex, key) {
  const response = await fetch('http://localhost:8080/decrypt', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ cipherHex: cipherHex, key: key }),
  });
  const data = await response.json();
  return data.plainText;
}
