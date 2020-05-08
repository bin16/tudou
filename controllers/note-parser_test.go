package controllers

import (
	"testing"
)

func TestNoteParser(t *testing.T) {
	s := `
	Hello World! * #Zoom **biu** ~~Helo ~~
	--------------------- ::ap ::alsoi Github (https://www.github.com)
	//#endregion zirj89ye8912iolaopdkdkoakodko sdjasdn
	 && **Zoome!** p[p[dnndnqdkq]] @7:00pm https://www.google.com/
	 sadu90-1 -p2po 121 [www](https://mail.google.com)  90ui u12 e
	`

	items := noteParser(s)
	result := ""
	for _, cell := range items {
		// if cell[0] != "text" {
		// 	log.Println(cell[0], cell[1])
		// }
		result += cell[1]
	}
	if result != s {
		t.Errorf("noteParser check failed: \n%s", result)
	}
}

func TestPickMDlink(t *testing.T) {
	s := `Lets's test [Github 🐱](https://www.github.com/? ) again`
	mdlink := pickMDlink(s)
	if !mdlink.Exists() {
		t.Errorf("%s should include MDlink", s)
	}
	if mdlink.Type != TypeMDlink {
		t.Errorf("%s type check failed", mdlink.Type)
	}
	if mdlink.Matched != "[Github 🐱](https://www.github.com/? )" {
		t.Errorf("%s matched check failed", mdlink.Matched)
	}
	if mdlink.Left != "Lets's test " {
		t.Errorf("%s left check failed", mdlink.Left)
	}
	if mdlink.Right != " again" {
		t.Errorf("%s right check failed", mdlink.Right)
	}
}

func TestPickHTTP(t *testing.T) {
	s := `Some words, Github (https://www.github.com/some-test_path?a=1&b=ZZZ) `
	httpURL := pickHTTP(s)
	if !httpURL.Exists() {
		t.Errorf("%s should include http://", s)
	}
	if httpURL.Type != TypeHTTP {
		t.Errorf("%s type check failed", httpURL.Type)
	}
	if httpURL.Matched != "https://www.github.com/some-test_path?a=1&b=ZZZ" {
		t.Errorf("%s matched check failed", httpURL.Matched)
	}
	if httpURL.Left != "Some words, Github (" {
		t.Errorf("%s left check failed", httpURL.Left)
	}
	if httpURL.Right != ") " {
		t.Errorf("%s right check failed", httpURL.Right)
	}
}

func TestPickHash(t *testing.T) {
	s := `0d9hj190-jh- 1hj -0-01 #109uj3-29 1092u3 jc0uf09124 093#$u903u`
	hashTag := pickHash(s)
	if !hashTag.Exists() {
		t.Errorf("%s should include #hash", s)
	}
	if hashTag.Type != TypeHash {
		t.Errorf("%s type check failed", hashTag.Type)
	}
	if hashTag.Matched != "#109uj3" {
		t.Errorf("%s matched check failed", hashTag.Matched)
	}
	if hashTag.Left != "0d9hj190-jh- 1hj -0-01 " {
		t.Errorf("%s left check failed", hashTag.Left)
	}
	if hashTag.Right != "-29 1092u3 jc0uf09124 093#$u903u" {
		t.Errorf("%s right check failed", hashTag.Right)
	}
}

func TestPickTime(t *testing.T) {
	s := `
	Hello @07:00am zikamkd `
	time := pickTime(s)
	if !time.Exists() {
		t.Errorf("%s should include @time", s)
	}
	if time.Type != TypeTime {
		t.Errorf("%s type check failed", time.Type)
	}
	if time.Matched != "@07:00am" {
		t.Errorf("%s matched check failed", time.Matched)
	}
	if time.Left != "\n\tHello " {
		t.Errorf("%s left check failed", time.Left)
	}
	if time.Right != " zikamkd " {
		t.Errorf("%s right check failed", time.Right)
	}
}

func TestPickBold(t *testing.T) {
	s := `hello World** Zoemjkaiji** `
	bolds := pickBold(s)
	if !bolds.Exists() {
		t.Errorf("%s should include **bold**", s)
	}
	if bolds.Type != TypeBold {
		t.Errorf("%s type check failed", bolds.Type)
	}
	if bolds.Matched != "** Zoemjkaiji**" {
		t.Errorf("%s matched check failed", bolds.Matched)
	}
	if bolds.Left != "hello World" {
		t.Errorf("%s left check failed", bolds.Left)
	}
	if bolds.Right != " " {
		t.Errorf("%s right check failed", bolds.Right)
	}
}

func TestPickMark(t *testing.T) {
	s := `0e812:hne8h09::9u01e2u0::`
	mark := pickMark(s)
	if !mark.Exists() {
		t.Errorf("%s should include ::mark::", s)
	}
	if mark.Type != TypeMark {
		t.Errorf("%s type check failed", mark.Type)
	}
	if mark.Matched != "::9u01e2u0::" {
		t.Errorf("%s matched check failed", mark.Matched)
	}
	if mark.Left != "0e812:hne8h09" {
		t.Errorf("%s left check failed", mark.Left)
	}
	if mark.Right != "" {
		t.Errorf("%s right check failed", mark.Right)
	}
}

func TestPickDel(t *testing.T) {
	s := `11e9iie1-0~ie1-09udfh981h~~ 39 (& ~~GT(G(G80891uhjr0819
zcaj09ajd 0k-`
	del := pickDel(s)
	if !del.Exists() {
		t.Errorf("%s should include ::del::", s)
	}
	if del.Type != TypeDel {
		t.Errorf("%s type check failed", del.Type)
	}
	if del.Matched != "~~ 39 (& ~~" {
		t.Errorf("%s matched check failed", del.Matched)
	}
	if del.Left != "11e9iie1-0~ie1-09udfh981h" {
		t.Errorf("%s left check failed", del.Left)
	}
	if del.Right != "GT(G(G80891uhjr0819\nzcaj09ajd 0k-" {
		t.Errorf("%s right check failed", del.Right)
	}
}
