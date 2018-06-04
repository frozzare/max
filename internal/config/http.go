package config

import (
	"io/ioutil"

	"github.com/frozzare/go/http2"
	"github.com/frozzare/max/internal/cache"
	"github.com/frozzare/max/internal/task"
	"gopkg.in/yaml.v2"
)

func includeHTTPTask(url string, cache *cache.Cache) (*task.Task, error) {
	client := http2.NewClient(nil)

	if cache != nil {
		if buf, err := cache.Get(url); len(buf) > 0 && err == nil {
			var t *task.Task

			if err := yaml.Unmarshal(buf, &t); err != nil {
				return nil, err
			}

			return t, nil
		}
	}

	res, err := client.Get(url)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var t *task.Task
	if err := yaml.Unmarshal(body, &t); err != nil {
		return nil, err
	}

	if cache != nil {
		if err := cache.Set(url, body); err != nil {
			return nil, err
		}
	}

	return t, nil
}
