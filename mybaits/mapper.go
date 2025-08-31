package mybaits

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/beevik/etree"
)

type Mapper struct {
	root map[string]*etree.Element
}

var queryTypes = map[string]struct{}{
	"sql":    struct{}{},
	"select": struct{}{},
	"insert": struct{}{},
	"update": struct{}{},
	"delete": struct{}{},
}

func NewMapper(xmlPath string) (mapper *Mapper, err error) {
	var data []byte
	if data, err = os.ReadFile(xmlPath); err != nil {
		err = fmt.Errorf("ReadFile fail. err: %v", err)
		return
	}
	rawText := replaceCDATA(string(data))

	doc := etree.NewDocument()
	if err = doc.ReadFromString(rawText); err != nil {
		err = fmt.Errorf("ReadFromString fail. err: %v", err)
		return
	}

	mapper = &Mapper{
		root: make(map[string]*etree.Element),
	}

	root := doc.Root()

	for _, child := range root.ChildElements() {
		if _, ok := queryTypes[child.Tag]; ok {
			id := child.SelectAttrValue("id", "")
			if id != "" {
				mapper.root[id] = child
			}
		}
	}

	return
}

type MapperStmt struct {
	ID   string
	Stmt string
}

func (m *Mapper) GetRawStatement() (fullStmt string, err error) {
	var mstmts []MapperStmt
	if mstmts, err = m.GetStatements(); err != nil {
		return
	}

	var stmsList []string
	for _, mstmt := range mstmts {
		stmsList = append(stmsList, mstmt.Stmt)
	}

	stms := &Statement{
		sql: strings.Join(stmsList, ";") + ";",
	}
	return stms.formatSQL()
}

func (m *Mapper) GetStatements() (mstmts []MapperStmt, err error) {
	for id, child := range m.root {
		if child.Tag != "sql" {
			cm := &childMapper{
				child: child,
				root:  m.root,
			}
			stmt, err := cm.getStatement()
			if err != nil {
				return nil, err
			}
			mstmts = append(mstmts, MapperStmt{
				ID:   id,
				Stmt: stmt,
			})
		}

	}
	return
}

func replaceCDATA(rawText string) string {
	cdataRegex := regexp.MustCompile(`<!\[CDATA\[([\s\S]*?)\]\]>`)
	return cdataRegex.ReplaceAllStringFunc(rawText, func(match string) string {
		content := cdataRegex.FindStringSubmatch(match)[1]
		return convertCDATA(content, true)
	})
}
