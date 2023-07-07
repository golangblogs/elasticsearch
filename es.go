package es

import (
	"encoding/json"
	"fmt"

	"github.com/astaxie/beego/httplib"
)

var esUrl string

func init() {
	//9200后面的/不要少
	esUrl = "http://127.0.0.1:9200/"
}

// 搜索功能
// sort []map[string]string 根据哪个字段排序 正序还是倒叙
func EsSearch(indexName string, query map[string]interface{}, from int, size int, sort []map[string]string) HitsData {
	searchQuery := map[string]interface{}{
		"query": query,
		"from":  from,
		"size":  size,
		"sort":  sort,
	}
	//httplib请求包
	req := httplib.Post(esUrl + indexName + "/_search")
	//searchQuery参数通过json的格式传过去
	req.JSONBody(searchQuery)
	//最后获取返回值
	str, err := req.String()
	fmt.Println(str)
	if err != nil {
		fmt.Println(err)
	}
	var stb ReqSearchData
	//把获取的json值解析一下
	err = json.Unmarshal([]byte(str), &stb)

	return stb.Hits
}

// 解析获取到的值
type ReqSearchData struct {
	Hits HitsData `json:"hits"`
}

type HitsData struct {
	Total TotalData     `json:"total"`
	Hits  []HitsTwoData `json:"hits"`
}

type HitsTwoData struct {
	Source json.RawMessage `json:"_source"`
}

type TotalData struct {
	Value    int
	Relation string
}

// 添加
func EsAdd(indexName string, id string, body map[string]interface{}) bool {
	//7.x版本type已经取消了，否则会报错
	req := httplib.Post(esUrl + indexName + "/_doc/" + id)

	//searchQuery参数通过json的格式传过去
	req.JSONBody(body)

	//最后获取返回值
	str, err := req.String()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(">>>>>>>>>>>" + str)
	return true
}

// 修改
func EsEdit(indexName string, id string, body map[string]interface{}) bool {

	bodyData := map[string]interface{}{
		//要把修改的内容放到doc中
		//要不然会报错
		"doc": body,
	}

	req := httplib.Post(esUrl + indexName + "/_doc/" + id + "/_update")
	req.JSONBody(bodyData)

	str, err := req.String()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(str)
	return true
}

// 删除
func EsDelete(indexName string, id string) bool {
	req := httplib.Delete(esUrl + indexName + "/_doc/" + id)
	//最后获取返回值
	str, err := req.String()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(str)
	return true
}
