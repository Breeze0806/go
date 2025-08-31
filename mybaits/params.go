package mybaits

import (
	"regexp"
	"strings"
)

var jdbcTypes = map[string][]string{
	"NUM":     {"TINYINT", "SMALLINT", "INTEGER", "BIGINT", "BIT", "DECIMAL", "DOUBLE", "FLOAT", "NUMERIC"},
	"BOOLEAN": {"BOOLEAN"},
	"DATE":    {"DATE", "TIME", "TIMESTAMP"},
	"STRING":  {"CHAR", "VARCHAR", "NCHAR", "NVARCHAR", "LONGNVARCHAR", "LONGVARCHAR"},
	"BINARY":  {"BINARY", "VARBINARY", "LONGVARBINARY", "BLOB"},
	"OTHER":   {"ARRAY", "CLOB", "CURSOR", "DATALINK", "DATETIMEOFFSET", "DISTINCT", "JAVA_OBJECT", "NCLOB", "NULL", "OTHER", "REAL", "REF", "ROWID", "SQLXML", "STRUCT", "UNDEFINED"},
}

type Param struct {
	FullName  string
	Name      string
	JdbcType  string
	JavaType  string
	MockValue string
}

func GetParams(childText, childTail string) map[string][]Param {
	params := map[string][]Param{
		"#": []Param{},
		"$": []Param{},
	}

	p := regexp.MustCompile(`\S`)
	if !p.MatchString(childText) {
		childText = ""
	}
	if !p.MatchString(childTail) {
		childTail = ""
	}
	convertString := childText + childTail

	for _, char := range []string{"#", "$"} {
		pattern := regexp.MustCompile(`\` + char + `\{.+?\}`)
		matches := pattern.FindAllString(convertString, -1)

		seen := make(map[string]bool)
		uniqueMatches := []string{}
		for _, m := range matches {
			if !seen[m] {
				seen[m] = true
				uniqueMatches = append(uniqueMatches, m)
			}
		}

		for _, match := range uniqueMatches {
			param := Param{FullName: match}
			inner := strings.TrimPrefix(strings.TrimSuffix(match, "}"), char+"{")
			parts := strings.Split(inner, ",")
			param.Name = parts[0]

			jdbcRegex := regexp.MustCompile(`\s*jdbcType\s*=\s*(\w+)`)
			if j := jdbcRegex.FindStringSubmatch(inner); len(j) > 1 {
				param.JdbcType = j[1]
			}

			javaRegex := regexp.MustCompile(`\s*javaType\s*=\s*(\w+)`)
			if j := javaRegex.FindStringSubmatch(inner); len(j) > 1 {
				param.JavaType = j[1]
			}

			param.MockValue = getMockValue(param.JdbcType)
			params[char] = append(params[char], param)
		}
	}

	return params
}

func getMockValue(jdbcType string) string {
	for _, types := range jdbcTypes {
		for _, t := range types {
			if t == jdbcType {
				return "?"
			}
		}
	}
	return "?"
}
