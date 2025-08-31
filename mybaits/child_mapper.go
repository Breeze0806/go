package mybaits

import (
	"regexp"
	"strings"

	"github.com/beevik/etree"
)

type childMapper struct {
	root       map[string]*etree.Element
	child      *etree.Element
	properties map[string]string
	native     bool
	whenCnt    int
}

// func GetChildStatement(mybatisMapper map[string]*etree.Element, childID string, kwargs map[string]interface{}) (string, error) {
// 	child, ok := mybatisMapper[childID]
// 	if !ok {
// 		return "", &ChildNotFoundError{childID}
// 	}

// 	stmt := ConvertChildren(mybatisMapper, child, kwargs)
// 	for _, c := range child.ChildElements() {
// 		stmt += ConvertChildren(mybatisMapper, c, kwargs)
// 	}

//		return formatSQL(stmt, kwargs), nil
//	}
func (cm *childMapper) getStatement() (stmt string, err error) {
	stmtB := &strings.Builder{}
	stmtB.WriteString(cm.convert())
	for _, c := range cm.child.ChildElements() {
		ccm := &childMapper{
			root:       cm.root,
			child:      c,
			properties: cm.properties,
			native:     cm.native,
			whenCnt:    cm.whenCnt,
		}

		stmtB.WriteString(ccm.convert())
	}

	myStmt := Statement{
		sql: stmtB.String(),
	}
	return myStmt.formatSQL()
}

func (cm *childMapper) convert() string {
	switch cm.child.Tag {
	case "sql", "select", "insert", "update", "delete":
		return cm.convertParameters(true, true)
	case "include":
		return cm.convertInclude()
	case "if":
		return cm.convertIf()
	case "choose", "when", "otherwise":
		return cm.convertChooseWhenOtherwise()
	case "trim", "where", "set":
		return cm.convertTrimWhereSet()
	case "foreach":
		return cm.convertForeach()
	case "bind":
		return cm.convertBind()
	default:
		return ""
	}
}

func (cm *childMapper) convertParameters(text, tail bool) (convertString string) {
	p := regexp.MustCompile(`\S`)

	childText := cm.child.Text()
	if !p.MatchString(childText) {
		childText = ""
	}
	childTail := cm.child.Tail()
	if !p.MatchString(childTail) {
		childTail = ""
	}

	if text && tail {
		convertString = childText + childTail
	} else if text {
		convertString = childText
	} else if tail {
		convertString = childTail
	} else {
		convertString = ""
	}
	paramsMap := GetParams(childText, childTail)
	allParams := append(paramsMap["#"], paramsMap["$"]...)
	for _, p := range allParams {
		convertString = strings.ReplaceAll(convertString, p.FullName, p.MockValue)
	}

	convertString = convertCDATA(convertString, false)
	return
}

