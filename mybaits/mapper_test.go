package mybaits

import (
	"encoding/xml"
	"os"
	"testing"
)

type expected struct {
	TestChoose       string `xml:"testChoose"`
	TestChooseNative string `xml:"testChooseNative"`
}

var (
	mapper     *Mapper
	myExpected expected
)

func initTest() {
	var err error
	mapper, err = NewMapper("testdata/test.xml")
	if err != nil {
		panic(err)
	}

	expectedContent, err := os.ReadFile("testdata/expected.xml")
	if err != nil {
		panic(err)
	}

	if err = xml.Unmarshal(expectedContent, &myExpected); err != nil {
		panic(err)
	}
}

func Test_Mapper_GetStatements(t *testing.T) {
	initTest()

	stmt, err := mapper.GetStatements()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(stmt)
}

func Test_Mapper_GetRawStatement(t *testing.T) {
	initTest()
	stmt, err := mapper.GetRawStatement()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(stmt)
}
