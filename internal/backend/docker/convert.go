package docker

import "strings"

func toEnv(env map[string]string) []string {
	var envs []string
	for k, v := range env {
		envs = append(envs, k+"="+v)
	}
	return envs
}

func toVolumes(paths []string) map[string]struct{} {
	set := map[string]struct{}{}
	for _, path := range paths {
		parts := strings.Split(path, ":")
		if len(parts) < 2 {
			continue
		}
		set[parts[1]] = struct{}{}
	}
	return set
}
