package helm

import (
	"fmt"

	"github.com/NexClipper/sudory/pkg/client/helm/repo"
)

func (c *Client) RepoAdd(args map[string]interface{}) (string, error) {
	var name, url string

	for _, a := range []string{"repo_name", "repo_url"} {
		v, ok := args[a]
		if !ok {
			return "", fmt.Errorf("%s is empty", a)
		}
		_, ok = v.(string)
		if !ok {
			return "", fmt.Errorf("failed to type assertion for %s", a)
		}
	}

	name = args["repo_name"].(string)
	url = args["repo_url"].(string)

	r := repo.NewRepos(c.settings)
	
	if err := r.AddRepo(name, url); err != nil {
		return "", err
	}

	return fmt.Sprintf("successfully added repo(%s)", name), nil
}
