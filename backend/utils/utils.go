package bomberman

import (
	"encoding/json"
	"math/rand"
)

func GenerateUniqueID() string {
	return "client_" + RandomString(8)
}

func RandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func EncodeJSON(data map[string]interface{}) (string, error) {
	// Convert JSON to JSON string
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	return string(jsonBytes), nil
}
