package gchatmain

import (
	"fmt"
	"github.com/linkinpark342/gchat/users"
)

func Main() {
	userMgr := users.NewManager()
	fmt.Printf("Magic %v\n", userMgr.GetById(0))
}
