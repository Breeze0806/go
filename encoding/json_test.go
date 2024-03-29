package encoding

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"
)

var basicJSON = `{
	"age": 100,
	"place": {
		"here": "B\\\"R"
	},
	"noop": {
		"what is a wren?": "a bird"
	},
	"happy": true,
	"immortal": false,
	"items": [1, 2, 3, {
		"tags": [1, 2, 3],
		"points": [
			[1, 2],
			[3, 4]
		]
	}, 4, 5, 6, 7],
	"arr": ["1", 2, "3", {
		"hello": "world"
	}, "4", 5],
	"vals": [1, 2, 3, {
		"sadf": "asdf"
	}],
	"name": {
		"first": "tom",
		"last": null
	},
	"created": "2014-05-16T08:28:06.989Z",
	"loggy": {
		"programmers": [{
				"firstName": "Brett",
				"lastName": "McLaughlin",
				"email": "aaaa",
				"tag": "good"
			},
			{
				"firstName": "Jason",
				"lastName": "Hunter",
				"email": "bbbb",
				"tag": "bad"
			},
			{
				"firstName": "Elliotte",
				"lastName": "Harold",
				"email": "cccc",
				"tag": "good"
			},
			{
				"firstName": 1002.3,
				"age": 101
			}
		]
	},
	"lastly": {
		"yay": "final"
	},
	"float": 1e1000
}`

var invlaidJSON = `{
	"age": 100,
	"place": {
		"here": "B\\\"R"
	},
	"noop": {
		"what is a wren?": ,"a bird"
	},
	"happy": true,
	"immortal": false,
	"items": [1, 2, 3, {
		"tags": [1, 2, 3],
		"points": [
			[1, 2],
			[3, 4]
		]
	}, 4, 5, 6, 7],
	"arr": ["1", 2, "3", {
		"hello": "world"
	}, "4", 5],
	"vals": [1, 2, 3, {
		"sadf": "asdf"
	}],
	"name": {
		"first": "tom",
		"last": null
	},
	"created": "2014-05-16T08:28:06.989Z",
	"loggy": {
		"programmers": [{
				"firstName": "Brett",
				"lastName": "McLaughlin",
				"email": "aaaa",
				"tag": "good"
			},
			{
				"firstName": "Jason",
				"lastName": "Hunter",
				"email": "bbbb",
				"tag": "bad"
			},
			{
				"firstName": "Elliotte",
				"lastName": "Harold",
				"email": "cccc",
				"tag": "good"
			},
			{
				"firstName": 1002.3,
				"age": 101
			}
		]
	},
	"lastly": {
		"yay": "final"
	}
}`

func testJSONFromString(s string) *JSON {
	return newJSONFromString(s)
}

