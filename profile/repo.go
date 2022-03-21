package profile

import (
	"context"
	"fmt"
	"sort"

	"github.com/google/go-github/v43/github"
)

func (p *Profile) getOwnedRepos(ctx context.Context) ([]*github.Repository, error) {
	var all []*github.Repository
	page := 1
	for {
		ownedRepos, resp, err := p.gh.Repositories.List(ctx, "", &github.RepositoryListOptions{
			Visibility:  "all",
			Affiliation: "owner",
			Sort:        "updated",
			ListOptions: github.ListOptions{
				PerPage: 100,
				Page:    page,
			},
		})
		if err != nil {
			return nil, fmt.Errorf("github api error while retrieving owned repos: %w", err)
		}
		all = append(all, ownedRepos...)
		if resp.NextPage == 0 {
			break
		}
		page++
	}
	sort.Slice(all, func(i, j int) bool {
		return *all[i].StargazersCount > *all[j].StargazersCount
	})
	return all, nil
}
