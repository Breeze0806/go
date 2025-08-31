package mybaits

import (
	"testing"
)

func Test_childMapper_getStatement(t *testing.T) {
	initTest()
	tests := []struct {
		name     string
		cm       *childMapper
		wantStmt string
		wantErr  bool
	}{
		{
			name: "testBasic",
			cm: &childMapper{
				root:  mapper.root,
				child: mapper.root["testBasic"],

				properties: nil,
				native:     false,
				whenCnt:    0,
			},
			wantStmt: "select name, category, price from fruits where category = 'apple' and price < 500",
		},
		{
			name: "testParameters",
			cm: &childMapper{
				root:  mapper.root,
				child: mapper.root["testParameters"],

				properties: nil,
				native:     false,
				whenCnt:    0,
			},
			wantStmt: "select name, category, price from fruits where category = :v1 and price > :v2 and type = :v3 and content = :v4",
		},
		{
			name: "testInclude",
			cm: &childMapper{
				root:  mapper.root,
				child: mapper.root["testInclude"],

				properties: nil,
				native:     false,
				whenCnt:    0,
			},
			wantStmt: "select name, category, price from fruits where category = :v1 order by name asc",
		},
		{
			name: "testIf",
			cm: &childMapper{
				root:  mapper.root,
				child: mapper.root["testIf"],

				properties: nil,
				native:     false,
				whenCnt:    0,
			},
			wantStmt: "select name, category, price from fruits where 1 = 1 and category = :v1 and price = :v2 and name = 'Fuji'",
		},
		{
			name: "testTrim",
			cm: &childMapper{
				root:  mapper.root,
				child: mapper.root["testTrim"],

				properties: nil,
				native:     false,
				whenCnt:    0,
			},
			wantStmt: "select name, category, price from fruits where 1 = 1 and category = :v1 and price = :v2 and name = 'Fuji'",
		},
		{
			name: "testWhere",
			cm: &childMapper{
				root:  mapper.root,
				child: mapper.root["testWhere"],

				properties: nil,
				native:     false,
				whenCnt:    0,
			},
			wantStmt: "select name, category, price from fruits where category = 'apple' and price = :v1 order by name asc",
		},

		{
			name: "testSet",
			cm: &childMapper{
				root:  mapper.root,
				child: mapper.root["testSet"],

				properties: nil,
				native:     false,
				whenCnt:    0,
			},
			wantStmt: "update fruits set category = :v1, price = :v2 where name = :v3",
		},
		{
			name: "testForeach",
			cm: &childMapper{
				root:  mapper.root,
				child: mapper.root["testForeach"],

				properties: nil,
				native:     false,
				whenCnt:    0,
			},
			wantStmt: "select name, category, price from fruits where category = 'apple' and (name = :v1 or name = :v2)",
		},
		{
			name: "testBind",
			cm: &childMapper{
				root:  mapper.root,
				child: mapper.root["testBind"],

				properties: nil,
				native:     false,
				whenCnt:    0,
			},
			wantStmt: "select name, category, price from fruits where name like :v1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotStmt, err := tt.cm.getStatement()
			if (err != nil) != tt.wantErr {
				t.Errorf("childMapper.getStatement() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotStmt != tt.wantStmt {
				t.Errorf("childMapper.getStatement() = %v, want %v", gotStmt, tt.wantStmt)
			}
		})
	}
}

func Test_childMapper_getStatement_choose(t *testing.T) {
	initTest()
	tests := []struct {
		name     string
		cm       *childMapper
		wantStmt string
		wantErr  bool
	}{
		{
			name: "testChoose",
			cm: &childMapper{
				root:  mapper.root,
				child: mapper.root["testChoose"],

				properties: nil,
				native:     false,
				whenCnt:    0,
			},
			wantStmt: "select name, category, price from fruits where name = :v1 and category = :v2 and price = :v3 and category = 'apple' and category is not null",
		},
		{
			name: "testChooseNative",
			cm: &childMapper{
				root:  mapper.root,
				child: mapper.root["testChooseNative"],

				properties: nil,
				native:     true,
				whenCnt:    0,
			},
			//wantStmt: "select name, category, price from fruits where name = :v1 and price = :v2 and category = 'apple' and category is not null",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotStmt, err := tt.cm.getStatement()
			if (err != nil) != tt.wantErr {
				t.Errorf("childMapper.getStatement() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotStmt != tt.wantStmt {
				t.Errorf("childMapper.getStatement() = %v, want %v", gotStmt, tt.wantStmt)
			}
		})
	}
}

func Test_replaceFirst(t *testing.T) {
	type args struct {
		s       string
		pattern string
		replace string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "and",
			args: args{
				s:       ` AND category = 'apple' AND price = ? -- if(price != null and price !='') ORDER BY name`,
				pattern: `^[\s*]?(AND|OR|and|or)`,
			},
			want: ` category = 'apple' AND price = ? -- if(price != null and price !='') ORDER BY name`,
		},
		{
			name: "set dot",
			args: args{
				s:       "SET -- if(category != null and category !='')\n category = ?,-- if(price != null and price !='')\nprice = ?,",
				pattern: `(,)[\s]*$`,
			},
			want: "SET -- if(category != null and category !='')\n category = ?,-- if(price != null and price !='')\nprice = ?",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := replaceFirst(tt.args.s, tt.args.pattern, tt.args.replace); got != tt.want {
				t.Errorf("replaceFirst() = %v, want %v", got, tt.want)
			}
		})
	}
}
