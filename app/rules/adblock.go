/*
 *  Copyright (c) 2020-2025 Mikhail Knyazhev <markus621@yandex.com>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

package rules

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"go.osspkg.com/algorithms/structs/bloom"
	"go.osspkg.com/goppy/v2/orm"
	"go.osspkg.com/goppy/v2/web"
	"go.osspkg.com/logx"
	"go.osspkg.com/routine"
	"go.osspkg.com/validate"

	"github.com/osspkg/fdns/app/database"
)

var (
	rex1 = regexp.MustCompile(`(?miU)^\|\|([a-z0-9-.]+)\^(\n|\r|\$)`)
	rex2 = regexp.MustCompile(`(?mUi)^\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3} ([a-z0-9-.]+)(\r|\n|$)`)
)

type AdBlock struct {
	bloom *bloom.Bloom
	cli   *web.ClientHttp
	repo  *database.RepoModels
}

func NewAdBlock(repo *database.RepoModels, cli web.ClientHttpPool) (*AdBlock, error) {
	ab := &AdBlock{
		cli:  cli.Create(),
		repo: repo,
	}
	var err error
	if ab.bloom, err = ab.createBloom(1000); err != nil {
		return nil, err
	}
	return ab, nil
}

func (v *AdBlock) createBloom(size uint64) (*bloom.Bloom, error) {
	return bloom.New(bloom.Quantity(size, 0.01))
}

func (v *AdBlock) Up(ctx context.Context) error {
	go routine.Interval(ctx, time.Hour, func(ctx context.Context) {
		routine.Retry(10, time.Second, func() error { //nolint:errcheck
			if err := v.ForceUpdate(ctx); err != nil {
				if !errors.Is(err, orm.ErrTagNotFound) {
					logx.Error("AdBlock update", "err", err)
				}
				return err
			}
			return nil
		})
	})

	//go routine.Interval(ctx, time.Hour*24, func(ctx context.Context) {
	//	time.Sleep(10 * time.Second)
	//	if err := v.UpgradeRules(ctx); err != nil {
	//		logx.Error("AdBlock update", "err", err)
	//	}
	//})

	return nil
}

func (v *AdBlock) Down() error {
	return nil
}

func (v *AdBlock) UpgradeRules(ctx context.Context) error {
	links, err := v.repo.ReadAdblockLinkByDisabled(ctx, []bool{false})
	if err != nil {
		return fmt.Errorf("fail load adblock links: %w", err)
	}

	for _, link := range links {
		count, err := func() (int, error) {

			var b []byte
			err0 := v.cli.Call(ctx, http.MethodGet, link.Link, nil, &b)
			if err0 != nil {
				return 0, err0
			}

			result := make([]string, 0, saveChunk)

			rexResult1 := rex1.FindAllSubmatch(b, -1)
			for _, rr := range rexResult1 {
				if len(rr) < 2 {
					continue
				}
				result = append(result, strings.TrimSpace(string(rr[1])))
				if len(result) >= saveChunk {
					if e := v.save(ctx, link.ID, result); e != nil {
						return 0, e
					}
					result = result[:0]
				}
			}

			rexResult2 := rex2.FindAllSubmatch(b, -1)
			for _, rr := range rexResult2 {
				if len(rr) < 2 {
					continue
				}
				result = append(result, strings.TrimSpace(string(rr[1])))
				if len(result) >= saveChunk {
					if e := v.save(ctx, link.ID, result); e != nil {
						return 0, e
					}
					result = result[:0]
				}
			}

			if len(result) > 0 {
				if e := v.save(ctx, link.ID, result); e != nil {
					return 0, e
				}
			}

			return len(rexResult1) + len(rexResult2), nil
		}()
		if err != nil {
			logx.Error("AdBlock upgrade", "uri", link.Link, "err", err)
		} else {
			logx.Info("AdBlock upgrade", "uri", link.Link, "count", count)
		}
	}

	return nil
}

const saveChunk = 10000

func (v *AdBlock) save(ctx context.Context, id int64, data []string) error {
	models := make([]*database.AdblockRule, 0, len(data))

	for _, datum := range data {
		domain, err := validate.NormalizeDomain(datum)
		if err != nil {
			logx.Warn("AdBlock validate domain", "domain", datum, "err", err)
			continue
		}
		models = append(models, &database.AdblockRule{
			LinkId:   id,
			Zone:     validate.GetDomainLevel(domain, 1),
			Rule:     domain,
			Disabled: false,
		})
	}

	return v.repo.CreateAdblockRule(ctx, models, database.ConflictIgnore())
}

func (v *AdBlock) ForceUpdate(ctx context.Context) error {
	links, err := v.repo.ReadAdblockLinkByDisabled(ctx, []bool{false})
	if err != nil {
		return fmt.Errorf("fail read adblock link: %w", err)
	}

	if len(links) == 0 {
		return nil
	}

	linksIds := make([]int64, 0, len(links))
	for _, link := range links {
		linksIds = append(linksIds, link.ID)
	}

	count, err := v.repo.CountAdblockRuleByLinkId(ctx, linksIds)
	if err != nil {
		return fmt.Errorf("fail get count adblock rule: %w", err)
	}
	if count <= 0 {
		return nil
	}

	rules, err := v.repo.ReadAdblockRuleByLinkId(ctx, linksIds)
	if err != nil {
		return fmt.Errorf("fail read adblock rule: %w", err)
	}

	var bf *bloom.Bloom
	if bf, err = v.createBloom(uint64(count)); err != nil {
		return err
	}

	cnt := 0
	for _, rule := range rules {
		if rule.Disabled {
			continue
		}
		cnt++
		bf.Add(rule.Rule)
	}

	logx.Info("AdBlock force update", "count", cnt)

	bf.CopyTo(v.bloom)

	return nil
}

func (v *AdBlock) Contain(name string) bool {
	levels := validate.CountDomainLevels(name)

	for i := levels; i >= 1; i-- {
		subDomain := validate.GetDomainLevel(name, i)

		if v.bloom.Contain(subDomain) {
			return true
		}
	}

	return false
}
