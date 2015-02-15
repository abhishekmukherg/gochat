package goscs

import (
	"crypto/cipher"
	"crypto/hmac"
	"crypto/sha1"
	"strings"
)

func (s *ScsMgr) cryptData(data, iv []byte) (out []byte) {
	out = make([]byte, len(data))
	mode := cipher.NewCBCEncrypter(s.aes, iv)
	mode.CryptBlocks(out, data)
	return out
}

func (s *ScsMgr) uncryptData(data, iv []byte) (out []byte) {
	out = make([]byte, len(data))
	mode := cipher.NewCBCDecrypter(s.aes, iv)
	mode.CryptBlocks(out, data)
	return out
}

func (s *ScsMgr) authToken(eData, eAtime, eTid, eIv string) []byte {
	unmac := box(eData, eAtime, eTid, eIv)
	macDaddy := hmac.New(sha1.New, s.key)
	macDaddy.Write([]byte(unmac))
	return macDaddy.Sum(nil)
}

func box(parts ...string) string {
	return strings.Join(parts, "|")
}

func unbox(input string) []string {
	return strings.Split(input, "|")
}
