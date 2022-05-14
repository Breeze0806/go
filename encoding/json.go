package encoding

import (
	"io/ioutil"
	"strconv"

	"github.com/pingcap/errors"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

var (
	//ErrJSONNotValid 不是合法的json
	ErrJSONNotValid = errors.NewNoStackError("json is not valid")
)

//JSON json格式编码
type JSON struct {
	res gjson.Result
}

//NewJSONFromString 从字符串s中获取JSON，该函数会检查格式合法性
func NewJSONFromString(s string) (*JSON, error) {
	if !gjson.Valid(s) {
		return nil, errors.Wrapf(ErrJSONNotValid, "json: %v", s)
	}
	return newJSONFromString(s), nil
}

//NewJSONFromString 从字符串s中获取JSON，该函数不会检查格式合法性
func newJSONFromString(s string) *JSON {
	return &JSON{
		res: gjson.Parse(s),
	}
}

//NewJSONFromBytes 从字符流中b中获取JSON，该函数会检查格式合法性
func NewJSONFromBytes(b []byte) (*JSON, error) {
	if !gjson.ValidBytes(b) {
		return nil, errors.Wrapf(ErrJSONNotValid, "json: %v", string(b))
	}
	return newJSONFromBytes(b), nil
}

//newJSONFromBytes 从字符流b中获取JSON，该函数不会检查格式合法性
func newJSONFromBytes(b []byte) *JSON {
	return &JSON{
		res: gjson.ParseBytes(b),
	}
}

//NewJSONFromFile 从文件filename中获取JSON，该函数会检查格式合法性
func NewJSONFromFile(filename string) (*JSON, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, errors.Wrapf(err, "read file %v fail.", filename)
	}
	return NewJSONFromBytes(data)
}

//GetJSON 获取path路径对应的值JOSN结构,对于下列json
//{
//  "a":{
//     "b":[{
//        c:"x"
//      }]
//	}
//}
//要访问到x字符串 path每层的访问路径为a,a.b,a.b.0，a.b.0.c
//如果path对应的不是json结构或者不存在，就会返回错误
func (j *JSON) GetJSON(path string) (*JSON, error) {
	res, err := j.getResult(path)
	if err != nil {
		return nil, errors.Wrapf(err, "getResult(%v) fail.", path)
	}
	if res.Type != gjson.JSON {
		return nil, errors.Errorf("path(%v) is not json", path)
	}

	return &JSON{
		res: res,
	}, nil
}

//GetBool 获取path路径对应的值bool值,对于下列json
//{
//  "a":{
//     "b":[{
//        c:"x"
//      }]
//	}
//}
//要访问到x字符串 path每层的访问路径为a,a.b,a.b.0，a.b.0.c
//如果path对应的不是bool值或者不存在，就会返回错误
func (j *JSON) GetBool(path string) (bool, error) {
	res, err := j.getResult(path)
	if err != nil {
		return false, errors.Wrapf(err, "getResult(%v) fail.", path)
	}
	switch res.Type {
	case gjson.False:
		return false, nil
	case gjson.True:
		return true, nil
	}
	return false, errors.Errorf("path(%v) is not bool", path)
}

//GetInt64 获取path路径对应的值int64值,对于下列json
//{
//  "a":{
//     "b":[{
//        c:"x"
//      }]
//	}
//}
//要访问到x字符串 path每层的访问路径为a,a.b,a.b.0，a.b.0.c
//如果path对应的不是int64值或者不存在，就会返回错误
func (j *JSON) GetInt64(path string) (int64, error) {
	res, err := j.getResult(path)
	if err != nil {
		return 0, errors.Wrapf(err, "getResult(%v) fail.", path)
	}
	switch res.Type {
	case gjson.Number:
		v, err := strconv.ParseInt(res.String(), 10, 64)
		if err != nil {
			return 0, errors.Wrapf(err, "path(%v) is not int64. val: %v", path, res.String())
		}
		return v, nil
	}
	return 0, errors.Errorf("path(%v) is not bool", path)
}

