/*
 *  Copyright (c) 2020-2025 Mikhail Knyazhev <markus621@yandex.com>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

package dnscli

import (
	"go.osspkg.com/goppy/v2/plugins"
	"go.osspkg.com/goppy/v2/xdns"

	"github.com/osspkg/fdns/app/database"
)

var Plugins = plugins.Inject(
	plugins.Plugin{
		Inject: func(repo *database.RepoModels, serv *xdns.Server, client *xdns.Client) *Client {
			cli := NewClient(repo, client)
			serv.HandleFunc(client)
			return cli
		},
	},
)
