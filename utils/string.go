package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/labstack/gommon/random"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

func GenerateToken(length uint8) string {
	if length-3 < 1 {
		return ""
	}
	return random.String(length-3) + strings.ToUpper(strconv.FormatInt(time.Now().UnixNano(), 36)[9:12])
}

func GenerateCode() string {
	return fmt.Sprintf("%06v", rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(1000000))
}

func HmacSHA256Sign(tokenKey, params []byte) (string, error) {
	mac := hmac.New(sha256.New, tokenKey)
	_, err := mac.Write(params)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(mac.Sum(nil)), nil
}

func GenerateStreamNo() string {
	return fmt.Sprintf("%s%s", MD5(GenerateToken(32))[14:], time.Now().Format(DateSecondString))
}
