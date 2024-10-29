package esClose

import "github.com/elastic/go-elasticsearch/v8/esapi"

func CloseResp(resp *esapi.Response) {
	if resp == nil || resp.Body == nil {
		return
	}
	_ = resp.Body.Close()
}
