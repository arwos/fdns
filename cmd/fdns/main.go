/*
 *  Copyright (c) 2020-2025 Mikhail Knyazhev <markus621@yandex.com>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

package main

import (
	"go.osspkg.com/goppy/v2"
	"go.osspkg.com/goppy/v2/orm"
	"go.osspkg.com/goppy/v2/web"
	"go.osspkg.com/goppy/v2/xdns"

	"github.com/osspkg/fdns/app/api"
	"github.com/osspkg/fdns/app/cache"
	"github.com/osspkg/fdns/app/database"
	"github.com/osspkg/fdns/app/dnscli"
	"github.com/osspkg/fdns/app/resolver"
	"github.com/osspkg/fdns/app/rules"
)

var Version = "v0.0.0-dev"

func main() {
	app := goppy.New("fDNS", Version, "filter dns")
	app.Plugins(
		web.WithServer(),
		web.WithClient(),
		orm.WithPgsqlClient(),
		orm.WithMigration(),
		orm.WithORM(),
		xdns.WithServer(),
		xdns.WithClient(),
	)
	app.Plugins(api.Plugins...)
	app.Plugins(cache.Plugins...)
	app.Plugins(rules.Plugins...)
	app.Plugins(database.Plugins...)
	app.Plugins(resolver.Plugins...)
	app.Plugins(dnscli.Plugins...)
	app.Run()
}
