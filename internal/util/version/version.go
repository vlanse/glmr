package version

import (
	"context"
	"encoding/json"
	"fmt"
	"runtime/debug"
	"strings"

	"golang.org/x/mod/semver"

	"github.com/vlanse/glmr/internal/util/request"
)

const githubAPIURL = "https://api.github.com/repos/vlanse/glmr"

var ghHeaders = map[string]string{"Accept": "application/vnd.github+json"}

type ghTag struct {
	Name   string `json:"name"`
	Commit struct {
		SHA string `json:"sha"`
	} `json:"commit"`
}

type ghCommit struct {
	Commit struct {
		Message string `json:"message"`
	} `json:"commit"`
}

func GetCurrent() string {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return ""
	}

	return strings.Split(strings.Split(info.Main.Version, "+")[0], "-")[0]
}

func CheckForUpdates(ctx context.Context) (string, string, error) {
	tags, err := fetchTags(ctx)
	if err != nil {
		return "", "", err
	}

	latest, ok := latestTag(tags)
	if !ok {
		return "", "", nil
	}

	currentVersion := GetCurrent()
	if !semver.IsValid(currentVersion) {
		return "", "", fmt.Errorf("check update version: current version %q is not a valid semantic version", currentVersion)
	}

	if semver.Compare(latest.Name, currentVersion) <= 0 {
		return "", "", nil
	}

	message, err := fetchCommitMessage(ctx, latest.Commit.SHA)
	if err != nil {
		return "", "", err
	}

	return latest.Name, message, nil
}

func fetchTags(ctx context.Context) ([]ghTag, error) {
	data, err := request.GET(ctx, githubAPIURL+"/tags", ghHeaders)
	if err != nil {
		return nil, fmt.Errorf("check update version: failed to list repo tags: %w", err)
	}

	var tags []ghTag
	if err = json.Unmarshal(data, &tags); err != nil {
		return nil, fmt.Errorf("check update version: failed to unmarshal repo tags: %w", err)
	}

	return tags, nil
}

func fetchCommitMessage(ctx context.Context, sha string) (string, error) {
	data, err := request.GET(ctx, fmt.Sprintf("%s/commits/%s", githubAPIURL, sha), ghHeaders)
	if err != nil {
		return "", fmt.Errorf("check update version: failed to get tag commit: %w", err)
	}

	var commit ghCommit
	if err = json.Unmarshal(data, &commit); err != nil {
		return "", fmt.Errorf("check update version: failed to unmarshal tag commit: %w", err)
	}

	return commit.Commit.Message, nil
}

func latestTag(tags []ghTag) (ghTag, bool) {
	var latest ghTag
	found := false
	for _, tag := range tags {
		name := normalizeVersion(tag.Name)
		if !semver.IsValid(name) {
			continue
		}
		if !found || semver.Compare(name, latest.Name) > 0 {
			latest = ghTag{Name: name, Commit: tag.Commit}
			found = true
		}
	}
	return latest, found
}

func normalizeVersion(tag string) string {
	if !strings.HasPrefix(tag, "v") {
		return "v" + tag
	}
	return tag
}
