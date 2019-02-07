package enc

import "encoding/json"

func HTTPResponseJSON(v interface{}, errs ...error) (contentType string, body []byte, err error) {
	if len(errs) > 0 && errs[0] != nil {
		err = errs[0]
		return
	}
	if err, _ = v.(error); err != nil {
		return
	}
	contentType = "text/json; charset=utf-8"
	body, err = json.Marshal(v)
	return
}
