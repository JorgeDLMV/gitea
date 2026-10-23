// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package integration

import (
	"fmt"
	"net/http"
	"net/url"
	"testing"

	auth_model "gitea.dev/models/auth"
	"gitea.dev/models/unittest"
	user_model "gitea.dev/models/user"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestActionsCacheDeletion(t *testing.T) {
	onGiteaRun(t, func(t *testing.T, _ *url.URL) {
		user2 := unittest.AssertExistsAndLoadBean(t, &user_model.User{ID: 2})
		session := loginUser(t, user2.Name)
		token := getTokenForLoggedInUser(t, session, auth_model.AccessTokenScopeWriteRepository, auth_model.AccessTokenScopeWriteUser)
		repo := createActionsTestRepo(t, token, "actions-cache-deletion", false)
		runner := newMockRunner()
		runner.registerAsRepoRunner(t, user2.Name, repo.Name, "mock-runner", []string{"ubuntu-latest"}, false)

		apiURL := fmt.Sprintf("/api/v1/repos/%s/%s/actions/caches", user2.Name, repo.Name)
		MakeRequest(t, NewRequest(t, "DELETE", apiURL).AddTokenAuth(token), http.StatusUnprocessableEntity)
		MakeRequest(t, NewRequest(t, "DELETE", apiURL+"?key=deps-linux&ref=refs/heads/main").AddTokenAuth(token), http.StatusNoContent)
		readToken := getTokenForLoggedInUser(t, loginUser(t, "user4"), auth_model.AccessTokenScopeWriteRepository)
		MakeRequest(t, NewRequest(t, "DELETE", apiURL+"?key=deps-linux").AddTokenAuth(readToken), http.StatusForbidden)

		req := NewRequestWithValues(t, "POST", fmt.Sprintf("/%s/%s/settings/actions/general/caches/delete", user2.Name, repo.Name), map[string]string{"key": "", "ref": ""})
		session.MakeRequest(t, req, http.StatusOK)

		workflowPath := ".gitea/workflows/cache.yml"
		createWorkflowFile(t, token, user2.Name, repo.Name, workflowPath, getWorkflowCreateFileOptions(user2, repo.DefaultBranch, "create "+workflowPath, `name: cache
on: push
jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - run: echo
`))
		deletions := runner.fetchTask(t).Context.GetFields()["gitea_cache_deletions"].GetListValue().GetValues()
		require.Len(t, deletions, 2)
		byKey := deletions[0].GetStructValue().GetFields()
		assert.Equal(t, "deps-linux", byKey["key"].GetStringValue())
		assert.Equal(t, "refs/heads/main", byKey["ref"].GetStringValue())
		assert.Positive(t, byKey["created"].GetNumberValue())
		all := deletions[1].GetStructValue().GetFields()
		assert.Empty(t, all["key"].GetStringValue())
		assert.Empty(t, all["ref"].GetStringValue())
	})
}
