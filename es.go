package es

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/httplib"
)

// elasticSearch连接地址
var esUrl string

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

func init() {
	//9200后面的/不要少
	esUrl = "http://127.0.0.1:9200/"
}

// 搜索功能
// sort []map[string]string 根据哪个字段排序 正序还是倒叙
func EsSearch(indexName string, query map[string]interface{}, from int, size int, sort []map[string]string) (HitsData, error) {
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
	var stb ReqSearchData
	if err != nil {
		return stb.Hits, err
	}

	//把获取的json值解析一下
	err = json.Unmarshal([]byte(str), &stb)
	if err != nil {
		return stb.Hits, err
	}
	return stb.Hits, err
}

// EsAdd  添加
func EsAdd(indexName string, id string, body map[string]interface{}) (bool, error) {
	//7.x版本type已经取消了，否则会报错
	req := httplib.Post(esUrl + indexName + "/_doc/" + id)

	//searchQuery参数通过json的格式传过去
	req.JSONBody(body)

	//最后获取返回值
	_, err := req.String()
	if err != nil {
		return false, err
	}
	return true, nil
}

// EsEdit 修改
func EsEdit(indexName string, id string, body map[string]interface{}) (bool, error) {

	bodyData := map[string]interface{}{
		//要把修改的内容放到doc中
		//要不然会报错
		"doc": body,
	}

	req := httplib.Post(esUrl + indexName + "/_doc/" + id + "/_update")
	req.JSONBody(bodyData)

	//最后获取返回值
	_, err := req.String()
	if err != nil {
		return false, err
	}
	return true, nil
}

// EsDelete 删除
func EsDelete(indexName string, id string) (bool, error) {
	req := httplib.Delete(esUrl + indexName + "/_doc/" + id)
	//最后获取返回值
	_, err := req.String()
	if err != nil {
		return false, err
	}
	return true, nil
}

// EsBulkAdd 实现批量添加数据
func EsBulkAdd(indexName string, data []interface{}) (bool, error) {
	//httplib请求包
	req := httplib.Post(esUrl + indexName + "/_bulk")
	// 设置请求头Content-Type
	req.Header("Content-Type", "application/x-ndjson")

	//参数通过json的格式传过去
	var bulkData string

	for _, item := range data {
		//这段代码首先尝试将每个数据项序列化为 JSON 字符串。如果序列化过程有误，它将跳过该数据项并继续处理其它数据
		itemJSON, err := json.Marshal(item)
		if err != nil {
			fmt.Printf("Error marshaling data item: %s\n", err)
			continue
		}
		bulkData += fmt.Sprintf("{\"index\":{}}\n%s\n", string(itemJSON))
	}

	req.Body([]byte(bulkData))

	//req.String()这个不能少
	_, err := req.String()
	if err != nil {
		return false, err
	}
	return true, nil
}
