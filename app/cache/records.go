/*
 *  Copyright (c) 2020-2025 Mikhail Knyazhev <markus621@yandex.com>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

package cache

import (
	"fmt"
	"time"

	"go.osspkg.com/ioutils/cache"
	"go.osspkg.com/xc"
)

type Record struct {
	Values []string
	TTL    uint32
}

func (r Record) Timestamp() int64 {
	return int64(r.TTL)
}

// -----------------------------------------------------------------------------

type Records struct {
	data cache.Cache[string, *Record]
}

func NewRecords(ctx xc.Context) *Records {
	return &Records{
		data: cache.New[string, *Record](
			cache.AutoClean[string, *Record](ctx.Context(), time.Minute),
		),
	}
}

func (v *Records) key(qtype uint16, name string) string {
	return fmt.Sprintf("%d %s", qtype, name)
}

func (v *Records) Set(qtype uint16, name string, ttl uint32, values ...string) {
	v.data.Set(
		v.key(qtype, name),
		&Record{Values: values, TTL: ttl},
	)
}

func (v *Records) Has(qtype uint16, name string) bool {
	return v.data.Has(v.key(qtype, name))
}

func (v *Records) Get(qtype uint16, name string) (*Record, bool) {
	return v.data.Get(v.key(qtype, name))
}

func (v *Records) Del(qtype uint16, name string) {
	v.data.Del(v.key(qtype, name))
}
