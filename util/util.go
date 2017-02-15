package util

import (
	"encoding/json"
	"log"
)

// PrettyPrint with marshaled json
func PrettyPrint(v interface{}) {
	out, err := json.MarshalIndent(v, "", "    ")
	if err != nil {
		log.Println(err.Error())
	}
	log.Println(string(out))
}
