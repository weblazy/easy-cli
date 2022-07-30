package template

import (
	"encoding/json"

	"github.com/google/uuid"
	"github.com/sunmi-OS/gocore/v2/tools/gocore/conf"
	"github.com/sunmi-OS/gocore/v2/tools/gocore/file"
)

func createPostman(root, name string, commonHeaders []conf.Header, apis []conf.Api) {
	result := map[string]interface{}{
		"info": map[string]interface{}{
			"_postman_id": uuid.New().String(),
			"name":        name,
			"schema":      "https://schema.getpostman.com/json/collection/v2.1.0/collection.json",
		},
		"response": []interface{}{},
	}
	headers := make([]map[string]interface{}, 0)
	for k1 := range commonHeaders {
		v1 := commonHeaders[k1]
		headers = append(headers, map[string]interface{}{
			"key":   v1.Key,
			"value": v1.Value,
			"type":  "text",
		})
	}

	items := make([]map[string]interface{}, 0)

	for k1 := range apis {
		v1 := apis[k1]
		for k2 := range v1.Handle {
			v2 := v1.Handle[k2]

			raw := map[string]string{}
			for k3 := range v2.RequestParams {
				v3 := v2.RequestParams[k3]
				raw[v3.Name] = ""
			}
			rawJson, _ := json.Marshal(raw)
			request := map[string]interface{}{
				"method": v2.Method,
				"header": headers,
				"body": map[string]interface{}{
					"mode": "raw",
					"raw":  string(rawJson),
				},
				"url": map[string]interface{}{
					"path":     []string{file.CamelToUnderline(v1.ModuleName), v1.Prefix, file.CamelToUnderline(v2.Name)},
					"protocol": "https",
					"host": []string{
						"{{" + name + "}}",
					},
				},
			}
			item := map[string]interface{}{
				"name":    file.CamelToUnderline(v2.Name),
				"request": request,
			}
			items = append(items, item)
		}

	}
	result["item"] = items
	resultByte, _ := json.Marshal(result)
	fileBuffer.Write(resultByte)
	fileForceWriter(fileBuffer, root+"/postman.json")
}
