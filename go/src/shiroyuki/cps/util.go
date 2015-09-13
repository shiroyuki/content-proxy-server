package cps

import "regexp"

func URESearch(pattern *regexp.Regexp, content string) map[string]string {
    matches := pattern.FindStringSubmatch(content)
    result  := make(map[string]string)

    for i, name := range pattern.SubexpNames() {
        result[name] = matches[i]
    }

    return result
}
