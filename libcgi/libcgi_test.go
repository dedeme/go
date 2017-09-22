package libcgi

import (
	"fmt"
	"testing"
  "github.com/dedeme/go/libcgi/cryp"
)

func TestErr(t *testing.T) {

}

func ExampleErr() {
	r, err := b64.DecodeString("eyJlcnJvciI6IkEgZmFpbCJ9")
	if err != nil {
		panic(err.Error())
	}

	fmt.Println()
	fmt.Println(string(r))

	// {"error":"A fail"}
}

func ExampleOk() {
	Init("a", "b")
	rp := make(map[string]interface{})
	rp["name"] = "Peter"
	rp["age"] = 23

	Ok(rp)

	r, err := b64.DecodeString("eyJhZ2UiOjIzLCJlcnJvciI6IiIsIm5hbWUiOiJQZXRlciJ9")
	if err != nil {
		panic(err.Error())
	}

	fmt.Println()
fmt.Println(cryp.Key("deme", 300))
	fmt.Println(string(r))
	fmt.Print(cryp.Decryp("b",
		"m++Vu56Dq71/nX+ypbartbPToeKnw5eXf92vzK7Hq9u0uJ+ghMGBpq2Oir+5sYKL"))

	// Output: m++Vu56Dq71/nX+ypbartbPToeKnw5eXf92vzK7Hq9u0uJ+ghMGBpq2Oir+5sYKL
	// {"age":23,"error":"","name":"Peter"}
	// {"age":23,"error":"","name":"Peter"}
}
