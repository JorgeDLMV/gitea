// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package v29

import (
	"context"

	"gitea.dev/modelmigration/base"

	"xorm.io/xorm"
)

func AddConcurrencyQueueMaxToActionRunAttemptAndJob(_ context.Context, x base.EngineMigration) error {
	type ActionRunAttempt struct {
		ConcurrencyQueueMax bool `xorm:"NOT NULL DEFAULT FALSE"`
	}
	type ActionRunJob struct {
		ConcurrencyQueueMax bool `xorm:"NOT NULL DEFAULT FALSE"`
	}
	_, err := x.SyncWithOptions(xorm.SyncOptions{
		IgnoreConstrains:  true,
		IgnoreDropIndices: true,
	}, new(ActionRunAttempt), new(ActionRunJob))
	return err
}
