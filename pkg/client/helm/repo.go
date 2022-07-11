package helm

import (
	"encoding/json"
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

	if err := r.Add(name, url); err != nil {
		if repo.ErrIsExistRepo(err) {
			return fmt.Sprintf("repository(%s) has already been added", name), nil
		}
		return "", err
	}

	return fmt.Sprintf("successfully added repo(%s)", name), nil
}

func (c *Client) RepoList(args map[string]interface{}) (string, error) {
	r := repo.NewRepos(c.settings)

	list, err := r.List()
	if err != nil {
		return "", err
	}

	b, err := json.Marshal(&list)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func (c *Client) RepoUpdate(args map[string]interface{}) (string, error) {
	var name string

	for _, a := range []string{"repo_name"} {
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

	r := repo.NewRepos(c.settings)

	if err := r.Update(name); err != nil {
		return "", err
	}

	return fmt.Sprintf("successfully updated repo(%s)", name), nil
}
