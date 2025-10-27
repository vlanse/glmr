package version

import (
	"context"
	"fmt"
	"runtime/debug"
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/google/go-github/v76/github"
)

func GetCurrent() string {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return ""
	}

	return strings.Split(strings.Split(info.Main.Version, "+")[0], "-")[0]
}

func CheckForUpdates(ctx context.Context) (string, error) {
	cl := github.NewClient(nil)

	tags, _, err := cl.Repositories.ListTags(ctx, "vlanse", "glmr", nil)
	if err != nil {
		return "", fmt.Errorf("check update version: failed to list repo tags: %w", err)
	}

	if len(tags) == 0 {
		return "", nil
	}

	updateVersion, err := semver.NewVersion(tags[0].GetName())
	if err != nil {
		return "", fmt.Errorf("check update version: parse tag: %w", err)
	}

	currentVersion, err := semver.NewVersion(GetCurrent())
	if err != nil {
		return "", fmt.Errorf("check update version: parse current version: %w", err)
	}

	if updateVersion.GreaterThan(currentVersion) {
		return fmt.Sprintf("v%s", updateVersion.String()), nil
	}

	return "", nil
}