//GetFloat64 获取path路径对应的值float64值,对于下列json
//{
//  "a":{
//     "b":[{
//        c:"x"
//      }]
//	}
//}
//要访问到x字符串 path每层的访问路径为a,a.b,a.b.0，a.b.0.c
//如果path对应的不是float64值或者不存在，就会返回错误
func (j *JSON) GetFloat64(path string) (float64, error) {
	res, err := j.getResult(path)
	if err != nil {
		return 0, errors.Wrapf(err, "getResult(%v) fail.", path)
	}
	switch res.Type {
	case gjson.Number:
		v, err := strconv.ParseFloat(res.String(), 64)
		if err != nil {
			return 0, errors.Wrapf(err, "path(%v) is not float64. val: %v", path, res.String())
		}
		return v, nil
	}
	return 0, errors.Errorf("path(%v) is not bool", path)
}

//GetString 获取path路径对应的值String值,对于下列json
//{
//  "a":{
//     "b":[{
//        c:"x"
//      }]
//	}
//}
//要访问到x字符串 path每层的访问路径为a,a.b,a.b.0，a.b.0.c
//如果path对应的不是String值或者不存在，就会返回错误
func (j *JSON) GetString(path string) (string, error) {
	res, err := j.getResult(path)
	if err != nil {
		return "", errors.Wrapf(err, "getResult(%v) fail.", path)
	}
	switch res.Type {
	case gjson.String:
		return res.String(), nil
	}
	return "", errors.Errorf("path(%v) is not string", path)
}

//GetArray 获取path路径对应的值数组,对于下列json
//{
//  "a":{
//     "b":[{
//        c:"x"
//      }]
//	}
//}
//要访问到x字符串 path每层的访问路径为a,a.b,a.b.0，a.b.0.c
//如果path对应的不是数组或者不存在，就会返回错误
func (j *JSON) GetArray(path string) ([]*JSON, error) {
	res, err := j.getResult(path)
	if err != nil {
		return nil, errors.Wrapf(err, "getResult(%v) fail.", path)
	}
	switch {
	case res.IsArray():
		var jsons []*JSON
		a := res.Array()
		for _, v := range a {
			jsons = append(jsons, &JSON{res: v})
		}
		return jsons, nil
	}
	return nil, errors.Errorf("path(%v) is not array", path)
}

//GetMap 获取path路径对应的值字符串映射,对于下列json
//{
//  "a":{
//     "b":[{
//        c:"x"
//      }]
//	}
//}
//要访问到x字符串 path每层的访问路径为a,a.b,a.b.0，a.b.0.c
//如果path对应的不是字符串映射或者不存在，就会返回错误
func (j *JSON) GetMap(path string) (map[string]*JSON, error) {
	res, err := j.getResult(path)
	if err != nil {
		return nil, errors.Wrapf(err, "getResult(%v) fail.", path)
	}
	switch {
	case res.IsObject():
		jsons := make(map[string]*JSON)
		m := res.Map()
		for k, v := range m {
			jsons[k] = &JSON{res: v}
		}
		return jsons, nil
	}
	return nil, errors.Errorf("path(%v) is not map", path)
}

//String 获取字符串表示
func (j *JSON) String() string {
	return j.res.Raw
}

//IsArray 判断path路径对应的值是否是数组,对于下列json
//{
//  "a":{
//     "b":[{
//        c:"x"
//      }]
//	}
//}
//要访问到x字符串 path每层的访问路径为a,a.b,a.b.0，a.b.0.c
func (j *JSON) IsArray(path string) bool {
	return j.res.Get(path).IsArray()
}

//IsNumber 判断path路径对应的值是否是数值,对于下列json
//{
//  "a":{
//     "b":[{
//        c:"x"
//      }]
//	}
//}
//要访问到x字符串 path每层的访问路径为a,a.b,a.b.0，a.b.0.c
func (j *JSON) IsNumber(path string) bool {
	return j.res.Get(path).Type == gjson.Number
}

//IsJSON 判断path路径对应的值是否是JSON,对于下列json
//{
//  "a":{
//     "b":[{
//        c:"x"
//      }]
//	}
//}
//要访问到x字符串 path每层的访问路径为a,a.b,a.b.0，a.b.0.c
func (j *JSON) IsJSON(path string) bool {
	return j.res.Get(path).IsObject()
}

//IsBool 判断path路径对应的值是否是BOOL,对于下列json
//{
//  "a":{
//     "b":[{
//        c:"x"
//      }]
//	}
//}
//要访问到x字符串 path每层的访问路径为a,a.b,a.b.0，a.b.0.c
func (j *JSON) IsBool(path string) bool {
	switch j.res.Get(path).Type {
	case gjson.False, gjson.True:
		return true
	}
	return false
}

