/*
 * Copyright (c) 2023 guojia99 All rights reserved.
 * Created: 2023/6/22 下午6:33.
 * Author:  guojia(https://github.com/guojia99)
 */

package main

import (
	"log"
	"time"

	"github.com/spf13/cobra"

	"github.com/guojia99/my-cubing/src"
	"github.com/guojia99/my-cubing/src/api"
	"github.com/guojia99/my-cubing/src/svc"
)

func NewAPIServerCmd() *cobra.Command {
	var config string

	cmd := &cobra.Command{
		Use:   "api",
		Short: "魔方赛事系统API",
		RunE: func(cmd *cobra.Command, args []string) error {
			cli := &src.Client{}
			return cli.Run(config)
		},
	}
	cmd.Flags().StringVarP(&config, "config", "c", "./etc/configs.json", "配置")
	return cmd
}

func NewAdminCmd() *cobra.Command {
	var config string

	var (
		name     string
		password string
	)

	cmd := &cobra.Command{
		Use:   "admin",
		Short: "添加管理员帐号密码",
		RunE: func(cmd *cobra.Command, args []string) error {
			svcCli, err := svc.NewContext(config)
			if err != nil {
				return err
			}

			if err = svcCli.DB.AutoMigrate(&api.Admin{}); err != nil {
				return err
			}
			err = svcCli.DB.Save(&api.Admin{
				UserName: name,
				Password: password,
				Timeout:  time.Now(),
			}).Error
			return err
		},
	}
	cmd.Flags().StringVarP(&config, "config", "c", "./etc/configs.json", "配置")
	cmd.Flags().StringVarP(&name, "name", "u", "admin", "用户名")
	cmd.Flags().StringVarP(&password, "password", "p", "admin", "帐号密码")
	return cmd
}

func main() {
	root := &cobra.Command{
		Use:   "my-cubing",
		Short: "魔方赛事系统",
	}

	root.AddCommand(NewAPIServerCmd(), NewAdminCmd())
	err := root.Execute()
	if err != nil {
		log.Println(err)
	}
}
