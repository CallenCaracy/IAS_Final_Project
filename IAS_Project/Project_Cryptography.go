package main

import (
	"bufio"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"fmt"
	"os"
	"strings"
)

func main() {
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
			var cipherTextNo3 string
			var cipherTextNo4 string
			var cipherTextNo5 []byte
			var keyStream []rune

			plainText := readLine("Enter text to encrypt: ")
			keyText := readLine("Enter text as key: ")

			plainText = normalizeText(plainText)
			keyText = normalizeText(keyText)
			keyStream = makeKeyStream(keyText, plainText)
			privateKey, publicKey := generateRSAKeys()

			fmt.Println("This is the keystream: ", keyStream)
			fmt.Println("This is the privateKey: ", privateKey)
			fmt.Println("This is the publicKey: ", publicKey)

			cipherTextNo1 = shiftingByFourEncryption(plainText, 4)
			cipherTextNo2 = vigenereEncrypt(cipherTextNo1, keyStream, Vigenère)
			cipherTextNo3 = vernamXOR(cipherTextNo2, keyStream)
			cipherTextNo4 = railFenceEncrypt(cipherTextNo3, 3)
			cipherTextNo5 = rsaEncrypt(publicKey, cipherTextNo4)

			fmt.Println("Cipher Text Layer 1:", cipherTextNo1)
			fmt.Println("Cipher Text Layer 2:", cipherTextNo2)
			fmt.Println("Cipher Text Layer 3:", cipherTextNo3)
			fmt.Println("Cipher Text Layer 4:", cipherTextNo4)
			fmt.Printf("Cipher Text Layer 5 (RSA): %x\n", cipherTextNo5)
			fmt.Println()

		case 2:
			// Later, implement decryption logic

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
