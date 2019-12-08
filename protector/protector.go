package main

import (
	"math/rand"
	"strconv"
)

func get_session_key() string {
	var result = ""
	for i := 0; i < 10; i++ {
		result += strconv.Itoa(rand.Intn(9) + 1)
	}
	return result
}

func get_hash_str() string {
	var li = ""
	for i := 0; i < 5; i++ {
		li += strconv.Itoa(rand.Intn(6) + 1)
	}
	return li
}

type SessionProtector struct {
	hash string
}

func NewSessionProtector(hash_str string) *SessionProtector {
	sp := new(SessionProtector)
	sp.hash = hash_str
	return sp
}

func (sp *SessionProtector) calc_hash(session_key string, vall int) (string, error) {
	var result = ""

	switch vall {
	case 1:
		number, err := strconv.Atoi(session_key[0:5])
		if err != nil {
			return result, err
		}
		result += "00" + strconv.Itoa(number%97)
		return result[len(result)-2:], nil
	case 2:

	}

	return result, nil

}

func main() {
	var str = "00" + "12"
	print(str[len(str)-2:])
}
