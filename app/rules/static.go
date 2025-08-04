/*
 *  Copyright (c) 2020-2025 Mikhail Knyazhev <markus621@yandex.com>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

package rules

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.osspkg.com/goppy/v2/orm"
	"go.osspkg.com/ioutils/cache"
	"go.osspkg.com/logx"
	"go.osspkg.com/routine"

	"github.com/osspkg/fdns/app/database"
)

type StaticRules struct {
	repo *database.RepoModels
	data cache.Cache[string, []string]
}

func NewStaticRules(repo *database.RepoModels) *StaticRules {
	return &StaticRules{
		repo: repo,
		data: cache.New[string, []string](),
	}
}

func (v *StaticRules) Up(ctx context.Context) error {
	go routine.Interval(ctx, time.Hour, func(ctx context.Context) {
		routine.Retry(10, time.Second, func() error { //nolint:errcheck
			if err := v.ForceUpdate(ctx); err != nil {
				if !errors.Is(err, orm.ErrTagNotFound) {
					logx.Error("StaticRules update", "err", err)
				}
				return err
			}
			return nil
		})
	})

	return nil
}

func (v *StaticRules) Down() error {
	return nil
}

func (v *StaticRules) key(qtype uint16, domain string) string {
	return fmt.Sprintf("%d %s", qtype, domain)
}

func (v *StaticRules) ForceUpdate(ctx context.Context) error {
	list, err := v.repo.ReadStaticRuleByDisabled(ctx, []bool{false})
	if err != nil {
		return fmt.Errorf("read static rules: %w", err)
	}

	result := make(map[string][]string, len(list))
	for _, item := range list {
		result[v.key(item.QType, item.Rule)] = item.Value
	}

	v.data.Replace(result)

	return nil
}

func (v *StaticRules) Convert(qtype uint16, domain string) ([]string, bool) {
	return v.data.Get(v.key(qtype, domain))
}
