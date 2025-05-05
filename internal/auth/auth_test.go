package auth

import (
	"testing"
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
