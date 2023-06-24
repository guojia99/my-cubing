/*
 *  * Copyright (c) 2023 guojia99 All rights reserved.
 *  * Created: 2023/6/24 下午12:59.
 *  * Author: guojia(https://github.com/guojia99)
 */

package api

import "time"

// https://www.worldcubeassociation.org/api/v0/export/public
type T struct {
	ExportDate time.Time `json:"export_date"`
	SqlUrl     string    `json:"sql_url"`
	TsvUrl     string    `json:"tsv_url"`
}
