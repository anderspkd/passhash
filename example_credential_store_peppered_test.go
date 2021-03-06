package passhash

import (
	"bytes"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
)

var pepperGrinder PepperGrinder

type PepperGrinder struct {
	cipher cipher.Block
}

func (p PepperGrinder) Encrypt(text []byte) (nonce, encrypted []byte) {
	nonce = make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err.Error())
	}
	aesgcm, err := cipher.NewGCM(p.cipher)
	if err != nil {
		panic(err.Error())
	}
	return nonce, aesgcm.Seal(nil, nonce, text, nil)
}

func (p PepperGrinder) Decrypt(nonce, text []byte) []byte {
	aesgcm, err := cipher.NewGCM(p.cipher)
	if err != nil {
		panic(err.Error())
	}
	decryptedText, err := aesgcm.Open(nil, nonce, text, nil)
	if err != nil {
		panic(err.Error())
	}
	return decryptedText
}

func NewPepperGrinder(key []byte) (PepperGrinder, error) {
	cipher, err := aes.NewCipher(key)
	if err != nil {
		return PepperGrinder{}, err
	}
	return PepperGrinder{cipher}, nil
}

func init() {
	var err error
	pepperGrinder, err = NewPepperGrinder([]byte(KeyForPepperID[DefaultPepperID]))
	if err != nil {
		panic(err.Error())
	}
}

type PepperID uint

const DefaultPepperID PepperID = 1

var KeyForPepperID = map[PepperID][]byte{
	1: []byte("AES256Key-32Characters1234567890"),
}

type StringCredentialPepperedStore struct {
	StringCredentialStore
	PepperID PepperID // Versions pepper key/method
	Nonce    []byte
}

func (store *StringCredentialPepperedStore) Store(credential *Credential) error {
	store.PepperID = DefaultPepperID
	store.Nonce, credential.Hash = pepperGrinder.Encrypt(credential.Hash)
	return store.StringCredentialStore.Store(credential)
}

func (store *StringCredentialPepperedStore) StoreContext(ctx context.Context, credential *Credential) error {
	return store.Store(credential)
}

func (store *StringCredentialPepperedStore) Load(id UserID) (*Credential, error) {
	credential, err := store.StringCredentialStore.Load(id)
	if err != nil {
		return nil, err
	}
	switch store.PepperID {
	case 1:
		credential.Hash = pepperGrinder.Decrypt(store.Nonce, credential.Hash)
	default:
		return nil, fmt.Errorf("Unsupported PepperID %v", store.PepperID)
	}
	return credential, nil
}

func (store *StringCredentialPepperedStore) LoadContext(id UserID) (*Credential, error) {
	return store.Load(id)
}

func ExampleCredentialStore_peppered() {
	userID := UserID(0)
	password := "insecurepassword"
	origCredential, err := NewCredential(userID, password)
	if err != nil {
		fmt.Println("Error creating credential.", err)
		return
	}

	store := StringCredentialPepperedStore{}
	store.Store(origCredential)
	newCredential, err := store.Load(userID)
	if err != nil {
		fmt.Println("Error loading credential.", err)
		return
	}

	credentialEqual := newCredential == origCredential
	kdfEqual := newCredential.Kdf == origCredential.Kdf
	cfEqual := newCredential.WorkFactor == origCredential.WorkFactor // Not equal due to pointer comparison
	saltEqual := bytes.Compare(newCredential.Salt, origCredential.Salt) == 0
	hashEqual := bytes.Compare(newCredential.Hash, origCredential.Hash) == 0
	matched, updated := newCredential.MatchesPassword(password)
	fmt.Println("credentialEqual:", credentialEqual)
	fmt.Println("kdfEqual:", kdfEqual)
	fmt.Println("cfEqual:", cfEqual)
	fmt.Println("saltEqual:", saltEqual)
	fmt.Println("hashEqual:", hashEqual) // Not equal due to peppering
	fmt.Println("newCredential.MatchesPassword (matched):", matched)
	fmt.Println("newCredential.MatchesPassword (updated):", updated)

	// Output:
	// credentialEqual: false
	// kdfEqual: true
	// cfEqual: false
	// saltEqual: true
	// hashEqual: false
	// newCredential.MatchesPassword (matched): true
	// newCredential.MatchesPassword (updated): false
}
