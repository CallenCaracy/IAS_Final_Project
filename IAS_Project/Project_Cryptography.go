package main

import (
	"bufio"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"strings"
)

func main() {
	var rsaPrivateKey *rsa.PrivateKey
	var rsaPublicKey rsa.PublicKey
	var choice int
	reader := bufio.NewReader(os.Stdin)
	var Vigenère [26][26]rune

	Vigenère = buildVigenereTable(Vigenère)

	for {
		displayOptions()
		fmt.Scan(&choice)
		reader.ReadString('\n')

		switch choice {
		case 1:
			var cipherTextNo1 string
			var cipherTextNo2 string
			var cipherTextNo3 []byte
			var cipherTextNo4 []byte
			var cipherTextNo5 []byte
			var keyStream []rune

			plainText := readLine("Enter text to encrypt: ")
			keyText := readLine("Enter text as key: ")

			plainText = normalizeText(plainText)
			keyText = normalizeText(keyText)
			keyStream = makeKeyStream(keyText, plainText)
			rsaPrivateKey, rsaPublicKey = generateRSAKeys()

			fmt.Println("This is the keystream: ", keyStream)

			cipherTextNo1 = shiftingByFourEncryption(plainText, 4)
			cipherTextNo2 = vigenereEncrypt(cipherTextNo1, keyStream, Vigenère)
			cipherTextNo3 = vernamXOREncode(cipherTextNo2, keyStream)
			cipherTextNo4 = railFenceEncryptBytes(cipherTextNo3, 3)
			cipherTextNo5 = rsaEncrypt(rsaPublicKey, cipherTextNo4)

			fmt.Println("Cipher Text Layer 1:", cipherTextNo1)
			fmt.Println("Cipher Text Layer 2:", cipherTextNo2)
			fmt.Println("Cipher Text Layer 3:", cipherTextNo3)
			fmt.Println("Cipher Text Layer 4:", cipherTextNo4)
			fmt.Printf("Cipher Text Layer 5: %x\n", cipherTextNo5)
			fmt.Println()

		case 2:
			var decryptedRSA []byte
			var keyStream []rune

			cipherHex := readLine("Enter RSA-encrypted hex text: ")
			keyText := readLine("Enter text as key: ")

			cipherBytes, _ := hex.DecodeString(cipherHex)
			keyText = normalizeText(keyText)

			decryptedRSA = rsaDecrypt(rsaPrivateKey, cipherBytes)
			decryptedRailFenceBytes := railFenceDecryptBytes([]byte(decryptedRSA), 3)
			keyStream = makeKeyStreamForBytes(keyText, len(decryptedRailFenceBytes))
			vernamDecrypted := vernamXORDecode(decryptedRailFenceBytes, keyStream)
			vigenereDecrypted := vigenereDecrypt(vernamDecrypted, keyStream, Vigenère)
			finalPlainText := shiftingByFourDencryption(vigenereDecrypted, 4)

			fmt.Println("Decrypted RSA (Layer 5):", decryptedRSA)
			fmt.Println("Transpositional Cipher (Layer 4):", decryptedRailFenceBytes)
			fmt.Println("Vernam Cipher (Layer 3):", vernamDecrypted)
			fmt.Println("Vigenere Cipher (Layer 2):", vigenereDecrypted)
			fmt.Println("Monoalphabetic Cipher Cipher (Layer 1):", finalPlainText)
		case 3:
			return

		default:
			fmt.Println("Invalid choice.")
			fmt.Println()
		}
	}
}

// This are helper functions
func displayOptions() {
	fmt.Println("Project: One Day I Am Gonna Grow Wings Cipher")
	fmt.Println("1. Encrypt A Text")
	fmt.Println("2. Decrypt A Text")
	fmt.Println("3. Exit")
	fmt.Print("Pick your choice: ")
}

