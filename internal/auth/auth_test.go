package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestPasswordHashing(t *testing.T) {
	tcases := []string{"12345", "myPAsword!", "MyPasword!"}
	for _, tcase := range tcases {
		hash, err := HashPassWord(tcase)
		if err != nil {
			t.Fatalf("Couldn't hash %s", tcase)
		}
		err = CheckPassWordHash(hash, tcase)
		if err != nil {
			t.Errorf("Hash and Password do not match. Hash: %s Password: %s", hash, tcase)
		}
	}
	hash, err := HashPassWord(tcases[1])
	if err != nil {
		t.Fatalf("Couldn't hash %s", tcases[1])
	}
	err = CheckPassWordHash(hash, tcases[2])
	if err == nil {
		t.Errorf("Hash and Password match eventhough they shouldn't. Hash: %s Password: %s", hash, tcases[2])
	}
}

func TestTokenValidation(t *testing.T) {
	id, err := uuid.NewRandom()
	if err != nil {
		t.Fatalf("couldn't generate uuid: %s", err)
	}
	ss, err := MakeJWT(id, "secret123", 2*time.Second)
	if err != nil {
		t.Fatalf("JWT Generation failed: %s", err)
	}
	ids, err := ValidateJWT(ss, "secret123")
	if err != nil {
		t.Fatalf("Validation failed: %s", err)
	}
	t.Logf("Validation successful for id: %s", ids)

	_, err = ValidateJWT(ss, "secret1234")
	if err == nil {
		t.Fatalf("Validation succeeded eventhough it should fail: %s", err)
	}
	time.Sleep(2 * time.Second)
	_, err = ValidateJWT(ss, "secret123")
	if err == nil {
		t.Fatalf("Validation succeeded eventhough the token shouldve expired: %s", err)
	}
}
