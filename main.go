/*
 * Copyright (c) 2023 guojia99 All rights reserved.
 * Created: 2023/6/22 下午6:33.
 * Author:  guojia(https://github.com/guojia99)
 */

package main

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/guojia99/my-cubing/src"
)

func main() {
	var config string
	cmd := &cobra.Command{
		Use:   "my-cubing",
		Short: "魔方赛事系统",
		RunE: func(cmd *cobra.Command, args []string) error {
			cli := &src.Client{}
			return cli.Run(config)
		},
	}

	cmd.Flags().StringVarP(&config, "config", "c", "./etc/configs.json", "配置")
	err := cmd.Execute()
	if err != nil {
		log.Println(err)
	}
}
