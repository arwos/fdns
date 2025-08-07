/*
 *  Copyright (c) 2020-2025 Mikhail Knyazhev <markus621@yandex.com>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

package api

import (
	"go.osspkg.com/goppy/v2/web"

	"github.com/osspkg/fdns/app/database"
)

type Api struct {
	router web.Router
	repo   *database.RepoModels
}

func NewApi(r web.RouterPool, repo *database.RepoModels) *Api {
	return &Api{
		router: r.Main(),
		repo:   repo,
	}
}

func (v *Api) Up() error {
	api := v.router.Collection("/api")
	api.Get("/adblock/list", v.AdblockList)

	return nil
}

func (v *Api) Down() error {
	return nil
}
