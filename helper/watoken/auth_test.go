package watoken

import (
	"fmt"
	"testing"
)

func TestTokenWa(t *testing.T) {
	req, err := EncodeforHours("628987140700", "Ulbi", "e4cb06d20bcce42bf4ac16c9b056bfaf1c6a5168c24692b38eb46d551777dc4147db091df55d64499fdf2ca85504ac4d320c4c645c9bef75efac0494314cae94", 43830)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	fmt.Println(req)
	tes, err := Decode("47db091df55d64499fdf2ca85504ac4d320c4c645c9bef75efac0494314cae94", req)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	fmt.Println(tes)
}

// func HashApiLogin(t *testing.T) {
// 	req, err := HashAndEncodeBase64("wss://gw.ulbi.ac.id/api/whatsauth/request")
// 	if err != nil {
// 		t.Errorf("Error: %v", err)
// 	}
// 	t.Logf("Hashed and Base64 Encoded: %s", req)
// }

// func TestMinute(t *testing.T) {
// 	fmt.Printf("TestMinute")
// 	req, err := watoken.EncodeforMinutes("6282184952582", "e4cb06d20bcce42bf4ac16c9b056bfaf1c6a5168c24692b38eb46d551777dc4147db091df55d64499fdf2ca85504ac4d320c4c645c9bef75efac0494314cae94", 120)
// 	if err != nil {
// 		t.Errorf("Error: %v", err)
// 	}
// 	fmt.Printf("Testing: %s", req)
// }

