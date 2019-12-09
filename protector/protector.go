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

func reverse(value string) string {
	data := []rune(value)
	result := []rune{}

	for i := len(data) - 1; i >= 0; i-- {
		result = append(result, data[i])
	}
	return string(result)
}

func (sp *SessionProtector) calc_hash(session_key string, vall int) string {
	switch vall {
	case 1:
		var result = ""
		key_int, err := strconv.Atoi(session_key[0:5])
		if err != nil {
			return ""
		}
		result = "00" + strconv.Itoa(key_int%97)
		return result[len(result)-2:]
	case 2:
		return reverse(session_key)
	case 3:
		return session_key[5:] + session_key[:5]
	case 4:
		var num = 0
		var key_int int
		var err error
		for i := 1; i < len(session_key)-1; i++ {
			key_int, err = strconv.Atoi(string(session_key[i]))
			if err != nil {
				return ""
			}
			num += key_int + 41
		}
		return strconv.Itoa(num)
	case 5:
		var num = 0
		for i := 0; i < len(session_key); i++ {
			num += int((session_key[i]) ^ 43)
		}
		return strconv.Itoa(num)
	}
	key_int, err := strconv.Atoi(session_key)
	if err != nil {
		return ""
	}
	return strconv.Itoa(key_int + vall)

}

func main() {
	var hash string = "13555"
	var key string = "7242985673"
	println(key)

	var protector1 = NewSessionProtector(hash)
	println("///////////////////////////")
	println(protector1.calc_hash(key, 1))
	println(protector1.calc_hash(key, 2))
	println(protector1.calc_hash(key, 3))
	println(protector1.calc_hash(key, 4))
	println(protector1.calc_hash(key, 5))
	println(protector1.calc_hash(key, 6))

}
