# Installation

```bash
go get github.com/golangblogs/elasticsearch
```

# Example EsBulkAdd

```go
import "github.com/golangblogs/elasticsearch"

/ EsBulkAdd 批量添加文档
// @router /es/bulk/add [post]
func (this *EsDemoController) EsBulkAdd() {
        //从数据库获取数据 是一个切片
        video := models.GetVideoList()
        //定义存储桶的数据
        var bulkBody []interface{}
        for _, v := range video {
                //拼接要导入的数据 得拼接
                body := map[string]interface{}{
                    "id":                   v.Id,
                    "title":                v.Title,
                    "sub_title":            v.SubTitle,
                    "add_time":             v.AddTime,
                    "img":                  v.Img,
                    "img1":                 v.Img1,
                    "episodes_count":       v.EpisodesCount,
                    "is_end":               v.IsEnd,
                    "channel_id":           v.ChannelId,
                    "status":               v.Status,
                    "region_id":            v.RegionId,
                    "type_id":              v.TypeId,
                    "episodes_update_time": v.EpisodesUpdateTime,
                    "comment":              v.Comment,
                    "user_id":              v.UserId,
                    "is_recommend":         v.IsRecommend,
                }
                bulkBody = append(bulkBody, body)
        }
		//index_name 索引的名字
        elasticsearch.EsBulkAdd("index_name", bulkBody)
        //to do something
}

```


# Example EsSearch

```go
import "github.com/golangblogs/elasticsearch"

// @router /es/search [post]
func (this *EsDemoController) EsSearch() {
	//查询的条件
	title := this.GetString("title")
	//第几页
	page, _ := this.GetInt("page", 1)
	//查几条，每页显示条数
	size, _ := this.GetInt("limit", 10)
	//偏移量
	from := (page - 1) * size
	//拼接查询条件
	query := map[string]interface{}{
		"bool": map[string]interface{}{
			"must": []map[string]interface{}{
				map[string]interface{}{
					"match": map[string]interface{}{
						"title": title,
					},
				},
			},
		},
	}

	//排序条件
	sort := []map[string]string{map[string]string{"id": "desc"}}
	res := elasticsearch.EsSearch("index_name", query, from, size, sort)
	//转结构体
	//符合搜索条件总条数
	total := res.Hits.Total.Value
	//定义一个视频切片，存放搜索完数据
	var videos []models.VideoData

	for _, v := range res.Hits.Hits {
		tmp := models.VideoData{}
		json.Unmarshal(v.Source, &tmp)
		videos = append(videos, tmp)
	}

	this.Data["json"] = map[string]interface{}{
		"code":  0,
		"total": total,
		"data":  videos,
		"msg":   "ok",
	}
	this.ServeJSON()

}

```