func readLine(prompt string) string {
	fmt.Print(prompt)
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	return strings.TrimSpace(text)
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
	text = strings.ReplaceAll(text, " ", "")
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

func makeKeyStreamForBytes(key string, length int) []rune {
	key = normalizeText(key)
	keyRunes := []rune(key)
	keyStream := make([]rune, length)
	for i := 0; i < length; i++ {
		keyStream[i] = keyRunes[i%len(keyRunes)]
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
// func vernamXOR(text string, keyStream []rune) string {
// 	var result strings.Builder

// 	for i, char := range text {
// 		if char >= 'A' && char <= 'Z' {
// 			xored := (char - 'A') ^ (keyStream[i] - 'A')
// 			result.WriteRune(xored + 'A')
// 		} else {
// 			result.WriteRune(char)
// 		}
// 	}

// 	return result.String()
// }

func vernamXOREncode(text string, keyStream []rune) []byte {
	xorBytes := make([]byte, len(text))
	for i, char := range text {
		xorBytes[i] = byte(char) ^ byte(keyStream[i])
	}
	return xorBytes
}

func vernamXORDecode(xorBytes []byte, keyStream []rune) string {
	result := make([]rune, len(xorBytes))
	for i, b := range xorBytes {
		result[i] = rune(b ^ byte(keyStream[i]))
	}
	return string(result)
}

// Transpositional Cipher
// func railFenceEncrypt(text string, rails int) string {
// 	if rails <= 1 || len(text) <= rails {
// 		return text
// 	}

// 	rail := make([][]rune, rails)
// 	row := 0
// 	down := true

// 	for _, char := range text {
// 		rail[row] = append(rail[row], char)

// 		if down {
// 			row++
// 			if row == rails-1 {
// 				down = false
// 			}
// 		} else {
// 			row--
// 			if row == 0 {
// 				down = true
// 			}
// 		}
// 	}

// 	var result strings.Builder
// 	for _, r := range rail {
// 		result.WriteString(string(r))
// 	}
// 	return result.String()
// }

func railFenceEncryptBytes(data []byte, rails int) []byte {
	if rails <= 1 || len(data) <= rails {
		return data
	}

	rail := make([][]byte, rails)
	row := 0
	down := true

	for _, b := range data {
		rail[row] = append(rail[row], b)
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

	var result []byte
	for _, r := range rail {
		result = append(result, r...)
	}
	return result
}

// func railFenceDecrypt(cipher string, rails int) string {
// 	if rails <= 1 || len(cipher) <= rails {
// 		return cipher
// 	}

// 	// Initialize the rail pattern
// 	pattern := make([]int, len(cipher))
// 	row, down := 0, true
// 	for i := range pattern {
// 		pattern[i] = row
// 		if down {
// 			row++
// 			if row == rails-1 {
// 				down = false
// 			}
// 		} else {
// 			row--
// 			if row == 0 {
// 				down = true
// 			}
// 		}
// 	}

// 	count := make([]int, rails)
// 	for _, r := range pattern {
// 		count[r]++
// 	}

// 	railsData := make([][]rune, rails)
// 	idx := 0
// 	for r := 0; r < rails; r++ {
// 		railsData[r] = []rune(cipher[idx : idx+count[r]])
// 		idx += count[r]
// 	}

// 	result := make([]rune, len(cipher))
// 	railIndex := make([]int, rails)
// 	for i, r := range pattern {
// 		result[i] = railsData[r][railIndex[r]]
// 		railIndex[r]++
// 	}

// 	return string(result)
// }

func railFenceDecryptBytes(cipher []byte, rails int) []byte {
	if rails <= 1 || len(cipher) <= rails {
		return cipher
	}

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

	railsData := make([][]byte, rails)
	idx := 0
	for r := 0; r < rails; r++ {
		railsData[r] = cipher[idx : idx+count[r]]
		idx += count[r]
	}

	result := make([]byte, len(cipher))
	railIndex := make([]int, rails)
	for i, r := range pattern {
		result[i] = railsData[r][railIndex[r]]
		railIndex[r]++
	}

	return result
}

// RSA Algorithm
func generateRSAKeys() (*rsa.PrivateKey, rsa.PublicKey) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(err)
	}
	return privateKey, privateKey.PublicKey
}

func rsaEncrypt(publicKey rsa.PublicKey, plaintext []byte) []byte {
	ciphertext, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, &publicKey, plaintext, nil)
	if err != nil {
		panic(err)
	}
	return ciphertext
}

func rsaDecrypt(privateKey *rsa.PrivateKey, ciphertext []byte) []byte {
	plaintext, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, privateKey, ciphertext, nil)
	if err != nil {
		panic(err)
	}
	return plaintext
}
