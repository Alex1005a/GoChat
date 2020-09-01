package hasher

import "testing"

func TestHasher_CheckPassword(t *testing.T) {
	hasher := NewHasher()
	password := "pass"
	hash, _ := hasher.PasswordToHash(password)
    t.Log(hash)
	if hasher.CheckPassword(password, hash) == false{
		t.Errorf("Check password is false")
	}

}