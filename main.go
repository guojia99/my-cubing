/*
 * Copyright (c) 2023 guojia99 All rights reserved.
 * Created: 2023/6/22 下午6:33.
 * Author:  guojia(https://github.com/guojia99)
 */

package main

import (
	"fmt"
	"os"

	"my-cubing/db"
	"my-cubing/web"
)

func main() {
	fmt.Println(os.Args)
	db.Init()
	web.NewClient().Run(os.Args[1])
}
