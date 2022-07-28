package service

import (
	"context"
	"demo/db"
	"demo/db/dao"
	"demo/db/model"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"os"
	"strconv"

	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"gorm.io/gorm"
)

// JsonResult 返回结构
type JsonResult struct {
	Code      int         `json:"code"`
	ErrorMsg  string      `json:"errorMsg,omitempty"`
	Data      interface{} `json:"data"`
	RedisData interface{} `json:"redis_data"`
	MongoData interface{} `json:"mongo_data"`
	Count     int         `json:"count"`
}

// FollowListResult  get_follow_list 返回结构
type FollowListResult struct {
	Data  FollowListData `json:"data"`
	Extra interface{}    `json:"extra"`
}

// FollowListData FollowListResult的Data子项结构
type FollowListData struct {
	ErrorCode   int           `json:"error_code"`
	Description string        `json:"description"`
	Cursor      int           `json:"cursor"`
	List        []interface{} `json:"list"`
}

// IndexHandler 计数器接口
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	data, err := getIndex()
	if err != nil {
		fmt.Fprint(w, "内部错误")
		return
	}
	fmt.Fprint(w, data)
}

// testHandlerr 测试接口，返回全部的header和query参数
func TestHandler(w http.ResponseWriter, r *http.Request) {
	data := make(map[string]string)
	fmt.Println("打印Header参数列表：")
	if len(r.Header) > 0 {
		for k, v := range r.Header {
			fmt.Printf("%s=%s\n", k, v[0])
			data[k] = v[0]
		}
	}
	fmt.Println("打印Form参数列表：")
	r.ParseForm()
	if len(r.Form) > 0 {
		for k, v := range r.Form {
			fmt.Printf("%s=%s\n", k, v[0])
			data[k] = v[0]
		}
	}

	scheme := "http://"
	if r.TLS != nil {
		scheme = "https://"
	}

	url := strings.Join([]string{scheme, r.Host, r.RequestURI}, "")
	data["url"] = url

	msg, err := json.Marshal(data)
	if err != nil {
		fmt.Fprint(w, "内部错误")
		return
	}
	w.Header().Set("content-type", "application/json")
	w.Write(msg)
}

// ErrorTestHandler 测试错误信息
func ErrorTestHandler(w http.ResponseWriter, r *http.Request) {
	statusStr := r.URL.Query().Get("status_id")
	statusId, err := strconv.Atoi(statusStr)
	if err != nil {
		statusId = 500
	}
	w.WriteHeader(statusId)
}

// CounterHandler 计数器接口
func CounterHandler(w http.ResponseWriter, r *http.Request) {
	res := &JsonResult{}

	if r.Method == http.MethodGet {
		redisCounter, err2 := getRedisCurrentCounter()
		upsertCounterRedis(r)
		mongoCounter, err3 := getMongoCurrentCounter()
		upsertCounterMongo(r)
		if err2 != nil || err3 != nil {
			res.Code = -1
			res.ErrorMsg = err2.Error()
		} else {
			res.RedisData = redisCounter.Count
			res.MongoData = mongoCounter.Count
		}
	} else if r.Method == http.MethodPost {
		modifyCounter(r)
		//counter, err := getCurrentCounter()
		redisCounter, err2 := getRedisCurrentCounter()
		mongoCounter, err3 := getMongoCurrentCounter()
		if err2 != nil || err3 != nil {
			res.Code = -1
			res.ErrorMsg = err2.Error()
		} else {
			//res.Data = counter.Count
			res.RedisData = redisCounter.Count
			res.MongoData = mongoCounter.Count
		}
	} else {
		res.Code = -1
		res.ErrorMsg = fmt.Sprintf("请求方法 %s 不支持", r.Method)
	}

	msg, err := json.Marshal(res)
	if err != nil {
		fmt.Fprint(w, "内部错误")
		return
	}
	w.Header().Set("content-type", "application/json")
	w.Write(msg)
}

// FollowListHandler 获取关注列表接口
func FollowListHandler(w http.ResponseWriter, r *http.Request) {
	res := &JsonResult{}

	domain := "http://open.douyin.com"
	path := "/following/list/"

	client := http.Client{Timeout: 1000 * time.Millisecond}

	req, err := http.NewRequest(http.MethodGet, domain+path, nil)
	if err != nil {
		fmt.Fprint(w, "内部错误")
		return
	}
	openId := r.Header.Get("X-Tt-Openid")
	if openId == "" {
		fmt.Fprint(w, "X-Tt-Openid 为空")
		return
	}
	query := req.URL.Query()
	query.Add("open_id", openId)
	query.Add("cursor", "0")
	query.Add("count", "30")
	req.URL.RawQuery = query.Encode()
	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		fmt.Fprint(w, "call openapi error")
		return
	}

	followBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprint(w, "内部错误")
		return
	}
	var followList FollowListResult
	err = json.Unmarshal(followBody, &followList)
	if err != nil {
		fmt.Fprint(w, "内部错误")
		return
	}
	if followList.Data.ErrorCode != 0 {
		res.Count = 0
		res.ErrorMsg = followList.Data.Description
	}
	if followList.Data.List != nil {
		res.Count = len(followList.Data.List)
	}
	msg, err := json.Marshal(res)
	if err != nil {
		fmt.Fprint(w, "内部错误")
		return
	}
	w.Header().Set("content-type", "application/json")
	w.Write(msg)
}

