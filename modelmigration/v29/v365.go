// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package v29

import (
	"context"

	"gitea.dev/modelmigration/base"
	"gitea.dev/modules/timeutil"
)

func AddActionCacheDeletion(_ context.Context, x base.EngineMigration) error {
	type ActionCacheDeletion struct {
		ID          int64              `xorm:"pk autoincr"`
		RepoID      int64              `xorm:"INDEX NOT NULL"`
		Key         string             `xorm:"VARCHAR(512) NOT NULL DEFAULT ''"`
		Ref         string             `xorm:"NOT NULL DEFAULT ''"`
		CreatedUnix timeutil.TimeStamp `xorm:"created NOT NULL"`
	}

	return x.Sync(new(ActionCacheDeletion))
}
