package index

import (
	"github.com/julienschmidt/httprouter"
	"github.com/raozhaofeng/zfeng/utils/body"
	"net/http"
)

func Test(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	body.SuccessJSON(w, "Hello World")
}
