package utils

import "golang.org/x/crypto/bcrypt"

func ComparePasswordHash(hashedPwd, pwd string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPwd), []byte(pwd))
	if err != nil {
		return false, err
	}
	return true, nil
}
