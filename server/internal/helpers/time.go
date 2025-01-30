package helpers

import "time"

func Sleep(seconds int) {
	time.Sleep(time.Duration(seconds) * time.Second)
}

func GetFormattedTime() string {
	return time.Now().Format(time.RFC1123Z)
}
