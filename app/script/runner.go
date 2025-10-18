package script

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

func RunAuto(script string, outputLines []string, envVars map[string]string, resolved map[string]string, key string) string {
	env := os.Environ()
	kvRe := regexp.MustCompile(`^([^=]+)=(.*)$`)
	for _, l := range outputLines {
		if m := kvRe.FindStringSubmatch(l); len(m) > 2 {
			env = append(env, fmt.Sprintf("%s=%s", m[1], m[2]))
		}
	}
	for k, v := range envVars {
		if _, ok := resolved[k]; !ok {
			env = append(env, fmt.Sprintf("%s=%s", k, v))
		}
	}
	cmd := exec.Command("sh", "-c", script)
	cmd.Env = env
	out, err := cmd.Output()
	if err != nil {
		fmt.Printf("Error running script for %s: %v\n", key, err)
		return ""
	}
	return strings.TrimSpace(string(out))
}
