package libcgi

import (
	"testing"
	"fmt"
)

func TestErr(t *testing.T) {

}

func ExampleErr() {
	Err("A fail")

	r, err := enc.DecodeString("eyJlcnJvciI6IkEgZmFpbCJ9")
	if err != nil {
		panic(err.Error())
	}

  fmt.Println()
	fmt.Println(string(r))

	// Output: eyJlcnJvciI6IkEgZmFpbCJ9
  // {"error":"A fail"}
}

func ExampleOk() {
  rp := make(map[string]interface{})
  rp["name"] = "Peter"
  rp["age"] = 23

  Ok(rp)

	r, err := enc.DecodeString("eyJhZ2UiOjIzLCJlcnJvciI6IiIsIm5hbWUiOiJQZXRlciJ9")
	if err != nil {
		panic(err.Error())
	}

  fmt.Println()
	fmt.Println(string(r))

  // Output: eyJhZ2UiOjIzLCJlcnJvciI6IiIsIm5hbWUiOiJQZXRlciJ9
  // {"age":23,"error":"","name":"Peter"}
}
