package util

import (
	"encoding/json"
	"fmt"
	"log"
	"time"
)

// PrettyPrint with marshaled json
func PrettyPrint(v interface{}) {
	out, err := json.MarshalIndent(v, "", "    ")
	if err != nil {
		log.Println(err.Error())
	}
	log.Println(string(out))
}

// time in a friendly format
func FriendlyString(duration time.Duration) string {
	if duration.Hours() >= 48 {
		return fmt.Sprintf("%.0f days ago", duration.Hours()/24)
	}
	if duration.Hours() >= 24 {
		return "1 day ago"
	}
	if duration.Hours() >= 2 {
		return fmt.Sprintf("%.0f hours ago", duration.Hours())
	}
	if duration.Hours() >= 1 {
		return "1 hour ago"
	}
	if duration.Minutes() >= 2 {
		return fmt.Sprintf("%.0f minutes ago", duration.Minutes())
	}
	if duration.Minutes() >= 1 {
		return "1 minute ago"
	}
	return "a couple seconds ago"
}
