package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"fmt"
	"strings"
)

type Cipher struct {
	privateKey     *rsa.PrivateKey
	publicKey      rsa.PublicKey
	vigenereTable  [26][26]rune
	railFenceRails int
	monoShift      int
}

func NewCipher() *Cipher {
	priv, pub := generateRSAKeys()
	return &Cipher{
		privateKey:     priv,
		publicKey:      pub,
		vigenereTable:  buildVigenereTable([26][26]rune{}),
		railFenceRails: 3,
		monoShift:      4,
	}
}

func (c *Cipher) Encrypt(plainText string, key string) ([]byte, error) {
	plainText = normalizeText(plainText)
	key = normalizeText(key)
	keyStream := makeKeyStream(key, plainText)

	fmt.Println("This is the keystream: ", keyStream)
	fmt.Println("This is the key: ", key)
	fmt.Println("This is the plainText: ", plainText)

	ct1 := shiftingByFourEncryption(plainText, c.monoShift)
	ct2 := vigenereEncrypt(ct1, keyStream, c.vigenereTable)
	ct3 := vernamXOR(ct2, keyStream)
	ct4 := railFenceEncrypt(ct3, c.railFenceRails)
	ct5 := rsaEncrypt(c.publicKey, ct4)

	return ct5, nil
}

func (c *Cipher) Decrypt(cipherBytes []byte, key string) (string, error) {
	key = normalizeText(key)
	decryptedRSA := rsaDecrypt(c.privateKey, cipherBytes)
	decryptedRailFence := railFenceDecrypt(decryptedRSA, c.railFenceRails)
	keyStream := makeKeyStream(key, decryptedRailFence)
	vernamDecrypted := vernamXOR(decryptedRailFence, keyStream)
	vigenereDecrypted := vigenereDecrypt(vernamDecrypted, keyStream, c.vigenereTable)
	finalPlainText := shiftingByFourDencryption(vigenereDecrypted, c.monoShift)
	return finalPlainText, nil
}

// Monoalphabetic Cipher
func shiftingByFourEncryption(plainText string, shifting int) string {
	cipherText := ""
	for _, char := range plainText {
		if char >= 'a' && char <= 'z' {
			cipherText += string((char-'a'+rune(shifting)+26)%26 + 'a')
		} else if char >= 'A' && char <= 'Z' {
			cipherText += string((char-'A'+rune(shifting)+26)%26 + 'A')
		} else {
			cipherText += string(char)
		}
	}
	return cipherText
}

func shiftingByFourDencryption(cipherText string, shifting int) string {
	plainText := ""
	for _, char := range cipherText {
		if char >= 'a' && char <= 'z' {
			plainText += string((char-'a'-rune(shifting)+26)%26 + 'a')
		} else if char >= 'A' && char <= 'Z' {
			plainText += string((char-'A'-rune(shifting)+26)%26 + 'A')
		} else {
			plainText += string(char)
		}
	}
	return plainText
}

// Polyalphabetic Cipher
func normalizeText(text string) string {
	text = strings.ToUpper(strings.TrimSpace(text))
	return text
}

func makeKeyStream(key string, text string) []rune {
	var keyStream []rune
	for i, char := range text {
		if char >= 'A' && char <= 'Z' {
			keyStream = append(keyStream, rune(key[i%len(key)]))
		} else {
			keyStream = append(keyStream, char)
		}
	}
	return keyStream
}

// Vigenère
func buildVigenereTable(Vigenère [26][26]rune) [26][26]rune {
	fmt.Print("Generating Tabula Recta...\n")
	for i := 0; i < 26; i++ {
		for j := 0; j < 26; j++ {
			Vigenère[i][j] = rune((i+j)%26 + 'A')
		}
	}
	return Vigenère
}

func vigenereEncrypt(text string, keyStream []rune, table [26][26]rune) string {
	var result strings.Builder

	for i, char := range text {
		if char >= 'A' && char <= 'Z' {
			row := keyStream[i] - 'A'
			col := char - 'A'
			result.WriteRune(table[row][col])
		} else {
			result.WriteRune(char)
		}
	}

	return result.String()
}

func vigenereDecrypt(cipherText string, keyStream []rune, table [26][26]rune) string {
	var result strings.Builder

	for i, char := range cipherText {
		if char >= 'A' && char <= 'Z' {
			row := keyStream[i] - 'A'
			col := 0
			// Find column in row where the character matches
			for j := 0; j < 26; j++ {
				if table[row][j] == char {
					col = j
					break
				}
			}
			result.WriteRune(rune(col) + 'A')
		} else {
			result.WriteRune(char)
		}
	}

	return result.String()
}

// Vernam
func vernamXOR(text string, keyStream []rune) string {
	var result strings.Builder

	for i, char := range text {
		if char >= 'A' && char <= 'Z' {
			xored := (char - 'A') ^ (keyStream[i] - 'A')
			result.WriteRune(xored + 'A')
		} else {
			result.WriteRune(char)
		}
	}

	return result.String()
}

// Transpositional Cipher
func railFenceEncrypt(text string, rails int) string {
	if rails <= 1 || len(text) <= rails {
		return text
	}

	rail := make([][]rune, rails)
	row := 0
	down := true

	for _, char := range text {
		rail[row] = append(rail[row], char)

		if down {
			row++
			if row == rails-1 {
				down = false
			}
		} else {
			row--
			if row == 0 {
				down = true
			}
		}
	}

	var result strings.Builder
	for _, r := range rail {
		result.WriteString(string(r))
	}
	return result.String()
}

func railFenceDecrypt(cipher string, rails int) string {
	if rails <= 1 || len(cipher) <= rails {
		return cipher
	}

	// Initialize the rail pattern
	pattern := make([]int, len(cipher))
	row, down := 0, true
	for i := range pattern {
		pattern[i] = row
		if down {
			row++
			if row == rails-1 {
				down = false
			}
		} else {
			row--
			if row == 0 {
				down = true
			}
		}
	}

	count := make([]int, rails)
	for _, r := range pattern {
		count[r]++
	}

	railsData := make([][]rune, rails)
	idx := 0
	for r := 0; r < rails; r++ {
		railsData[r] = []rune(cipher[idx : idx+count[r]])
		idx += count[r]
	}

	result := make([]rune, len(cipher))
	railIndex := make([]int, rails)
	for i, r := range pattern {
		result[i] = railsData[r][railIndex[r]]
		railIndex[r]++
	}

	return string(result)
}

// RSA Algorithm
func generateRSAKeys() (*rsa.PrivateKey, rsa.PublicKey) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(err)
	}
	return privateKey, privateKey.PublicKey
}

func rsaEncrypt(publicKey rsa.PublicKey, plaintext string) []byte {
	ciphertext, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, &publicKey, []byte(plaintext), nil)
	if err != nil {
		panic(err)
	}
	return ciphertext
}

func rsaDecrypt(privateKey *rsa.PrivateKey, ciphertext []byte) string {
	plaintext, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, privateKey, ciphertext, nil)
	if err != nil {
		panic(err)
	}
	return string(plaintext)
}
