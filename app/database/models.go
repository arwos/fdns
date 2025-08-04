/*
 *  Copyright (c) 2020-2025 Mikhail Knyazhev <markus621@yandex.com>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

package database

import (
	"time"

	"github.com/google/uuid"

	"github.com/osspkg/fdns/app/types"
)

//go:generate goppy gen --type=orm:pgsql --db-read=slave --db-write=master --index=1000 --sql-dir=../../migrations

//gen:orm table=setting index=uniq:key
type Settings struct {
	ID        int64       // col=id index=pk
	Key       string      // col=key len=20
	Value     types.JSONB // col=value
	UpdatedAt time.Time   // col=updated_at auto=time.Now()
}

//gen:orm table=api_token index=uniq:token,domain
type ApiToken struct {
	ID        int64     // col=id index=pk
	Token     uuid.UUID // col=token
	QType     uint16    // col=qtype
	Domain    string    // col=domain
	Disabled  bool      // col=disabled
	UpdatedAt time.Time // col=updated_at auto=time.Now()
}

//gen:orm table=name_space index=uniq:zone
type NS struct {
	ID        int64     // col=id index=pk
	Zone      string    // col=zone len=256
	Value     []string  // col=value
	Disabled  bool      // col=disabled
	UpdatedAt time.Time // col=updated_at auto=time.Now()
}

//gen:orm table=dns_history index=uniq:domain,qtype
type History struct {
	ID        int64     // col=id index=pk
	Domain    string    // col=domain len=256
	QType     uint16    // col=qtype
	Value     []string  // col=value
	Disabled  bool      // col=disabled
	UpdatedAt time.Time // col=updated_at auto=time.Now()
}

//gen:orm table=adblock_link index=uniq:link
type AdblockLink struct {
	ID        int64     // col=id index=pk
	Link      string    // col=link len=2049
	Disabled  bool      // col=disabled
	UpdatedAt time.Time // col=updated_at auto=time.Now()
}

//gen:orm table=adblock_rule index=uniq:rule
type AdblockRule struct {
	ID        int64     // col=id index=pk
	LinkId    int64     // col=link_id index=fk:adblock_link.id
	Zone      string    // col=zone len=256
	Rule      string    // col=rule len=256
	Disabled  bool      // col=disabled
	UpdatedAt time.Time // col=updated_at auto=time.Now()
}

//gen:orm table=regexp_rule index=uniq:rule,qtype
type RegexpRule struct {
	ID        int64     // col=id index=pk
	Rule      string    // col=rule len=1024
	QType     uint16    // col=qtype
	Value     []string  // col=value
	Disabled  bool      // col=disabled
	UpdatedAt time.Time // col=updated_at auto=time.Now()
}

//gen:orm table=static_rule index=uniq:rule,qtype
type StaticRule struct {
	ID        int64     // col=id index=pk
	Rule      string    // col=rule len=256
	QType     uint16    // col=qtype
	Value     []string  // col=value
	Disabled  bool      // col=disabled
	UpdatedAt time.Time // col=updated_at auto=time.Now()
}