// TestEndGateHandler 获取关注列表接口(直接返回抖开的body)
func TestEndGateHandler(w http.ResponseWriter, r *http.Request) {

	// domain := "http://douyincloud.gateway.egress.ivolces.com"
	domain := "https://developer.toutiao.com"
	path := "/api/v2/tags/text/antidirt"

	//payload := strings.NewReader(`{"access_token": "0801121846765a5a4d2f6b385a68307237534d43397a667865513d3d","appname": "douyin"}`)
	payloadWithoutToken := strings.NewReader(`{"tasks": [{"content": "要检测的文本"}]}`)

	client := http.Client{Timeout: 1000 * time.Millisecond}

	req, err := http.NewRequest(http.MethodPost, domain+path, payloadWithoutToken)
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		fmt.Fprint(w, "内部错误")
		return
	}
	openId := r.Header.Get("X-Tt-Openid")
	if openId == "" {
		fmt.Fprint(w, "X-Tt-Openid 为空")
		return
	}
	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		fmt.Fprint(w, "call openapi error")
		return
	}

	followBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprint(w, "内部错误")
		return
	}
	w.Header().Set("content-type", "application/json")
	w.Write(followBody)
}

// TestFollowListHandler 获取关注列表接口(直接返回抖开的body)
func TestFollowListHandler(w http.ResponseWriter, r *http.Request) {

	domain := "https://open.douyin.com"
	path := "/following/list/"

	client := http.Client{Timeout: 1000 * time.Millisecond}

	req, err := http.NewRequest(http.MethodGet, domain+path, nil)
	if err != nil {
		fmt.Fprint(w, "内部错误")
		return
	}
	openId := r.Header.Get("X-Tt-Openid")
	if openId == "" {
		fmt.Fprint(w, "X-Tt-Openid 为空")
		return
	}
	query := req.URL.Query()
	query.Add("open_id", openId)
	query.Add("cursor", "0")
	query.Add("count", "30")
	req.URL.RawQuery = query.Encode()
	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		fmt.Fprint(w, "call openapi error")
		return
	}

	followBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprint(w, "内部错误")
		return
	}
	w.Header().Set("content-type", "application/json")
	w.Write(followBody)
}
func TestSleepHandler(w http.ResponseWriter, r *http.Request) {
	sleepTimeStr := r.URL.Query().Get("sleep_time")
	sleepTimeId, err := strconv.Atoi(sleepTimeStr)
	if err != nil {
		sleepTimeId = 10
	}
	time.Sleep(time.Duration(sleepTimeId) * time.Second)
	str := fmt.Sprintf("sleep %d s done\n", sleepTimeId)
	w.Write([]byte(str))
}

func PingHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "pong!\n")
}

func GetOsEnvHandler(w http.ResponseWriter, r *http.Request) {
	redisAddr := os.Getenv("REDIS_ADDRESS")
	mongoUrl := os.Getenv("MONGO_ADDRESS")
	osEnv := fmt.Sprintf("REDIS_ADDRESS=%s\nMONGO_URL=%s\n", redisAddr, mongoUrl)
	w.Write([]byte(osEnv))
}

// modifyCounter 更新计数，自增或者清零
func modifyCounter(r *http.Request) (int32, error) {
	action, dbType, err := getAction(r)
	fmt.Println(action, dbType, err)
	if err != nil {
		return 0, err
	}
	fmt.Println(action, dbType, err)
	var count int32
	if action == "inc" {
		if dbType == "mysql" {
			fmt.Println("inc redis count")
			//count, err = upsertCounter(r)
			//if err != nil {
			//	return 0, err
			//}
		} else if dbType == "redis" {
			fmt.Println("inc redis count")
			count, err = upsertCounterRedis(r)
			if err != nil {
				return 0, err
			}
		} else if dbType == "mongo" {
			fmt.Println("inc mongo count")
			count, err = upsertCounterMongo(r)
			if err != nil {
				return 0, err
			}
		}

	} else if action == "clear" {
		err = clearCounter()
		if err != nil {
			return 0, err
		}
		count = 0
	} else {
		err = fmt.Errorf("参数 action : %s 错误", action)
	}

	return count, err
}

