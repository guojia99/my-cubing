/*
 *  * Copyright (c) 2023 guojia99 All rights reserved.
 *  * Created: 2023/7/10 下午4:07.
 *  * Author: guojia(https://github.com/guojia99)
 */

package src

import (
	"github.com/guojia99/my-cubing/src/api"
	"github.com/guojia99/my-cubing/src/svc"
)

type Client struct {
	svc *svc.Context
}

func (c *Client) Run(config string) (err error) {
	c.svc, err = svc.NewContext(config)
	if err != nil {
		return err
	}
	return api.NewClient(c.svc).Run()
}
