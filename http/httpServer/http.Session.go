package httpServer

import (
	"encoding/json"
	uuid "github.com/satori/go.uuid"
	"net/http"
	"strings"
)

func GetSessionId(r *http.Request, sid string) (string, error) {
	cookie, err := r.Cookie(sid)
	if err != nil {

		return strings.ReplaceAll(uuid.NewV4().String(), "-", ""), err
	}
	return cookie.Value, nil
}

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
