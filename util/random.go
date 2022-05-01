package util

import (
	"math/rand"
	"strconv"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandomInt(min, max int32) int32 {
	return min + rand.Int31n(max-min+1)
}

func RandomInt64(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}
	return sb.String()
}

func RandomDepartmentName() string {
	return RandomString(30)
}

func RandomJobTitle() string {
	return RandomString(6)
}

func RandomMinSalary() int64 {
	return RandomInt64(100, 1000)
}

func RandomMaxSalary() int64 {
	return RandomInt64(2000, 5000)
}

func RandomFirstName() string {
	return RandomString(4) + " " + RandomString(6)
}

func RandomLastName() string {
	return RandomString(5)
}

func RandomEmail() string {
	return RandomString(3) + "." + RandomLastName() + "@" + RandomString(3) + ".com"
}

func RandomPhoneNumber() string {
	return "+" + strconv.FormatInt(RandomInt64(800, 830), 10) + strconv.FormatInt(RandomInt64(11111111, 99999999), 10)
}