//IsString 判断path路径对应的值是否是字符串,对于下列json
//{
//  "a":{
//     "b":[{
//        c:"x"
//      }]
//	}
//}
//要访问到x字符串 path每层的访问路径为a,a.b,a.b.0，a.b.0.c
func (j *JSON) IsString(path string) bool {
	return j.res.Get(path).Type == gjson.String
}

//IsNull 判断path路径对应的值值是否为空,对于下列json
//{
//  "a":{
//     "b":[{
//        c:"x"
//      }]
//	}
//}
//要访问到x字符串 path每层的访问路径为a,a.b,a.b.0，a.b.0.c
func (j *JSON) IsNull(path string) bool {
	return j.res.Get(path).Type == gjson.Null
}

//Exists 判断path路径对应的值值是否存在,对于下列json
//{
//  "a":{
//     "b":[{
//        c:"x"
//      }]
//	}
//}
//要访问到x字符串 path每层的访问路径为a,a.b,a.b.0，a.b.0.c
func (j *JSON) Exists(path string) bool {
	return j.res.Get(path).Exists()
}

//Set 将path路径对应的值设置成v,会返回错误error,对于下列json
//{
//  "a":{
//     "b":[{
//        c:"x"
//      }]
//	}
//}
//要访问到x字符串 path每层的访问路径为a,a.b,a.b.0，a.b.0.c
func (j *JSON) Set(path string, v interface{}) error {
	s, err := sjson.Set(j.String(), path, v)
	if err != nil {
		return errors.Wrapf(err, "path(%v) set fail. val: %v", path, v)
	}
	j.fromString(s)
	return nil
}

//SetRawBytes 将path路径对应的值设置成b,会返回错误error,对于下列json
//{
//  "a":{
//     "b":[{
//        c:"x"
//      }]
//	}
//}
//要访问到x字符串 path每层的访问路径为a,a.b,a.b.0，a.b.0.c
func (j *JSON) SetRawBytes(path string, b []byte) error {
	return j.SetRawString(path, string(b))
}

//SetRawString 将path路径对应的值设置成s,会返回错误error,对于下列json
//{
//  "a":{
//     "b":[{
//        c:"x"
//      }]
//	}
//}
//要访问到x字符串 path每层的访问路径为a,a.b,a.b.0，a.b.0.c
func (j *JSON) SetRawString(path string, s string) error {
	ns, err := sjson.SetRaw(j.String(), path, s)
	if err != nil {
		return errors.Wrapf(err, "path(%v) set fail. val: %v", path, s)
	}
	j.fromString(ns)
	return nil
}

//Remove 将path路径对应的值删除,会返回错误error,对于下列json
//{
//  "a":{
//     "b":[{
//        c:"x"
//      }]
//	}
//}
//要访问到x字符串 path每层的访问路径为a,a.b,a.b.0，a.b.0.c
func (j *JSON) Remove(path string) error {
	s, err := sjson.Delete(j.String(), path)
	if err != nil {
		return errors.Wrapf(err, "path(%v) remove fail.", path)
	}
	j.fromString(s)
	return nil
}

//FromString 将字符串s设置成JSON,会返回错误error
func (j *JSON) FromString(s string) error {
	new, err := NewJSONFromString(s)
	if err != nil {
		return err
	}
	j.res = new.res
	return nil
}

//FromBytes 将字节流b设置成JSON,会返回错误error
func (j *JSON) FromBytes(b []byte) error {
	new, err := NewJSONFromBytes(b)
	if err != nil {
		return err
	}
	j.res = new.res
	return nil
}

//FromFile 从文件名为filename的文件中读取JSON,会返回错误error
func (j *JSON) FromFile(filename string) error {
	new, err := NewJSONFromFile(filename)
	if err != nil {
		return err
	}
	j.res = new.res
	return nil
}

//Clone 克隆JSON
func (j *JSON) Clone() *JSON {
	return &JSON{
		res: j.res,
	}
}

//MarshalJSON 序列化JSON
func (j *JSON) MarshalJSON() ([]byte, error) {
	return []byte(j.res.Raw), nil
}

func (j *JSON) fromString(s string) {
	new := newJSONFromString(s)
	j.res = new.res
}

func (j *JSON) getResult(path string) (gjson.Result, error) {
	res := j.res.Get(path)
	if res.Exists() {
		return res, nil
	}
	return gjson.Result{}, errors.Errorf("path(%v) does not exist", path)
}
