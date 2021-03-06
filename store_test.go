package passhash

import (
	"context"
	"testing"
)

func TestDummyCredentialStoreStore(t *testing.T) {
	store := DummyCredentialStore{}
	credential := &Credential{}
	if err := store.Store(credential); err != nil {
		t.Error("Got error storing credential.", err)
	}
}

func TestDummyCredentialStoreStoreContext(t *testing.T) {
	store := DummyCredentialStore{}
	credential := &Credential{}
	if err := store.StoreContext(context.Background(), credential); err != nil {
		t.Error("Got error storing credential.", err)
	}
}

func TestDummyCredentialStoreLoad(t *testing.T) {
	store := DummyCredentialStore{}
	userID := UserID(0)
	credential, err := store.Load(userID)
	if err == nil {
		t.Error("Got error loading credential.", err)
	}
	if credential != nil {
		t.Error("DummyCredentialStore provided credential.", credential)
	}
}

func TestDummyCredentialStoreLoadContext(t *testing.T) {
	store := DummyCredentialStore{}
	userID := UserID(0)
	credential, err := store.LoadContext(context.Background(), userID)
	if err == nil {
		t.Error("Got error loading credential.", err)
	}
	if credential != nil {
		t.Error("DummyCredentialStore provided credential.", credential)
	}
}
