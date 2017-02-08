package util

import (
	"encoding/json"
	"fmt"
)

// Check panics if err is not nil
func Check(err error) {
	if err != nil {
		panic(err)
	}
}

// PrettyPrint with marshaled json
func PrettyPrint(v interface{}) {
	out, err := json.MarshalIndent(v, "", "    ")
	Check(err)
	fmt.Println(string(out))
}
