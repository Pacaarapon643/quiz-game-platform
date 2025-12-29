package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
)

func main() {
	// สร้าง keys directory ถ้ายังไม่มี
	if err := os.MkdirAll("keys", 0755); err != nil {
		fmt.Printf("Error creating keys directory: %v\n", err)
		return
	}

	// สร้าง RSA private key (2048 bits)
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		fmt.Printf("Error generating private key: %v\n", err)
		return
	}

	// บันทึก private key
	privateFile, err := os.Create("keys/private.pem")
	if err != nil {
		fmt.Printf("Error creating private.pem: %v\n", err)
		return
	}
	defer privateFile.Close()

	privateBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	if err := pem.Encode(privateFile, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateBytes,
	}); err != nil {
		fmt.Printf("Error encoding private key: %v\n", err)
		return
	}

	// บันทึก public key
	publicFile, err := os.Create("keys/public.pem")
	if err != nil {
		fmt.Printf("Error creating public.pem: %v\n", err)
		return
	}
	defer publicFile.Close()

	publicBytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		fmt.Printf("Error marshaling public key: %v\n", err)
		return
	}

	if err := pem.Encode(publicFile, &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicBytes,
	}); err != nil {
		fmt.Printf("Error encoding public key: %v\n", err)
		return
	}

	fmt.Println("✅ RSA key pair generated successfully!")
	fmt.Println("📄 keys/private.pem - Keep this SECRET! (added to .gitignore)")
	fmt.Println("📄 keys/public.pem - Can be shared")
	fmt.Println("\n⚠️  IMPORTANT: Do not commit private.pem to git!")
}
