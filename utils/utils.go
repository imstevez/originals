package utils

import (
	"crypto/md5"
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/pborman/uuid"
)

func RegMatch(str string, pattern string) bool {
	reg := regexp.MustCompile(pattern)
	return reg.MatchString(str)
}

func Password(password string) (pwd string, salt string) {
	salt = uuid.New()
	pwd = Hash(password, salt)
	return
}

func Hash(password, salt string) string {
	return Md5(strings.NewReader(Md5(strings.NewReader(password)) + salt))
}

// Md5 function
func Md5(reader io.Reader) string {
	hash := md5.New()
	_, _ = io.Copy(hash, reader)
	return fmt.Sprintf("%X", hash.Sum(nil))
}
