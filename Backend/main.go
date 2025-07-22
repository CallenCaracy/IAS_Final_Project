package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/rs/cors"
)

var cipher *Cipher

func main() {
	cipher = NewCipher()

	http.HandleFunc("/encrypt", func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			PlainText string `json:"plainText"`
			Key       string `json:"key"`
		}
		json.NewDecoder(r.Body).Decode(&req)
		encrypted, err := cipher.Encrypt(req.PlainText, req.Key)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		resp := map[string]string{"cipherHex": fmt.Sprintf("%x", encrypted)}
		json.NewEncoder(w).Encode(resp)
	})

	http.HandleFunc("/decrypt", func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			CipherHex string `json:"cipherHex"`
			Key       string `json:"key"`
		}
		json.NewDecoder(r.Body).Decode(&req)
		cipherBytes, _ := hex.DecodeString(req.CipherHex)
		decrypted, err := cipher.Decrypt(cipherBytes, req.Key)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		resp := map[string]string{"plainText": decrypted}
		json.NewEncoder(w).Encode(resp)
	})

	// 👇 CORS middleware wrapping the default mux
	corsHandler := cors.New(cors.Options{
		AllowedOrigins: []string{
			"http://localhost:3000",
			"http://localhost:5173",
		},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type"},
		AllowCredentials: true,
	}).Handler(http.DefaultServeMux)

	log.Println("Go backend running on :8080")
	log.Fatal(http.ListenAndServe(":8080", corsHandler))
}
