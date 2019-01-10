// +build all

package aes

// ignoring travis because it doesn't like the new file creation thing.
// TODO: fix this
import (
	"log"
	"testing"
)

func TestAes(t *testing.T) {
	password := "Cool"
	ciphertext, err := Encrypt([]byte("Hello World"), password)
	if err != nil {
		t.Fatal(err)
	}
	log.Printf("Encrypted: %x\n", ciphertext)
	plaintext, err := Decrypt(ciphertext, password)
	if err != nil {
		t.Fatal(err)
	}
	log.Printf("Decrypted: %s\n", plaintext)
	if string(plaintext) != "Hello World" {
		t.Errorf("Problem with Decryption")
	}
	_, err = Decrypt(ciphertext, "Notcool")
	if err == nil {
		t.Fatalf("Didn't catch error during decrpytion, exiting!")
	}
	_, err = Decrypt([]byte(""), "")
	if err == nil {
		t.Fatalf("Didn't catch error during decrpytion, exiting!")
	}
	data := []byte("This is test data")
	log.Println("ULKBKPASQWEQWEQWE")
	err = EncryptFile("test_files/text.txt", data, password)
	if err != nil {
		t.Fatal(err)
	}
	decryptedSlice, err := DecryptFile("test_files/text.txt", password)
	if err != nil {
		t.Fatal(err)
	}
	if string(decryptedSlice) != "This is test data" {
		t.Fatalf("Can't decrypt file, exiting!")
	}
	_, err = DecryptFile("test_files/text.txt", "Notcool")
	if err == nil {
		t.Fatal(err)
	}
}
