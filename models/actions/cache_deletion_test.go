// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package actions

import (
	"strings"
	"testing"
	"time"

	"gitea.dev/models/db"
	"gitea.dev/models/unittest"
	"gitea.dev/modules/timeutil"
	"gitea.dev/modules/util"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCacheDeletions(t *testing.T) {
	require.NoError(t, unittest.PrepareTestDatabase())

	expired := &ActionCacheDeletion{RepoID: 1, Key: "old", CreatedUnix: timeutil.TimeStampNow().AddDuration(-cacheDeletionLifetime - time.Hour)}
	_, err := db.GetEngine(t.Context()).NoAutoTime().Insert(expired)
	require.NoError(t, err)

	require.NoError(t, CreateCacheDeletion(t.Context(), 1, "deps", "refs/heads/main"))
	require.NoError(t, CreateCacheDeletion(t.Context(), 2, "", ""))
	assert.ErrorIs(t, CreateCacheDeletion(t.Context(), 1, strings.Repeat("k", CacheKeyMaxLength+1), ""), util.ErrInvalidArgument)

	deletions, err := FindCacheDeletions(t.Context(), 1)
	require.NoError(t, err)
	require.Len(t, deletions, 1)
	assert.Equal(t, "deps", deletions[0].Key)
	assert.Equal(t, "refs/heads/main", deletions[0].Ref)
	unittest.AssertNotExistsBean(t, &ActionCacheDeletion{ID: expired.ID})
}
