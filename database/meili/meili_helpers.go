package meili

import "strings"

func format(username string) string {
	username = strings.ToLower(username)
	if username != "" && username[0] != '@' {
		username = "@" + username
	}
	return username
}
