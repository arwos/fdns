/*
 *  Copyright (c) 2020-2025 Mikhail Knyazhev <markus621@yandex.com>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

package api

//go:generate easyjson

import (
	"time"

	"go.osspkg.com/goppy/v2/web"
)

type (
	//easyjson:json
	AdblockListModel struct {
		Links []AdblockLinkModel `json:"links"`
		Rules []AdblockRuleModel `json:"rules"`
	}

	//easyjson:json
	AdblockLinkModel struct {
		Id        int64     `json:"id"`
		Link      string    `json:"link"`
		Disabled  bool      `json:"disabled"`
		UpdatedAt time.Time `json:"updatedAt"`
	}

	//easyjson:json
	AdblockRuleModel struct {
		Id        int64     `json:"id"`
		LinkId    int64     `json:"linkId"`
		Rule      string    `json:"rule"`
		Disabled  bool      `json:"disabled"`
		UpdatedAt time.Time `json:"updatedAt"`
	}
)

func (v *Api) AdblockList(ctx web.Context) {
	links, err := v.repo.ReadAdblockLinkAll(ctx.Context())
	if err != nil {
		ctx.Error(500, err)
		return
	}

	items, err := v.repo.ReadAdblockRuleAll(ctx.Context())
	if err != nil {
		ctx.Error(500, err)
		return
	}

	result := AdblockListModel{
		Links: make([]AdblockLinkModel, 0, len(links)),
		Rules: make([]AdblockRuleModel, 0, len(items)),
	}

	for _, link := range links {
		result.Links = append(result.Links, AdblockLinkModel{
			Id:        link.ID,
			Link:      link.Link,
			Disabled:  link.Disabled,
			UpdatedAt: link.UpdatedAt,
		})
	}

	for _, item := range items {
		result.Rules = append(result.Rules, AdblockRuleModel{
			Id:        item.ID,
			LinkId:    item.LinkId,
			Rule:      item.Rule,
			Disabled:  item.Disabled,
			UpdatedAt: item.UpdatedAt,
		})
	}

	ctx.JSON(200, result)
}
