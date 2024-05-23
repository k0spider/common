package googleAuth

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base32"
	"fmt"
	"net/url"
	"time"
)

func GetQRCodeGoogleUrl(name, secret, title, level string, width, height uint16) (string, error) {
	data, err := url.QueryUnescape(fmt.Sprintf("otpauth://totp/%s?secret=%s&issuer=%s", name, secret, title))
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("https://api.qrserver.com/v1/create-qr-code/?data=%s&size=%dx%d&ecc=%s", data, width, height, level), nil
}

func VerifyCode(secret string, code string) bool {
	if getCode(secret, 0) == code {
		return true
	}
	return false
}

func getCode(secret string, offset int64) string {
	key, err := base32.StdEncoding.DecodeString(secret)
	if err != nil {
		return ""
	}
	// generate a one-time password using the time at 30-second intervals
	epochSeconds := time.Now().Unix() + offset
	return fmt.Sprintf("%06d", oneTimePassword(key, toBytes(epochSeconds/30)))
}

func toBytes(value int64) []byte {
	var result []byte
	mask := int64(0xFF)
	shifts := [8]uint16{56, 48, 40, 32, 24, 16, 8, 0}
	for _, shift := range shifts {
		result = append(result, byte((value>>shift)&mask))
	}
	return result
}

func toUint32(bytes []byte) uint32 {
	return (uint32(bytes[0]) << 24) + (uint32(bytes[1]) << 16) +
		(uint32(bytes[2]) << 8) + uint32(bytes[3])
}

func oneTimePassword(key []byte, value []byte) uint32 {
	// sign the value using HMAC-SHA1
	hmacSha1 := hmac.New(sha1.New, key)
	hmacSha1.Write(value)
	hash := hmacSha1.Sum(nil)

	// We're going to use a subset of the generated hash.
	// Using the last nibble (half-byte) to choose the index to start from.
	// This number is always appropriate as it's maximum decimal 15, the hash will
	// have the maximum index 19 (20 bytes of SHA1) and we need 4 bytes.
	offset := hash[len(hash)-1] & 0x0F

	// get a 32-bit (4-byte) chunk from the hash starting at offset
	hashParts := hash[offset : offset+4]

	// ignore the most significant bit as per RFC 4226
	hashParts[0] = hashParts[0] & 0x7F

	number := toUint32(hashParts)

	// size to 6 digits
	// one million is the first number with 7 digits so the remainder
	// of the division will always return < 7 digits
	pwd := number % 1000000
	return pwd
}
