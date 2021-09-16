package git

import (
	"regexp"
	"strings"

	"github.com/pkg/errors"
)

func projectNameFromUrl(url string) (string, error) {
	submatchParts := regexp.MustCompile(`/(.+?).git`).FindStringSubmatch(strings.Trim(url, "\n"))
	if len(submatchParts) != 2 {
		return "", errors.New("unknown format of `remote get-url`")
	}

	return submatchParts[1], nil
}