func TestNewJsonFromString(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		want    *JSON
		wantErr bool
	}{
		{
			name: "1",
			args: args{
				s: invlaidJSON,
			},
			wantErr: true,
		},
		{
			name: "2",
			args: args{
				s: basicJSON,
			},
			want:    testJSONFromString(basicJSON),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewJSONFromString(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewJsonFromString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got == nil && tt.want == nil {
				return
			}

			if got.String() != tt.want.String() {
				t.Errorf("NewJsonFromString() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestNewJsonFromBytes(t *testing.T) {
	type args struct {
		b []byte
	}
	tests := []struct {
		name    string
		args    args
		want    *JSON
		wantErr bool
	}{
		{
			name: "1",
			args: args{
				b: []byte(invlaidJSON),
			},
			wantErr: true,
		},

		{
			name: "2",
			args: args{
				b: []byte(basicJSON),
			},
			want:    testJSONFromString(basicJSON),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewJSONFromBytes(tt.args.b)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewJsonFromBytes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got == nil && tt.want == nil {
				return
			}

			if got.String() != tt.want.String() {
				t.Errorf("NewJsonFromBytes() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestNewJsonFromFile(t *testing.T) {
	type args struct {
		filename string
	}
	tests := []struct {
		name    string
		args    args
		want    *JSON
		wantErr bool
	}{
		{
			name: "1",
			args: args{
				"1213231k3j21kl3.dadsadasda",
			},
			wantErr: true,
		},
		{
			name: "2",
			args: args{
				"test_data",
			},
			want: testJSONFromString(strings.ReplaceAll(basicJSON, "\n", "\r\n")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewJSONFromFile(tt.args.filename)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewJsonFromFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got == nil && tt.want == nil {
				return
			}

			if got.String() != tt.want.String() {
				t.Errorf("NewJsonFromFile() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestJson_GetJson(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		j       *JSON
		args    args
		want    *JSON
		wantErr bool
	}{
		{
			name: "1",
			j:    testJSONFromString(basicJSON),
			args: args{
				path: "loggy.programmers.0",
			},
			want: testJSONFromString(`{
				"firstName": "Brett",
				"lastName": "McLaughlin",
				"email": "aaaa",
				"tag": "good"
			}`),
			wantErr: false,
		},

		{
			name: "2",
			j:    testJSONFromString(basicJSON),
			args: args{
				path: "loggy.programmers.0.firstName",
			},
			wantErr: true,
		},
		{
			name: "3",
			j:    testJSONFromString(basicJSON),
			args: args{
				path: "loggy.programmers.0.1111",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.j.GetJSON(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("JSON.GetJson() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got == nil && tt.want == nil {
				return
			}

			if got.String() != tt.want.String() {
				t.Errorf("JSON.GetJson() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestJson_GetBool(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		j       *JSON
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "1",
			j:    testJSONFromString(basicJSON),
			args: args{
				path: "happy",
			},
			want: true,
		},

		{
			name: "2",
			j:    testJSONFromString(basicJSON),
			args: args{
				path: "immortal",
			},
		},
		{
			name: "3",
			j:    testJSONFromString(basicJSON),
			args: args{
				path: "loggy.programmers.0",
			},
			wantErr: true,
		},
		{
			name: "4",
			j:    testJSONFromString(basicJSON),
			args: args{
				path: "loggy.programmers.0.1111",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.j.GetBool(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("JSON.GetBool() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("JSON.GetBool() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJson_GetInt64(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		j       *JSON
		args    args
		want    int64
		wantErr bool
	}{
		{
			name: "1",
			j:    testJSONFromString(basicJSON),
			args: args{
				path: "age",
			},
			want: 100,
		},
		{
			name: "2",
			j:    testJSONFromString(basicJSON),
			args: args{
				path: "arr",
			},
			wantErr: true,
		},
		{
			name: "3",
			j:    testJSONFromString(basicJSON),
			args: args{
				path: "loggy.programmers.3.firstname",
			},
			wantErr: true,
		},
		{
			name: "4",
			j:    testJSONFromString(basicJSON),
			args: args{
				path: "loggy.programmers.3.firstName",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.j.GetInt64(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("JSON.GetInt64() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("JSON.GetInt64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJson_GetFloat64(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		j       *JSON
		args    args
		want    float64
		wantErr bool
	}{
		{
			name: "1",
			j:    testJSONFromString(basicJSON),
			args: args{
				path: "loggy.programmers.3.firstName",
			},
			want: 1002.3,
		},
		{
			name: "2",
			j:    testJSONFromString(basicJSON),
			args: args{
				path: "arr",
			},
			wantErr: true,
		},
		{
			name: "3",
			j:    testJSONFromString(basicJSON),
			args: args{
				path: "loggy.programmers.3.firstname",
			},
			wantErr: true,
		},
		// {
		// 	name: "4",
		// 	j:    testJsonFromString(basicJSON),
		// 	args: args{
		// 		path: "float",
		// 	},
		// 	wantErr: false,
		// },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.j.GetFloat64(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("JSON.GetFloat64() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("JSON.GetFloat64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJson_GetString(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		j       *JSON
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "1",
			j:    testJSONFromString(basicJSON),
			args: args{
				path: "name.first",
			},
			want: "tom",
		},
		{
			name: "2",
			j:    testJSONFromString(basicJSON),
			args: args{
				path: "arr",
			},
			wantErr: true,
		},
		{
			name: "3",
			j:    testJSONFromString(basicJSON),
			args: args{
				path: "loggy.programmers.3.firstname",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.j.GetString(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("JSON.GetString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("JSON.GetString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJson_GetArray(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		j       *JSON
		args    args
		want    []*JSON
		wantErr bool
	}{
		{
			name: "1",
			j:    testJSONFromString(basicJSON),
			args: args{
				path: "vals",
			},
			want: []*JSON{
				testJSONFromString(`1`),
				testJSONFromString(`2`),
				testJSONFromString(`3`),
				testJSONFromString(`{
		"sadf": "asdf"
	}`),
			},
		},
		{
			name: "2",
			j:    testJSONFromString(basicJSON),
			args: args{
				path: "name",
			},
			wantErr: true,
		},
		{
			name: "3",
			j:    testJSONFromString(basicJSON),
			args: args{
				path: "loggy.programmers.3.firstname",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.j.GetArray(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("JSON.GetArray() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if len(got) != len(tt.want) {
				goto ERROR
			}
			for i := range got {
				if got[i].String() != tt.want[i].String() {
					goto ERROR
				}
			}
			return
		ERROR:
			t.Errorf("JSON.GetArray() = %v, want %v", got, tt.want)

		})
	}
}

func TestJson_GetMap(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		j       *JSON
		args    args
		want    map[string]*JSON
		wantErr bool
	}{
		{
			name: "1",
			j:    testJSONFromString(basicJSON),
			args: args{
				path: "name",
			},
			want: map[string]*JSON{
				"first": testJSONFromString(`"tom"`),
				"last":  testJSONFromString(`null`),
			},
		},
		{
			name: "2",
			j:    testJSONFromString(basicJSON),
			args: args{
				path: "vals",
			},
			wantErr: true,
		},
		{
			name: "3",
			j:    testJSONFromString(basicJSON),
			args: args{
				path: "loggy.programmers.3.firstname",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.j.GetMap(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("JSON.GetMap() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}

			if len(got) != len(tt.want) {
				t.Errorf("JSON.GetMap() = %v, want %v", got, tt.want)
				return
			}

			for k := range got {
				if got[k].String() != tt.want[k].String() {
					t.Errorf("JSON.GetMap() = %v, want %v", got, tt.want)
					return
				}
			}
		})
	}
}

func TestJson_IsArray(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name string
		j    *JSON
		args args
		want bool
	}{
		{
			name: "1",
			j:    testJSONFromString(basicJSON),
			args: args{
				path: "vals",
			},
			want: true,
		},
		{
			name: "2",
			j:    testJSONFromString(basicJSON),
			args: args{
				path: "name",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.j.IsArray(tt.args.path); got != tt.want {
				t.Errorf("JSON.IsArray() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJson_IsNumber(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name string
		j    *JSON
		args args
		want bool
	}{
		{
			name: "1",
			j:    testJSONFromString(basicJSON),
			args: args{
				path: "loggy.programmers.3.firstName",
			},
			want: true,
		},
		{
			name: "2",
			j:    testJSONFromString(basicJSON),
			args: args{
				path: "arr",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.j.IsNumber(tt.args.path); got != tt.want {
				t.Errorf("JSON.IsNumber() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJson_IsJson(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name string
		j    *JSON
		args args
		want bool
	}{
		{
			name: "1",
			j:    testJSONFromString(basicJSON),
			args: args{
				path: "loggy.programmers.0",
			},
			want: true,
		},

		{
			name: "2",
			j:    testJSONFromString(basicJSON),
			args: args{
				path: "loggy.programmers.0.firstName",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.j.IsJSON(tt.args.path); got != tt.want {
				t.Errorf("JSON.IsJson() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJson_IsBool(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name string
		j    *JSON
		args args
		want bool
	}{
		{
			name: "1",
			j:    testJSONFromString(basicJSON),
			args: args{
				path: "happy",
			},
			want: true,
		},

		{
			name: "2",
			j:    testJSONFromString(basicJSON),
			args: args{
				path: "immortal",
			},
			want: true,
		},
		{
			name: "3",
			j:    testJSONFromString(basicJSON),
			args: args{
				path: "loggy.programmers.0",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.j.IsBool(tt.args.path); got != tt.want {
				t.Errorf("JSON.IsBool() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJson_IsString(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name string
		j    *JSON
		args args
		want bool
	}{
		{
			name: "1",
			j:    testJSONFromString(basicJSON),
			args: args{
				path: "name.first",
			},
			want: true,
		},
		{
			name: "2",
			j:    testJSONFromString(basicJSON),
			args: args{
				path: "arr",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.j.IsString(tt.args.path); got != tt.want {
				t.Errorf("JSON.IsString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJson_IsNull(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name string
		j    *JSON
		args args
		want bool
	}{
		{
			name: "1",
			j:    testJSONFromString(basicJSON),
			args: args{
				path: "name.last",
			},
			want: true,
		},
		{
			name: "2",
			j:    testJSONFromString(basicJSON),
			args: args{
				path: "arr",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.j.IsNull(tt.args.path); got != tt.want {
				t.Errorf("JSON.IsNull() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJson_Exists(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name string
		j    *JSON
		args args
		want bool
	}{
		{
			name: "1",
			j:    testJSONFromString(basicJSON),
			args: args{
				path: "name.last",
			},
			want: true,
		},
		{
			name: "2",
			j:    testJSONFromString(basicJSON),
			args: args{
				path: "arr.100",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.j.Exists(tt.args.path); got != tt.want {
				t.Errorf("JSON.Exists() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJson_Set(t *testing.T) {
	type args struct {
		path string
		v    interface{}
	}
	tests := []struct {
		name    string
		j       *JSON
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "1",
			j:    testJSONFromString(`{"name":"x"}`),
			args: args{
				path: "name.first",
				v:    nil,
			},
			want:    `{"name":{"first":null}}`,
			wantErr: false,
		},
		{
			name: "2",
			j:    testJSONFromString(`{"name":"x"}`),
			args: args{
				path: "name.first.1",
				v:    nil,
			},
			want:    `{"name":{"first":[null,null]}}`,
			wantErr: false,
		},
		{
			name: "3",
			j:    testJSONFromString(`{"name":"x"}`),
			args: args{
				path: "name.first.*",
				v:    nil,
			},
			want:    `{"name":"x"}`,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.j.Set(tt.args.path, tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("JSON.Set() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.j.String() != tt.want {
				t.Errorf("j.String() %v, want %v", tt.j, tt.want)
			}
		})
	}
}

func TestJson_SetRawBytes(t *testing.T) {
	type args struct {
		path string
		b    []byte
	}
	tests := []struct {
		name    string
		j       *JSON
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "1",
			j:    testJSONFromString(`{"name":"x"}`),
			args: args{
				path: "name.first",
				b:    []byte(`[1,2,3,4,5,67]`),
			},
			want:    `{"name":{"first":[1,2,3,4,5,67]}}`,
			wantErr: false,
		},
		{
			name: "2",
			j:    testJSONFromString(`{"name":"x"}`),
			args: args{
				path: "name.first.*",
				b:    []byte("1"),
			},
			want:    `{"name":"x"}`,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.j.SetRawBytes(tt.args.path, tt.args.b); (err != nil) != tt.wantErr {
				t.Errorf("JSON.SetRawBytes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.j.String() != tt.want {
				t.Errorf("j.String() %v, want %v", tt.j, tt.want)
			}
		})
	}
}

func TestJson_SetRawString(t *testing.T) {
	type args struct {
		path string
		s    string
	}
	tests := []struct {
		name    string
		j       *JSON
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "1",
			j:    testJSONFromString(`{"name":"x"}`),
			args: args{
				path: "name.first",
				s:    `[1,2,3,4,5,67]`,
			},
			want:    `{"name":{"first":[1,2,3,4,5,67]}}`,
			wantErr: false,
		},
		{
			name: "2",
			j:    testJSONFromString(`{"name":"x"}`),
			args: args{
				path: "name.first.*",
				s:    "1",
			},
			want:    `{"name":"x"}`,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.j.SetRawString(tt.args.path, tt.args.s); (err != nil) != tt.wantErr {
				t.Errorf("JSON.SetRawString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.j.String() != tt.want {
				t.Errorf("j.String() %v, want %v", tt.j, tt.want)
			}
		})
	}
}

func TestJson_Remove(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		j       *JSON
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "1",
			j:    testJSONFromString(`{"name":"x"}`),
			args: args{
				path: "name",
			},
			want:    `{}`,
			wantErr: false,
		},
		{
			name: "2",
			j:    testJSONFromString(`{"name":"x"}`),
			args: args{
				path: "name.first.*",
			},
			want:    `{"name":"x"}`,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.j.Remove(tt.args.path); (err != nil) != tt.wantErr {
				t.Errorf("JSON.Remove() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.j.String() != tt.want {
				t.Errorf("j.String() %v, want %v", tt.j, tt.want)
			}
		})
	}
}

func TestJson_FromString(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		j       *JSON
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "1",
			j:    &JSON{},
			args: args{
				s: invlaidJSON,
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "2",
			j:    &JSON{},
			args: args{
				s: basicJSON,
			},
			want:    basicJSON,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.j.FromString(tt.args.s); (err != nil) != tt.wantErr {
				t.Errorf("JSON.FromString() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.j.String() != tt.want {
				t.Errorf("j.String() %v, want %v", tt.j, tt.want)
			}
		})
	}
}

func TestJson_FromBytes(t *testing.T) {
	type args struct {
		b []byte
	}
	tests := []struct {
		name    string
		j       *JSON
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "1",
			j:    &JSON{},
			args: args{
				b: []byte(invlaidJSON),
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "2",
			j:    &JSON{},
			args: args{
				b: []byte(basicJSON),
			},
			want:    basicJSON,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.j.FromBytes(tt.args.b); (err != nil) != tt.wantErr {
				t.Errorf("JSON.FromBytes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.j.String() != tt.want {
				t.Errorf("j.String() %v, want %v", tt.j, tt.want)
			}
		})
	}
}

func TestJson_FromFile(t *testing.T) {
	type args struct {
		filename string
	}
	tests := []struct {
		name    string
		j       *JSON
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "1",
			j:    &JSON{},
			args: args{
				"1213231k3j21kl3.dadsadasda",
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "2",
			j:    &JSON{},
			args: args{
				"test_data",
			},
			want: strings.ReplaceAll(basicJSON, "\n", "\r\n"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.j.FromFile(tt.args.filename); (err != nil) != tt.wantErr {
				t.Errorf("JSON.FromFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.j.String() != tt.want {
				t.Errorf("j.String() %v, want %v", tt.j, tt.want)
			}
		})
	}
}

func TestJson_Clone(t *testing.T) {
	tests := []struct {
		name string
		j    *JSON
		want *JSON
	}{
		{
			name: "1",
			j:    testJSONFromString(basicJSON),
			want: testJSONFromString(basicJSON),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.j.Clone()
			if got == tt.j {
				t.Errorf("JSON.Clone() = %p, j %p", got, tt.want)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("JSON.Clone() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJson_MarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		j       *JSON
		want    []byte
		wantErr bool
	}{
		{
			name: "1",
			j:    testJSONFromString(basicJSON),
			want: []byte(basicJSON),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.j.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("JSON.MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("JSON.MarshalJSON() = %v, want %v", got, tt.want)
			}
		})
	}

	testJSON := struct {
		MyJSON *JSON `JSON:"myJson"`
		Age    int   `JSON:"age"`
	}{
		MyJSON: testJSONFromString(basicJSON),
		Age:    11,
	}

	b, err := json.Marshal(testJSON)
	if err != nil {
		t.Error(err)
	}

	t.Log(string(b))
}
