/*
 *  Copyright (c) 2020-2025 Mikhail Knyazhev <markus621@yandex.com>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

package rules

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"time"

	"github.com/osspkg/fdns/app/vars"
	"go.osspkg.com/goppy/v2/orm"
	"go.osspkg.com/logx"
	"go.osspkg.com/routine"
	"go.osspkg.com/syncing"

	"github.com/osspkg/fdns/app/database"
)

type regexpRule struct {
	rule  *regexp.Regexp
	qtype uint16
	data  []string
}

func newRegexpRule(rule string, qtype uint16, data []string) (*regexpRule, error) {
	rx, err := regexp.Compile(rule)
	if err != nil {
		return nil, err
	}
	return &regexpRule{
		rule:  rx,
		qtype: qtype,
		data:  data,
	}, nil
}

func (v *regexpRule) Compile(qtype uint16, domain string) []string {
	if v.rule == nil {
		return nil
	}

	result := make([]string, 0, len(v.data))
	matches := v.rule.FindStringSubmatchIndex(domain)
	if matches == nil {
		return append(result, v.data...)
	}

	for _, s := range v.data {
		value := v.rule.ExpandString([]byte{}, s, domain, matches)
		result = append(result, vars.BuildDNSRecord(domain, qtype, string(value)))
	}

	return result
}

func (v *regexpRule) Match(name string) bool {
	if v.rule == nil {
		return false
	}
	return v.rule.MatchString(name)
}

// -------------------------------------------------------------------------------------

type RegexpRules struct {
	repo *database.RepoModels
	data map[uint16][]*regexpRule
	mux  syncing.Lock
}

func NewRegexpRules(repo *database.RepoModels) *RegexpRules {
	return &RegexpRules{
		repo: repo,
		data: make(map[uint16][]*regexpRule, 100),
		mux:  syncing.NewLock(),
	}
}

func (v *RegexpRules) Up(ctx context.Context) error {
	go routine.Interval(ctx, time.Hour, func(ctx context.Context) {
		routine.Retry(10, time.Second, func() error { //nolint:errcheck
			if err := v.ForceUpdate(ctx); err != nil {
				if !errors.Is(err, orm.ErrTagNotFound) {
					logx.Error("RegexpRules update", "err", err)
				}
				return err
			}
			return nil
		})
	})

	return nil
}

func (v *RegexpRules) Down() error {
	return nil
}

func (v *RegexpRules) ForceUpdate(ctx context.Context) error {
	list, err := v.repo.ReadRegexpRuleByDisabled(ctx, []bool{false})
	if err != nil {
		return fmt.Errorf("fail read regexp rule: %w", err)
	}

	result := make(map[uint16][]*regexpRule, len(list))

	for _, item := range list {
		rr, e := newRegexpRule(item.Rule, item.QType, item.Value)
		if e != nil {
			return fmt.Errorf("fail parse regexp rule `%s`: %w", item.Rule, e)
		}

		if _, ok := result[item.QType]; !ok {
			result[item.QType] = make([]*regexpRule, 0, 2)
		}
		result[item.QType] = append(result[item.QType], rr)
	}

	v.mux.Lock(func() {
		v.data = result
	})

	return nil
}

func (v *RegexpRules) Convert(qtype uint16, domain string) (out []string, ok bool) {
	v.mux.RLock(func() {
		if _, has := v.data[qtype]; !has {
			return
		}

		for _, datum := range v.data[qtype] {
			if datum.Match(domain) {
				out, ok = datum.Compile(qtype, domain), true
				return
			}
		}
	})

	return
}
