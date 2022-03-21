package profile

import (
	"context"
	"strings"

	"github.com/google/go-github/v43/github"
)

type Stats struct {
	User              *github.User
	OwnedRepositories []*github.Repository
	TotalStars        int
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

	ownedRepositores, err := p.getOwnedRepos(ctx)
	if err != nil {
		return nil, err
	}
	stats.OwnedRepositories = ownedRepositores

	for _, repo := range ownedRepositores {
		stats.TotalStars += *repo.StargazersCount
	}

	p.stats = &stats
	return p.stats, nil
}
