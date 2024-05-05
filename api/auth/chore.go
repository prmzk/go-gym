package auth

import (
	"context"
	"fmt"
	"time"
)

func (authApi *authApi) choresTicker() {
	ticker := time.NewTicker(time.Hour * 1)
	go func() {
		for range ticker.C {
			if _, err := authApi.DB.PurgeExpiredTokens(context.Background()); err != nil {
				fmt.Println("Error purging expired tokens: ", err)
			}
		}
	}()
}
