package alertmanager

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/prometheus/alertmanager/api/v2/models"
)

func (c *Client) GetSilence(apiPath string, params map[string]interface{}) (string, error) {
	var silenceId string

	if found, err := FindCastFromMap(params, "silence_id", &silenceId); found && err != nil {
		return "", err
	} else if !found {
		return "", err
	} else if silenceId == "" {
		return "", fmt.Errorf("silence_id is empty")
	}

	ctx, cancel := context.WithTimeout(context.Background(), defaultQueryTimeout)
	defer cancel()

	body, err := c.client.Get(apiPath + "/silence/" + silenceId).Do(ctx).Raw()
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func (c *Client) GetSilences(apiPath string, params map[string]interface{}) (string, error) {
	var filterInfs []interface{}

	// optional
	if found, err := FindCastFromMap(params, "filter", &filterInfs); found && err != nil {
		return "", err
	}

	var filter []string
	for _, inf := range filterInfs {
		matcher, ok := inf.(string)
		if !ok {
			return "", fmt.Errorf("type of '%s' must be string, not %T", "filter's item", inf)
		}
		filter = append(filter, matcher)
	}

	ctx, cancel := context.WithTimeout(context.Background(), defaultQueryTimeout)
	defer cancel()

	body, err := c.client.Get(apiPath+"/silences").SetParam("filter", filter...).Do(ctx).Raw()
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func (c *Client) CreateSilences(apiPath string, params map[string]interface{}) (string, error) {
	var author string
	var start string
	var end string
	var comment string
	var matchersInfs []interface{}

	// required
	if found, err := FindCastFromMap(params, "author", &author); found && err != nil {
		return "", err
	} else if !found {
		return "", err
	} else if author == "" {
		return "", fmt.Errorf("author is empty")
	}

	if found, err := FindCastFromMap(params, "start", &start); found && err != nil {
		return "", err
	} else if !found {
		return "", err
	} else if start == "" {
		return "", fmt.Errorf("start is empty")
	}

	if found, err := FindCastFromMap(params, "end", &end); found && err != nil {
		return "", err
	} else if !found {
		return "", err
	} else if end == "" {
		return "", fmt.Errorf("end is empty")
	}

	if found, err := FindCastFromMap(params, "matchers", &matchersInfs); found && err != nil {
		return "", err
	} else if !found {
		return "", err
	} else if len(matchersInfs) == 0 {
		return "", fmt.Errorf("matchers is empty")
	}

	// optional
	if found, err := FindCastFromMap(params, "comment", &comment); found && err != nil {
		return "", err
	}

	var matchers []string
	for _, inf := range matchersInfs {
		matcher, ok := inf.(string)
		if !ok {
			return "", fmt.Errorf("type of '%s' must be string, not %T", "matchers's item", inf)
		}
		matchers = append(matchers, matcher)
	}

	modelsMatchers, err := ConvertMathcersToModels(matchers)
	if err != nil {
		return "", err
	}

	startsAt, err := strfmt.ParseDateTime(start)
	if err != nil {
		return "", err
	}

	endsAt, err := strfmt.ParseDateTime(end)
	if err != nil {
		return "", err
	}

	if time.Time(startsAt).After(time.Time(endsAt)) {
		return "", fmt.Errorf("start must be before end")
	}

	silence := &models.Silence{
		Matchers:  modelsMatchers,
		CreatedBy: &author,
		Comment:   &comment,
		StartsAt:  &startsAt,
		EndsAt:    &endsAt,
	}

	if err := silence.Validate(strfmt.Default); err != nil {
		return "", err
	}

	b, err := json.Marshal(silence)
	if err != nil {
		return "", err
	}

	ctx, cancel := context.WithTimeout(context.Background(), defaultQueryTimeout)
	defer cancel()

	body, err := c.client.Post(apiPath+"/silences").SetBody("application/json", b).Do(ctx).Raw()
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func (c *Client) DeleteSilence(apiPath string, params map[string]interface{}) (string, error) {
	var silenceId string

	// required
	if found, err := FindCastFromMap(params, "silence_id", &silenceId); found && err != nil {
		return "", err
	} else if !found {
		return "", err
	} else if silenceId == "" {
		return "", fmt.Errorf("silence_id is empty")
	}

	ctx, cancel := context.WithTimeout(context.Background(), defaultQueryTimeout)
	defer cancel()

	body, err := c.client.Delete(apiPath + "/silence/" + silenceId).Do(ctx).Raw()
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func (c *Client) UpdateSilence(apiPath string, params map[string]interface{}) (string, error) {
	var silenceId string
	var start string
	var end string
	var comment string

	// required
	if found, err := FindCastFromMap(params, "silence_id", &silenceId); found && err != nil {
		return "", err
	} else if !found {
		return "", err
	} else if silenceId == "" {
		return "", fmt.Errorf("silence_id is empty")
	}

	// optional
	if found, err := FindCastFromMap(params, "start", &start); found && err != nil {
		return "", err
	}

	if found, err := FindCastFromMap(params, "end", &end); found && err != nil {
		return "", err
	}

	if found, err := FindCastFromMap(params, "comment", &comment); found && err != nil {
		return "", err
	}

	startsAt, err := strfmt.ParseDateTime(start)
	if err != nil {
		return "", err
	}

	endsAt, err := strfmt.ParseDateTime(end)
	if err != nil {
		return "", err
	}

	// get silence
	str, err := c.GetSilence(apiPath, map[string]interface{}{"silence_id": silenceId})
	if err != nil {
		return "", err
	}

	existSilence := &models.GettableSilence{}
	if err := json.Unmarshal([]byte(str), existSilence); err != nil {
		return "", err
	}

	if !startsAt.Equal(strfmt.NewDateTime()) {
		existSilence.StartsAt = &startsAt
	}

	if !endsAt.Equal(strfmt.NewDateTime()) {
		existSilence.EndsAt = &endsAt
	}

	if time.Time(*existSilence.StartsAt).After(time.Time(*existSilence.EndsAt)) {
		return "", fmt.Errorf("start must be before end")
	}

	if comment != "" {
		existSilence.Comment = &comment
	}

	silence := &models.PostableSilence{
		ID:      silenceId,
		Silence: existSilence.Silence,
	}

	if err := silence.Validate(strfmt.Default); err != nil {
		return "", err
	}

	b, err := json.Marshal(silence)
	if err != nil {
		return "", err
	}

	ctx, cancel := context.WithTimeout(context.Background(), defaultQueryTimeout)
	defer cancel()

	body, err := c.client.Post(apiPath+"/silences").SetBody("application/json", b).Do(ctx).Raw()
	if err != nil {
		return "", err
	}

	return string(body), nil
}
