package ginhelper

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"mime/multipart"
	"regexp"
	"strings"
)

func GetRequestBody(c *gin.Context) (map[string]interface{}, error) {
	o := make(map[string]interface{})

	form, _ := c.MultipartForm()

	if form != nil {
		// Process Form Data
		parseValues(form.Value, o)
		parseFiles(form.File, o)
	} else {
		x, _ := ioutil.ReadAll(c.Request.Body)

		err := json.Unmarshal(x, &o)

		if err != nil {
			return nil, err
		}

		json.Unmarshal(x, &o)
		o["_json"] = []byte(x)
	}

	return o, nil
}

func parseValues(values map[string][]string, o map[string]interface{}) {
	for k, v := range values {
		parseKeyValue(k, v, o)
	}
}

func parseFiles(values map[string][]*multipart.FileHeader, o map[string]interface{}) {
	for k, v := range values {
		parseKeyValue(k, v, o)
	}
}

func parseKeyValue(k string, v interface{}, o map[string]interface{}) {
	if k == "_json" {
		jsonStr := &v.([]string)[0]
		json.Unmarshal([]byte(*jsonStr), &o)
		o["_json"] = []byte(*jsonStr)
		return
	}

	re := regexp.MustCompile("^.*\\[.*\\]$")
	if re.MatchString(k) {
		i := strings.Index(k, "[")
		lk := k[:i]
		rk := k[i+1 : len(k)-1]

		if o[lk] == nil {
			o[lk] = make(map[string]interface{})
		}

		switch v.(type) {
		case []string:
			if len(v.([]string)) > 1 {
				o[lk].(map[string]interface{})[rk] = v
			} else {
				o[lk].(map[string]interface{})[rk] = v.([]string)[0]
			}
		case []*multipart.FileHeader:
			if len(v.([]*multipart.FileHeader)) > 1 {
				o[lk].(map[string]interface{})[rk] = v
			} else {
				o[lk].(map[string]interface{})[rk] = v.([]*multipart.FileHeader)[0]
			}
		}
	} else {
		switch v.(type) {
		case []string:
			if len(v.([]string)) > 1 {
				o[k] = v
			} else {
				o[k] = v.([]string)[0]
			}
		case []*multipart.FileHeader:
			if len(v.([]*multipart.FileHeader)) > 1 {
				o[k] = v
			} else {
				o[k] = v.([]*multipart.FileHeader)[0]
			}
		}
	}
}
