package auth

import (
	"fmt"
	"net/http"
	"strings"
)

func GetAPIKey(header http.Header) (string, error) {
	val := header.Get("Authorization")

	if val == "" {
		return "",fmt.Errorf("Authorization header missing")
	}

	result := strings.Split(val, " ") 
	if len(result) != 2 {
		return "",fmt.Errorf("Authorization header is broken when recieved")
	}

	if result[0] != "APIKey" {
		return "",fmt.Errorf("Authorization does'nt start with APIKey")
	}
	
	return result[1],nil
}
