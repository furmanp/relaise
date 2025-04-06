package services

import (
	"fmt"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/storer"
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

type ReleasePayload struct {
	TagName  string
	Messages []string
}

func GetGitRepository(repoPath string) (*git.Repository, error) {
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		return nil, fmt.Errorf("no Git repository found in '%s': %w", repoPath, err)
	}
	return repo, nil
}

func GetLatestSemanticTag(repo *git.Repository) (*object.Tag, error) {
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

	return latestSemanticTag, nil
}

func GetCommitMessagesSinceLastTag(repo *git.Repository, lastTag *object.Tag) ([]string, error) {
	commit, err := repo.CommitObject(lastTag.Target)

	if err != nil {
		return nil, fmt.Errorf("error getting commit associated with provided Tag: %w", err)
	}

	head, err := repo.Head()

	if err != nil {
		return nil, fmt.Errorf("error getting HEAD: %w", err)
	}

	headCommit, err := repo.CommitObject(head.Hash())

	if err != nil {
		return nil, fmt.Errorf("error getting HEAD commit: %w", err)

	}

	commitIter, err := repo.Log(&git.LogOptions{Order: git.LogOrderCommitterTime, From: headCommit.Hash})
	if err != nil {
		return nil, fmt.Errorf("error fetching commit log: %w", err)
	}

	var commitMessages []string
	err = commitIter.ForEach(func(c *object.Commit) error {
		if c.Hash == commit.Hash {
			return storer.ErrStop
		}
		commitMessages = append(commitMessages, c.Message)

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("error iterating over commits: %w", err)
	}

	return commitMessages, nil
}

func GetReleasePayload(repo *git.Repository) (*ReleasePayload, error) {
	latestTag, err := GetLatestSemanticTag(repo)
	if err != nil {
		return nil, fmt.Errorf("error getting latest semantic tag: %w", err)
	}

	commitMessages, err := GetCommitMessagesSinceLastTag(repo, latestTag)
	if err != nil {
		return nil, fmt.Errorf("error getting commit messages since last tag: %w", err)
	}

	return &ReleasePayload{
		TagName:  latestTag.Name,
		Messages: commitMessages,
	}, nil
}
