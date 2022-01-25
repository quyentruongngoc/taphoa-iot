package internal

import (
	"fmt"
	"log"
	"time"
)

const MAX_TOKEN_TIME = 480 // minutes

// type TokenData struct {
// 	Token string
// 	Time  time.Time
// }

type TokenData struct {
	Time time.Time
	Role uint
	User string
}

var TokenList map[string]TokenData

func InitInternal() {
	TokenList = make(map[string]TokenData)

	go func() {
		for {
			if len(TokenList) == 0 {
				continue
			}

			tempArr := make([]string, 0)
			log.Println("Check token list", time.Now().UTC().Format(time.RFC3339))
			time.Sleep(time.Duration(10) * time.Minute)
			// search in arry to find expired token
			now := time.Now().Minute()
			for k, v := range TokenList {
				var duration = v.Time.Minute() - now
				if duration > MAX_TOKEN_TIME {
					tempArr = append(tempArr, k)
				}
			}

			// remove expired token in list
			for i := range tempArr {
				log.Println("Remove token: ", tempArr[i])
				delete(TokenList, tempArr[i])
			}
		}
	}()
}

func IsTokenExist(token string) bool {
	_, ok := TokenList[token]
	return ok
}

func AddToTokenList(token string, role uint, user string) {
	now := time.Now()
	item := TokenData{
		Time: now,
		Role: role,
		User: user,
	}
	TokenList[token] = item
}

func UpdateToken(token string) {
	item, ok := TokenList[token]
	if !ok {
		return
	}
	item.Time = time.Now()

	TokenList[token] = item
}

func GetTokenRole(token string) (uint, error) {
	item, ok := TokenList[token]
	if !ok {
		log.Printf("Failed to get role for token: %v", token)
		return 0, fmt.Errorf("Failed to get role for token: %v", token)
	}

	return item.Role, nil
}

func GetTokenUser(token string) (string, error) {
	item, ok := TokenList[token]
	if !ok {
		log.Printf("Failed to get user for token: %v", token)
		return "", fmt.Errorf("Failed to get user for token: %v", token)
	}

	return item.User, nil
}
