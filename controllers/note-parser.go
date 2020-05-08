package controllers

import (
	"regexp"
	"strings"
)

// Constants
const (
	PrefixHTTP  = "http://"
	PrefixHTTPS = "https://"
	TypeBold    = "**bold**"
	TypeDel     = "~~del~~"
	TypeMark    = "::mark::"
	TypeTime    = "@time"
	TypeHash    = "#hash"
	TypeHTTP    = "http://"
	TypeText    = "text"
	TypeMDlink  = "[]()"
)

type picked struct {
	Type    string
	Left    string
	Right   string
	Matched string
}

func (p picked) Exists() bool {
	return p.Matched != ""
}

func noteParser(s string) []([2]string) {
	var cur string
	head := 0
	sl := strings.Split(s, "")
	var results []([2]string)
	for {
		ch := sl[head]
		cur += ch

		switch ch {
		case "*":
			bold := pickBold(cur)
			if bold.Exists() {
				cur = ""
				results = append(results, [2]string{TypeText, bold.Left})
				results = append(results, [2]string{bold.Type, bold.Matched})
				results = append(results, [2]string{TypeText, bold.Right})
			}
		case ":":
			mark := pickBold(cur)
			if mark.Exists() {
				cur = ""
				results = append(results, [2]string{TypeText, mark.Left})
				results = append(results, [2]string{mark.Type, mark.Matched})
				results = append(results, [2]string{TypeText, mark.Right})
			}
		case "~":
			del := pickBold(cur)
			if del.Exists() {
				cur = ""
				results = append(results, [2]string{TypeText, del.Left})
				results = append(results, [2]string{del.Type, del.Matched})
				results = append(results, [2]string{TypeText, del.Right})
			}
		case "/":
			// check https?://
			if strings.HasPrefix(PrefixHTTP, cur) || strings.HasPrefix(PrefixHTTPS, cur) {
				for {
					if sl[head+1] == " " || sl[head+1] == "\n" {
						break
					}

					cur += sl[head+1]
					head++
				}
				if httpURL := pickHTTP(cur); httpURL.Exists() {
					cur = ""
					results = append(results, [2]string{TypeText, httpURL.Left})
					results = append(results, [2]string{httpURL.Type, httpURL.Matched})
					results = append(results, [2]string{TypeText, httpURL.Right})
				}
			}
		case " ":
			if timeTag := pickTime(cur); timeTag.Exists() {
				cur = ""
				results = append(results, [2]string{TypeText, timeTag.Left})
				results = append(results, [2]string{timeTag.Type, timeTag.Matched})
				results = append(results, [2]string{TypeText, timeTag.Right})
			} else if hashTag := pickHash(cur); hashTag.Exists() {
				cur = ""
				results = append(results, [2]string{TypeText, hashTag.Left})
				results = append(results, [2]string{hashTag.Type, hashTag.Matched})
				results = append(results, [2]string{TypeText, hashTag.Right})
			}
		case ")":
			if mdlink := pickMDlink(cur); mdlink.Exists() {
				cur = ""
				results = append(results, [2]string{TypeText, mdlink.Left})
				results = append(results, [2]string{mdlink.Type, mdlink.Matched})
				results = append(results, [2]string{TypeText, mdlink.Right})
			}
		}

		head++
		if head == len(sl) {
			results = append(results, [2]string{TypeText, cur})
			break
		}
	}

	return results
}

func pickMDlink(s string) picked {
	pr := picked{}
	rx := regexp.MustCompile(`\[[^\n]+\]\([:0-9a-zA-Z\.\-_&\?=/%\s]+\)`)
	matched := rx.FindString(s)
	if matched != "" {
		parts := strings.Split(s, matched)
		pr.Type = TypeMDlink
		pr.Matched = matched
		pr.Left = parts[0]
		pr.Right = parts[1]
	}

	return pr
}

// ...**XXX**
func pickBold(s string) picked {
	pr := picked{}
	rx := regexp.MustCompile(`\*\*[^\n]+\*\*`)
	matched := rx.FindString(s)
	if matched != "" {
		parts := strings.Split(s, matched)
		pr.Type = TypeBold
		pr.Matched = matched
		pr.Left = parts[0]
		pr.Right = parts[1]
	}

	return pr
}

// ...~~XXX~~
func pickDel(s string) picked {
	pr := picked{}
	rx := regexp.MustCompile(`~~[^\n]+~~`)
	matched := rx.FindString(s)
	if matched != "" {
		parts := strings.Split(s, matched)
		pr.Type = TypeDel
		pr.Matched = matched
		pr.Left = parts[0]
		pr.Right = parts[1]
	}

	return pr
}

// ...::XXX::
func pickMark(s string) picked {
	pr := picked{}
	rx := regexp.MustCompile(`::[^\n]+::`)
	matched := rx.FindString(s)
	if matched != "" {
		parts := strings.Split(s, matched)
		pr.Type = TypeMark
		pr.Matched = matched
		pr.Left = parts[0]
		pr.Right = parts[1]
	}

	return pr
}

func pickHTTP(s string) picked {
	pr := picked{}
	rx := regexp.MustCompile(`https?://[a-zA-Z0-9\.\-_&\?\=\/]+`)
	matched := rx.FindString(s)
	if matched != "" {
		parts := strings.Split(s, matched)
		pr.Type = TypeHTTP
		pr.Matched = matched
		pr.Left = parts[0]
		pr.Right = parts[1]
	}

	return pr
}

func pickTime(s string) picked {
	pr := picked{}
	rx := regexp.MustCompile(`@[0-9:\.ampAMP]+`)
	matched := rx.FindString(s)
	if matched != "" {
		parts := strings.Split(s, matched)
		pr.Type = TypeTime
		pr.Matched = matched
		pr.Left = parts[0]
		pr.Right = parts[1]
	}

	return pr
}

func pickHash(s string) picked {
	pr := picked{}
	rx := regexp.MustCompile(`#[^\n\s\-]+`)
	matched := rx.FindString(s)
	if matched != "" {
		parts := strings.Split(s, matched)
		pr.Type = TypeHash
		pr.Matched = matched
		pr.Left = parts[0]
		pr.Right = parts[1]
	}

	return pr
}
