// +build all travis

package xlm

import (
	"testing"
)

func TestXLM(t *testing.T) {
	seed, address, err := GetKeyPair()
	if err != nil {
		t.Fatal(err)
	}
	err = GetXLM(address)
	if err != nil {
		t.Fatal(err)
	}
	err = GetXLM("blah")
	if err == nil {
		t.Fatal("Invalid Address, shouldn't work!")
	}
	// now we have coins, so we can send to another account
	// need to create an account through
	_, destPubKey, err := GetKeyPair()
	if err != nil {
		t.Fatal(err)
	}
	_, _, err = SendXLMCreateAccount(destPubKey, "2", "wrong seed")
	// create the destiantion account by sendignb some coins to bootstrap
	if err == nil {
		t.Fatalf("Wrong seed, shouldn't work!")
	}
	_, _, err = SendXLMCreateAccount("wrong pubkey", "2", seed)
	// create the destiantion account by sendignb some coins to bootstrap
	if err == nil {
		t.Fatalf("Wrong pubkey, shouldn't work!")
	}
	_, _, err = SendXLMCreateAccount(destPubKey, "2", seed)
	// create the destiantion account by sendignb some coins to bootstrap
	if err != nil {
		t.Fatal(err)
	}
	_, _, err = SendXLM(destPubKey, "1", seed, "")
	if err != nil {
		t.Fatal(err)
	}
	_, _, err = SendXLM(destPubKey, "1", "wrong seed", "")
	if err == nil {
		t.Fatalf("Wrong seed, shouldn't work!")
	}
	_, _, err = SendXLM("wrong pubkey", "1", seed, "")
	if err == nil {
		t.Fatalf("Wrong pubkey, shouldn't work!")
	}
	// don't test the reverse becuase apparently there's some problem in catching the next block
	// or something
}
