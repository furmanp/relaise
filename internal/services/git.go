package services

import (
	"errors"
	"fmt"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/storer"
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

var ErrNoTagsFound = errors.New("no semantic version tags found")

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
	foundTag := false

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
			foundTag = true
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("error iterating over tags: %w", err)
	}

	if !foundTag {
		return nil, ErrNoTagsFound
	}

	return latestSemanticTag, nil
}

func GetCommitMessagesSinceLastTag(repo *git.Repository, lastTag *object.Tag) ([]string, error) {
	if lastTag == nil {
		return nil, fmt.Errorf("no valid tag found. Exiting program")
	}
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

func GetAllCommitMessages(repo *git.Repository) ([]string, error) {
	head, err := repo.Head()
	if err != nil {
		return nil, fmt.Errorf("error getting HEAD: %w", err)
	}

	commitIter, err := repo.Log(&git.LogOptions{Order: git.LogOrderCommitterTime, From: head.Hash()})
	if err != nil {
		return nil, fmt.Errorf("error fetching commit log: %w", err)
	}

	var commitMessages []string
	err = commitIter.ForEach(func(c *object.Commit) error {
		commitMessages = append(commitMessages, c.Message)
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("error iterating over all commits: %w", err)
	}

	for i, j := 0, len(commitMessages)-1; i < j; i, j = i+1, j-1 {
		commitMessages[i], commitMessages[j] = commitMessages[j], commitMessages[i]
	}

	return commitMessages, nil
}
