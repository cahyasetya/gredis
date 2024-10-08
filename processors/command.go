package processors

import (
	"fmt"

	"github.com/cahyasetya/gredis/storage"
)

const (
	GET = "GET"
	SET = "SET"
	DEL = "DEL"
	FAILED = "FAILED"
)

func HandleCommand(commands []string) []byte {
	if commands[0] == GET {
		result, err := storage.Get(commands[1])
		if err != nil {
			return []byte(fmt.Sprintf("FAILED. err: %v", err))
		}
		return []byte(result)
	}else if commands[0] == SET {
		storage.Set(commands[1], commands[2])
		return []byte("SUCCESS")
	}else if commands[0] == DEL {
		err := storage.Del(commands[1])
		if err != nil {
			return []byte(fmt.Sprintf("FAILED. err: %v", err))
		}
		return []byte("SUCCESS")
	}
	return []byte(FAILED)
}
