package common

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"login_service/models"
	"time"
)

func TokenGenerator(user models.User) string {
	SecretKey := "TjWnZr4u7x!A%D*G-JaNdRgUkXp2s5v8y/B?E(H+MbPeShVmYq3t6w9z$C&F)J@NcRfTjWnZr4u7x!A%D*G-KaPdSgVkXp2s5v8y/B?E(H+MbQeThWmZq3t6w9z$C&F)J@NcRfUjXn2r5u7x!A%D*G-KaPdSgVkYp3s6v9y/B?E(H+MbQeThWmZq4t7w!z%C*F)J@NcRfUjXn2r5u8x/A?D(G+KaPdSgVkYp3s6v9y$B&E)H@McQeThWmZq4t7w!"
	hashString := []byte(fmt.Sprintf("%s-%s-%s", user.UserId, user.Password, time.Now().String()))
	tokenHash := md5.New()
	tokenHash.Write([]byte(SecretKey))
	return hex.EncodeToString(tokenHash.Sum(hashString))
}