func (cm *childMapper) convertInclude() string {
	properties := make(map[string]string)
	if cm.properties != nil {
		for k, v := range cm.properties {
			properties[k] = v
		}
	}

	for _, c := range cm.child.ChildElements() {
		if c.Tag == "property" {
			name := c.SelectAttrValue("name", "")
			value := c.SelectAttrValue("value", "")
			properties[name] = value
		}
	}

	refID := cm.child.SelectAttrValue("refid", "")
	for _, char := range []string{"#", "$"} {
		pattern := regexp.MustCompile(`\` + char + `\{(.+?)\}`)
		if matches := pattern.FindStringSubmatch(refID); len(matches) > 1 {
			if val, ok := properties[matches[1]]; ok {
				refID = val
				break
			}
		}
	}
	cb := &strings.Builder{}
	includeChild, ok := cm.root[refID]
	if !ok {
		return ""
	}

	includeCM := &childMapper{
		root:       cm.root,
		child:      includeChild,
		properties: cm.properties,
		native:     cm.native,
		whenCnt:    cm.whenCnt,
	}

	cb.WriteString(includeCM.convert())

	cb.WriteString(cm.convertParameters(true, false))
	for _, c := range includeCM.child.ChildElements() {
		ccm := &childMapper{
			root:       cm.root,
			child:      c,
			properties: properties,
			native:     cm.native,
			whenCnt:    cm.whenCnt,
		}
		cb.WriteString(ccm.convert())
	}
	cb.WriteString(cm.convertParameters(false, true))

	return cb.String()
}

func (cm *childMapper) convertIf() string {

	cb := &strings.Builder{}
	// test := cm.child.SelectAttrValue("test", "")
	// cb.WriteString("-- if(")
	// cb.WriteString(test)
	// cb.WriteString(")\n")
	str := cm.convertParameters(true, false)

	cb.WriteString(str)

	for _, c := range cm.child.ChildElements() {
		ccm := &childMapper{
			root:  cm.root,
			child: c,

			properties: cm.properties,
			native:     cm.native,
			whenCnt:    cm.whenCnt,
		}
		str = ccm.convert()
		cb.WriteString(str)
	}

	cb.WriteString(cm.convertParameters(false, true))
	return cb.String()
}

func (cm *childMapper) convertChooseWhenOtherwise() string {
	native := cm.native
	whenCnt := cm.whenCnt
	cb := &strings.Builder{}

	if cm.child.Tag == "choose" {
		cb.WriteString(cm.convertParameters(true, false))
	}

	for _, c := range cm.child.ChildElements() {
		ccm := &childMapper{
			root:       cm.root,
			child:      c,
			properties: cm.properties,
			native:     cm.native,
		}
		if c.Tag == "when" {
			if !(native && whenCnt >= 1) {
				// test := c.SelectAttrValue("test", "")
				// cb.WriteString("-- if(")
				// cb.WriteString(test)
				// cb.WriteString(")\n")
				cb.WriteString(ccm.convertParameters(true, true))
				whenCnt++
				ccm.whenCnt = whenCnt
			}
		} else if c.Tag == "otherwise" {
			//cb.WriteString("-- otherwise\n")
			cb.WriteString(ccm.convertParameters(true, true))

		}
		cb.WriteString(ccm.convert())
	}

	if cm.child.Tag == "choose" {
		cb.WriteString(cm.convertParameters(false, true))
	}
	return cb.String()
}

func (cm *childMapper) convertTrimWhereSet() string {
	var prefix, suffix, prefixOverrides, suffixOverrides string

	switch cm.child.Tag {
	case "trim":
		prefix = cm.child.SelectAttrValue("prefix", "")
		suffix = cm.child.SelectAttrValue("suffix", "")
		prefixOverrides = cm.child.SelectAttrValue("prefixOverrides", "")
		suffixOverrides = cm.child.SelectAttrValue("suffixOverrides", "")
	case "set":
		prefix = "SET"
		suffixOverrides = ","
	case "where":
		prefix = "WHERE"
		prefixOverrides = "AND|and|OR|or"
	default:
		return ""
	}

	cb := &strings.Builder{}
	cb.WriteString(cm.convertParameters(true, false))
	for _, c := range cm.child.ChildElements() {
		ccm := &childMapper{
			root:       cm.root,
			child:      c,
			properties: cm.properties,
			native:     cm.native,
			whenCnt:    cm.whenCnt,
		}
		cb.WriteString(ccm.convert())
	}

	convertString := cb.String()
	if prefixOverrides != "" {
		convertString = replaceFirst(convertString, `^[\s]*?(`+prefixOverrides+`)`, "")
	}

	if suffixOverrides != "" {
		pattern := `(` + suffixOverrides + `)[\s]*$`
		convertString = replaceFirst(convertString, pattern, "")
	}

	cb.Reset()
	if regexp.MustCompile(`\S`).MatchString(convertString) {
		if prefix != "" {
			cb.WriteString(prefix)
			cb.WriteString(" ")
			cb.WriteString(convertString)
		}
		if suffix != "" {
			cb.WriteString(" ")
			cb.WriteString(suffix)
		}
	}

	cb.WriteString(cm.convertParameters(false, true))
	return cb.String()
}

func (cm *childMapper) convertForeach() string {
	open := cm.child.SelectAttrValue("open", "")
	close := cm.child.SelectAttrValue("close", "")
	separator := cm.child.SelectAttrValue("separator", "")

	cb := &strings.Builder{}
	cb.WriteString(cm.convertParameters(true, false))
	for _, c := range cm.child.ChildElements() {
		ccm := &childMapper{
			root:       cm.root,
			child:      c,
			properties: cm.properties,
			native:     cm.native,
			whenCnt:    cm.whenCnt,
		}
		cb.WriteString(ccm.convert())
	}
	convertString := cb.String()
	cb.Reset()
	cb.WriteString(open)
	cb.WriteString(convertString)
	cb.WriteString(separator)
	cb.WriteString(convertString)
	cb.WriteString(close)
	cb.WriteString(cm.convertParameters(false, true))
	return cb.String()
}

func (cm *childMapper) convertBind() string {
	name := cm.child.SelectAttrValue("name", "")
	value := cm.child.SelectAttrValue("value", "")
	convertString := cm.convertParameters(false, true)
	return strings.ReplaceAll(convertString, name, value)
}

func convertCDATA(s string, reverse bool) string {
	if reverse {
		s = strings.ReplaceAll(s, "&", "&amp;")
		s = strings.ReplaceAll(s, "<", "&lt;")
		s = strings.ReplaceAll(s, ">", "&gt;")
		s = strings.ReplaceAll(s, "\"", "&quot;")
	} else {
		s = strings.ReplaceAll(s, "&amp;", "&")
		s = strings.ReplaceAll(s, "&lt;", "<")
		s = strings.ReplaceAll(s, "&gt;", ">")
		s = strings.ReplaceAll(s, "&quot;", "\"")
	}
	return s
}

func findFirst(s, pattern string) []int {
	re := regexp.MustCompile(pattern)

	return re.FindStringIndex(s)
}

func replaceFirst(s, pattern, replace string) string {
	loc := findFirst(s, pattern)
	if len(loc) == 0 {
		return s
	}
	if loc[1] < len(s) {
		return s[:loc[0]] + replace + s[loc[1]:]
	}

	return s[:loc[0]] + replace
}
