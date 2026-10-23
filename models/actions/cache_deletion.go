// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package actions

import (
	"context"
	"time"
	"unicode/utf8"

	"gitea.dev/models/db"
	"gitea.dev/modules/timeutil"
	"gitea.dev/modules/util"
)

// ActionCacheDeletion drops the cache entries of a repository saved before it was created. Runners
// hold the cache, so it travels with every job of the repository and each runner applies it to its
// own store. An empty Key matches every key and an empty Ref every ref.
type ActionCacheDeletion struct {
	ID          int64              `xorm:"pk autoincr"`
	RepoID      int64              `xorm:"INDEX NOT NULL"`
	Key         string             `xorm:"VARCHAR(512) NOT NULL DEFAULT ''"`
	Ref         string             `xorm:"NOT NULL DEFAULT ''"`
	CreatedUnix timeutil.TimeStamp `xorm:"created NOT NULL"`
}

// CacheKeyMaxLength is the longest key actions/cache accepts.
const CacheKeyMaxLength = 512

// cacheDeletionLifetime outlasts the runners' default retention of 7 days without use, after which
// an entry no job restores is gone anyway.
const cacheDeletionLifetime = 30 * 24 * time.Hour

func init() {
	db.RegisterModel(new(ActionCacheDeletion))
}

func CreateCacheDeletion(ctx context.Context, repoID int64, key, ref string) error {
	if utf8.RuneCountInString(key) > CacheKeyMaxLength {
		return util.NewInvalidArgumentErrorf("cache key is longer than %d characters", CacheKeyMaxLength)
	}
	return db.WithTx(ctx, func(ctx context.Context) error {
		if _, err := db.GetEngine(ctx).Where("repo_id = ? AND created_unix < ?", repoID, cacheDeletionCutoff()).Delete(new(ActionCacheDeletion)); err != nil {
			return err
		}
		return db.Insert(ctx, &ActionCacheDeletion{RepoID: repoID, Key: key, Ref: ref})
	})
}

func FindCacheDeletions(ctx context.Context, repoID int64) ([]*ActionCacheDeletion, error) {
	var deletions []*ActionCacheDeletion
	return deletions, db.GetEngine(ctx).Where("repo_id = ? AND created_unix >= ?", repoID, cacheDeletionCutoff()).OrderBy("id").Find(&deletions)
}

func cacheDeletionCutoff() timeutil.TimeStamp {
	return timeutil.TimeStampNow().AddDuration(-cacheDeletionLifetime)
}
