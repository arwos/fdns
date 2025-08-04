/*
 *  Copyright (c) 2020-2025 Mikhail Knyazhev <markus621@yandex.com>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

package resolver

import (
	"context"
	"fmt"
	"slices"
	"time"

	"github.com/miekg/dns"
	"github.com/osspkg/fdns/app/vars"
	"go.osspkg.com/goppy/v2/xdns"
	"go.osspkg.com/logx"
	"go.osspkg.com/random"

	"github.com/osspkg/fdns/app/cache"
	"github.com/osspkg/fdns/app/database"
	"github.com/osspkg/fdns/app/dnscli"
	"github.com/osspkg/fdns/app/rules"
)

type Resolver struct {
	cli   *dnscli.Client
	cache *cache.Records
	rules *rules.Rules
	repo  *database.RepoModels
}

func NewResolver(dc *dnscli.Client, ds *xdns.Server, cr *cache.Records, rr *rules.Rules, rm *database.RepoModels) *Resolver {
	res := &Resolver{
		cli:   dc,
		cache: cr,
		rules: rr,
		repo:  rm,
	}

	ds.HandleFunc(res)

	return res
}

func (v *Resolver) Exchange(question dns.Question) ([]dns.RR, error) {
	if c, ok := v.cache.Get(question.Qtype, question.Name); ok {
		random.Shuffle(c.Values)
		return ParseRR(c.Values...)
	}

	if v.rules.IsBlocked(question.Name) {
		return nil, fmt.Errorf("%s is blocked", question.Name)
	}

	if values := v.rules.GetUserRules(question.Qtype, question.Name); len(values) > 0 {
		random.Shuffle(values)
		return ParseRR(values...)
	}

	response, err := v.cli.Exchange(question)
	if err != nil {
		return nil, fmt.Errorf("exchange %s error: %w", question.Name, err)
	}

	if len(response) > 0 {
		values := make([]string, 0, len(response))
		for _, rr := range response {
			rr.Header().Ttl = vars.DefaultTTl
			values = append(values, rr.String())
		}

		slices.Sort(values)
		v.cache.Set(question.Qtype, question.Name, vars.DefaultTTl, values...)
		go v.saveHistory(question.Name, question.Qtype, values...)

		random.Shuffle(values)
		return ParseRR(values...)
	}

	values, err := v.loadHistory(question.Name, question.Qtype)
	if err != nil {
		return nil, fmt.Errorf("form history %s error: %w", question.Name, err)
	}

	v.cache.Set(question.Qtype, question.Name, vars.DefaultTTl, values...)

	random.Shuffle(values)
	return ParseRR(values...)
}

func (v *Resolver) loadHistory(name string, qtype uint16) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	list, err := v.repo.ReadHistoryByDomain(ctx, []string{name})
	if err != nil {
		return nil, err
	}

	response := make([]string, 0, len(list))
	for _, item := range list {
		if item.QType != qtype {
			continue
		}
		response = append(response, item.Value...)
	}

	return response, nil
}

func (v *Resolver) saveHistory(name string, qtype uint16, values ...string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	rec := &database.History{Domain: name, QType: qtype, Value: values}
	err := v.repo.CreateHistory(ctx,
		[]*database.History{rec},
		database.ConflictUpdate([]string{"domain", "qtype"}, []string{"value"}),
	)
	if err != nil {
		logx.Error("History save", "err", err, "name", name, "qtype", qtype, "values", values)
	}
}
