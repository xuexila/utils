package httpServer

import (
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

func GetLoginInfo(session string) (LoginInfo, bool) {
	tmp, ok := LoginSessionMap.Load(session)
	if !ok {
		return LoginInfo{}, false
	}
	return tmp.(LoginInfo), true
}
