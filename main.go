package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

const BASE_URL = `https://capi-v2.sankakucomplex.com/posts/keyset?lang=en&limit=100&tags=`

type resBody struct {
	Meta struct {
		Next string `json:"next"`
	} `json:"meta"`
	Data []struct {
		Id       int    `json:"id"`
		FileUrl  string `json:"file_url"`
		FileSize int    `json:"file_size"`
		FileType string `json:"file_type"`
		Md5      string `json:"md5"`
	} `json:"data"`
}

func (v *resBody) json() (s string) {
	if v == nil {
		return
	}
	j, err := json.MarshalIndent(v, "", "    ")
	if err != nil {
		return
	}
	return string(j)
}
func (v *resBody) urls() (s *[]string) {
	if v == nil {
		return
	}
	var u []string
	for _, datum := range v.Data {
		if datum.FileUrl == "" {
			continue
		}
		u = append(u, datum.FileUrl)
	}
	return &u
}

func main() {
	tag := `file_type%3Avideo%20tifa_lockhart%20high_resolution`
	//tag := `file_type%3Avideo%20tracer_%28overwatch%29%20high_resolution`
	//tag := `file_type%3Avideo%20tracer_%28overwatch%29%20very_high_resolution`
	//tag := `file_type%3Avideo%20sombra_%28overwatch%29%20high_resolution`
	//tag := `file_type%3Avideo%20sombra_%28overwatch%29%20very_high_resolution`
	//tag := `file_type%3Avideo%20d.va_%28overwatch%29%20high_resolution`
	//tag := `file_type%3Avideo%20d.va_%28overwatch%29%20very_high_resolution`
	//tag := `file_type%3Avideo%20mei_%28overwatch%29%20high_resolution`
	//tag := `file_type%3Avideo%20mei_%28overwatch%29%20very_high_resolution`
	//tag := `file_type%3Avideo%20mercy_%28overwatch%29%20high_resolution`
	//tag := `file_type%3Avideo%20mercy_%28overwatch%29%20very_high_resolution`
	//tag := `file_type%3Avideo%20widowmaker_%28overwatch%29%20very_high_resolution`
	//tag := `file_type%3Avideo%20widowmaker_%28overwatch%29%20high_resolution`
	//tag := `file_type%3Avideo%20k%2Fda_%28league_of_legends%29%20very_high_resolution`
	//tag := `file_type%3Avideo%20k%2Fda_%28league_of_legends%29%20high_resolution`
	//tag := `file_type%3Avideo%20marie_rose%20very_high_resolution`
	//tag := `file_type%3Avideo%20marie_rose%20high_resolution`
	url := BASE_URL + tag
	method := "GET"

	body := fetch(method, url)
	var res resBody
	json.Unmarshal(body, &res)
	for {
		if res.Meta.Next == "" {
			break
		}
		b := fetch(method, url+`&next=`+res.Meta.Next)
		var r resBody
		json.Unmarshal(b, &r)
		res.Meta.Next = r.Meta.Next
		res.Data = append(res.Data, r.Data...)

	}

	urls := *res.urls()
	fmt.Println("=======================================")
	fmt.Println(strings.Join(urls, "\n"))
	fmt.Println("url count : ", len(urls))
	fmt.Println("=======================================")

}

func fetch(method, url string) (body []byte) {
	status := 0
	defer func() {
		err, _ := recover().(error)
		if err != nil || status != 200 {
			body = nil
			log.Println(err)
		}
	}()

	client := &http.Client{}
	req, _ := http.NewRequest(method, url, nil)

	res, _ := client.Do(req)

	defer res.Body.Close()
	status = res.StatusCode
	body, _ = ioutil.ReadAll(res.Body)
	return
}
