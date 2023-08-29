package main

import (
	"bufio"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
)

func getKey() []byte {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Key: ")
	scanner.Scan()
	data := scanner.Text()
	hash := sha256.Sum256([]byte(data))

	return hash[:]
}

func encrypt(plaintextFile, encryptedFile string) error {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Note: ")
	scanner.Scan()
	data := scanner.Text()
	note := []byte(data)
	noteStr := fmt.Sprintf("%s\n", note)
	file, err := os.OpenFile(encryptedFile, os.O_CREATE|os.O_WRONLY, 0644)
	file.WriteString(noteStr)
	if err != nil {
		return err
	}
	//
	key := getKey()
	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	plaintext, err := ioutil.ReadFile(plaintextFile)
	if err != nil {
		return err
	}

	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	file.Write(ciphertext)
	return nil
}

func decrypt(encryptedFile, decryptedFile string) error {
	ciphertext, err := ioutil.ReadFile(encryptedFile)
	if err != nil {
		return err
	}

	newlineIndex := -1
	for i, b := range ciphertext {
		if b == '\n' {
			newlineIndex = i
			break
		}
	}
	if newlineIndex != -1 {
		fmt.Println("Note:", string(ciphertext[:newlineIndex]))
		ciphertext = ciphertext[newlineIndex+1:]
	}

	key := getKey()
	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	if len(ciphertext) < aes.BlockSize {
		return err
	}

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)

	err = ioutil.WriteFile(decryptedFile, ciphertext, 0644)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	if len(os.Args) == 1 {
		fmt.Printf("Usage: %s [e | d]\n", os.Args[0])
		return
	}
	action := os.Args[1]
	switch action {
	case "e":
		if err := encrypt("entry", "cipher"); err != nil {
			log.Fatalf("Encryption error: %v", err)
		}
	case "d":
		if err := decrypt("cipher", "entry"); err != nil {
			log.Fatalf("Decryption error: %v", err)
		}
	default:
		log.Fatalf("Unknown action %s. Use 'e' or 'd'.", action)
	}
}
