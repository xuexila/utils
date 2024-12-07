package httpServer

import (
	"encoding/json"
	"net/http"
)

func (this Router) GetLoginInfo(w http.ResponseWriter, r *http.Request) (LoginInfo, error) {
	info := LoginInfo{}
	raw, err := this.Store.GetUp(w, r, this.SessionLoginName)
	if err != nil {
		return info, err
	}
	if err = json.Unmarshal([]byte(raw), &info); err != nil {
		return info, err
	}
	return info, nil
}
