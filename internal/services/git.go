package services

import (
	"fmt"
	"github.com/go-git/go-git/v5/plumbing/object"
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

func GetLatestTagSha(repoPath string) (*object.Tag, error) {
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		return nil, fmt.Errorf("no Git repository found in '%s': %w", repoPath, err)
	}

	tags, err := repo.Tags()
	if err != nil {
		return nil, fmt.Errorf("error fetching tags: %w", err)
	}

	var latestVersion *semver.Version
	var latestSemanticTag *object.Tag

	err = tags.ForEach(func(ref *plumbing.Reference) error {
		tagObj, err := repo.TagObject(ref.Hash())
		if err != nil {
			return nil
		}

		tagName := ref.Name().Short()
		verStr := strings.TrimPrefix(tagName, "v")

		v, err := semver.NewVersion(verStr)
		if err != nil {
			return nil
		}

		if latestVersion == nil || v.GreaterThan(latestVersion) {
			latestVersion = v
			latestSemanticTag = tagObj
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("error iterating over tags: %w", err)
	}

	if latestSemanticTag == nil {
		return nil, fmt.Errorf("no valid annotated semver tags found")
	}

	fmt.Println("Points to commit SHA:", latestSemanticTag.Target.String())

	return latestSemanticTag, nil
}
