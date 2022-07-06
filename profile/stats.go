package profile

import (
	"context"
	"sort"
	"strings"

	"github.com/google/go-github/v43/github"
)

type Stats struct {
	User              *github.User
	OwnedRepositories []*github.Repository
	TotalStars        int
	TotalFollowers    int
}

func (p *Profile) Stats(ctx context.Context) (*Stats, error) {

	username := strings.Split(p.config.Context.Repository, "/")[0]

	if p.stats != nil {
		return p.stats, nil
	}

	var stats Stats

	user, _, err := p.gh.Users.Get(ctx, username)
	if err != nil {
		return nil, err
	}
	stats.User = user
	stats.TotalFollowers = user.GetFollowers()

	ownedRepositores, err := p.getOwnedRepos(ctx)
	if err != nil {
		return nil, err
	}
	stats.OwnedRepositories = ownedRepositores

	for _, repo := range ownedRepositores {
		stats.TotalStars += repo.GetStargazersCount()
	}

	if p.config.ExtraRepo != "" {
		owner, repo, found := strings.Cut(p.config.ExtraRepo, "/")
		if found {
			if repository, _, err := p.gh.Repositories.Get(ctx, owner, repo); err == nil {
				if p.config.ExtraRepoDescription != "" {
					repository.Description = &p.config.ExtraRepoDescription
				}
				stats.OwnedRepositories = append(stats.OwnedRepositories, repository)
			}
		}
	}

	sort.Slice(stats.OwnedRepositories, func(i, j int) bool {
		return stats.OwnedRepositories[i].GetStargazersCount() < stats.OwnedRepositories[j].GetStargazersCount()
	})

	p.stats = &stats
	return p.stats, nil
}