// upsertCounter 更新或修改计数器
func upsertCounter(r *http.Request) (int32, error) {
	currentCounter, err := getCurrentCounter()
	var count int32
	createdAt := time.Now()
	if err != nil && err != gorm.ErrRecordNotFound {
		return 0, err
	} else if err == gorm.ErrRecordNotFound {
		count = 1
		createdAt = time.Now()
	} else {
		count = currentCounter.Count + 1
		createdAt = currentCounter.CreatedAt
	}

	counter := &model.CounterModel{
		Id:        1,
		Count:     count,
		CreatedAt: createdAt,
		UpdatedAt: time.Now(),
	}
	err = dao.Imp.UpsertCounter(counter)
	if err != nil {
		return 0, err
	}
	return counter.Count, nil
}

// upsertCounter 更新或修改计数器
func upsertCounterRedis(r *http.Request) (int32, error) {
	var ctx = context.Background()
	rdb := db.GetRedis()
	key := "count"
	val, err := rdb.Get(ctx, key).Result()
	if err != nil {
		panic(err)
	}
	ans, err := strconv.Atoi(val)
	if err != nil {
		panic(err)
	}
	err = rdb.Set(ctx, key, ans+1, 0).Err()
	if err != nil {
		panic(err)
	}
	return int32(ans + 1), nil
}

// upsertCounter 更新或修改计数器
func upsertCounterMongo(r *http.Request) (int32, error) {
	client := db.GetMongo()
	coll := client.Database(db.DataBase).Collection("count")

	doc := &model.MongoCount{}

	filter := bson.M{"type": "mongodb"}
	result := coll.FindOne(context.TODO(), filter)

	fmt.Printf("documents upsertCounterMongo with result:%v\n", result)
	if err := result.Decode(doc); err != nil {
		fmt.Println("upsertCounterMongo decode err", err)
		return 0, err
	}
	fmt.Println(doc)

	update := bson.M{"$set": model.MongoCount{Type: "mongodb", Count: doc.Count + 1}}

	uResult, err := coll.UpdateMany(context.TODO(), filter, update)
	if err != nil {
		fmt.Println("upsertCounterMongo update err", err)
	}
	fmt.Println(uResult.MatchedCount)

	return int32(doc.Count + 1), nil
}

func clearCounter() error {
	return dao.Imp.ClearCounter(1)
}

// getCurrentCounter 查询当前计数器
func getCurrentCounter() (*model.CounterModel, error) {
	counter, err := dao.Imp.GetCounter(1)
	if err != nil {
		return nil, err
	}

	return counter, nil
}

// getCurrentCounter redis 查询当前计数器
func getRedisCurrentCounter() (*model.CounterModel, error) {
	var ctx = context.Background()
	rdb := db.GetRedis()
	fmt.Println(rdb)
	key := "count"
	val, err := rdb.Get(ctx, key).Result()
	if err != nil {
		fmt.Println("getRedisCurrentCounter err ", err)
		return nil, err
	}
	ans, err := strconv.Atoi(val)
	if err != nil {
		return nil, err
	}
	counter := &model.CounterModel{
		Id:    0,
		Count: int32(ans),
	}
	return counter, nil
}

// getCurrentCounter mongo 查询当前计数器
func getMongoCurrentCounter() (*model.CounterModel, error) {
	client := db.GetMongo()
	coll := client.Database(db.DataBase).Collection("count")

	doc := &model.MongoCount{}

	filter := bson.M{"type": "mongodb"}
	result := coll.FindOne(context.TODO(), filter)

	fmt.Printf("documents upsertCounterMongo with result:%v\n", result)
	if err := result.Decode(doc); err != nil {
		fmt.Println("upsertCounterMongo decode err", err)
		return nil, err
	}

	fmt.Println(doc)
	counter := &model.CounterModel{
		Id:    0,
		Count: int32(doc.Count),
	}
	return counter, nil
}

// getAction 获取action
func getAction(r *http.Request) (string, string, error) {
	decoder := json.NewDecoder(r.Body)
	body := make(map[string]interface{})
	if err := decoder.Decode(&body); err != nil {
		return "", "", err
	}
	defer r.Body.Close()

	fmt.Println("get data from body")
	fmt.Println(body)
	action, ok := body["action"]
	if !ok {
		return "", "", fmt.Errorf("缺少 action 参数")
	}

	dbType, ok := body["dbtype"]
	if !ok {
		return "", "", fmt.Errorf("缺少 dbtype 参数")
	}

	return action.(string), dbType.(string), nil
}

// getIndex 获取主页
func getIndex() (string, error) {
	b, err := ioutil.ReadFile("./index.html")
	if err != nil {
		return "", err
	}
	return string(b), nil
}
