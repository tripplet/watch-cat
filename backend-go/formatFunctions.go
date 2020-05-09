package main

import (
	"html/template"
	"strings"
	"time"

	"github.com/hako/durafmt"
)

func formatToLower(s string) string {
	return strings.ToLower(s)
}

func formatTimespan(d time.Duration) string {
	return durafmt.Parse(d).String()
}

func formatDatetime(t time.Time) string {
	return t.In(tzLocation).Format("15:04:05 - 02.01.2006")
}

func formatAddBreakChars(s string) template.HTML {
	if strings.Contains(s, ":") {
		return template.HTML(strings.ReplaceAll(s, ":", ":<wbr/>"))
	} else if strings.Contains(s, ".") {
		return template.HTML(strings.ReplaceAll(s, ".", ".<wbr/>"))
	}

	return template.HTML(s)
}
