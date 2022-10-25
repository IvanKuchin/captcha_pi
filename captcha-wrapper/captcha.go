package captcha_wrapper

import (
	"net/http"

	"github.com/dchest/captcha"
)

func MakeNew(tID string) (string, error) {
	return captcha.NewLen(4), nil
}

func CaptchaHandler() http.Handler {
	return captcha.Server(captcha.StdWidth*2/3, captcha.StdHeight)
}

func Verify(captcha_id, solution string) (bool, error) {
	return captcha.VerifyString(captcha_id, solution), nil
}
