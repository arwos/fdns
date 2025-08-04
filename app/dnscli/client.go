/*
 *  Copyright (c) 2020-2025 Mikhail Knyazhev <markus621@yandex.com>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

package dnscli

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/miekg/dns"
	"go.osspkg.com/goppy/v2/orm"
	"go.osspkg.com/goppy/v2/xdns"
	"go.osspkg.com/logx"
	"go.osspkg.com/network/address"
	"go.osspkg.com/routine"

	"github.com/osspkg/fdns/app/database"
)

type Client struct {
	ns   *Rules
	cli  *xdns.Client
	repo *database.RepoModels
}

func NewClient(repo *database.RepoModels, cli *xdns.Client) *Client {
	c := &Client{
		cli:  cli,
		ns:   NewRules(),
		repo: repo,
	}

	c.cli.SetZoneResolver(c.ns)

	return c
}

func (v *Client) Up(ctx context.Context) error {
	go routine.Interval(ctx, time.Hour, func(ctx context.Context) {
		routine.Retry(10, time.Second, func() error { //nolint:errcheck
			if err := v.ForceUpdate(ctx); err != nil {
				if !errors.Is(err, orm.ErrTagNotFound) {
					logx.Error("DNS Client update dns list", "err", err)
				}
				return err
			}
			return nil
		})
	})

	return nil
}

func (v *Client) Down() error {
	return nil
}

func (v *Client) ForceUpdate(ctx context.Context) error {
	list, err := v.repo.ReadNSByDisabled(ctx, []bool{false})
	if err != nil {
		return fmt.Errorf("could not read nameservers: %w", err)
	}

	values := make(map[string][]string, len(list))

	for _, item := range list {
		values[item.Zone] = address.FixIPPort("53", item.Value...)
		logx.Info("DNS Client update", "zone", item.Zone, "ips", item.Value)
	}

	v.ns.Replace(values)

	return nil
}

func (v *Client) Exchange(question dns.Question) ([]dns.RR, error) {
	return v.cli.Exchange(question)
}
