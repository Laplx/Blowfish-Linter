package formatter

import (
	"regexp"
)

func FixLatex(content string, inline string, multiline string) string {
	if inline != "false" {
		re := regexp.MustCompile(`\$(.+?)\$`)
		if inline == "true" {
			content = re.ReplaceAllString(content, `\\($1\\)`)
		} else {
			delimitRe := regexp.MustCompile(`(.+?)to(.+?)`)
			if match := delimitRe.FindStringSubmatch(inline); match != nil {
				front, back := match[1], match[2]
				re = regexp.MustCompile(regexp.QuoteMeta(front) + `(.+?)` + regexp.QuoteMeta(back))
				content = re.ReplaceAllString(content, front+`${1}`+back)
			} else {
				content = re.ReplaceAllString(content, `$1`)
			}
		}
	}

	if multiline != "false" {
		reBlock := regexp.MustCompile(`\$\$(?s)(.*?)\$\$`)
		content = reBlock.ReplaceAllStringFunc(content, func(match string) string {
			inner := match[2 : len(match)-2]
			reSlash := regexp.MustCompile(`\\{2}`)
			if multiline == "true" {
				inner = reSlash.ReplaceAllString(inner, `\\\`)
			} else {
				inner = reSlash.ReplaceAllString(inner, multiline)
			}
			return `$$` + inner + `$$`
		})
	}
	return content
}
