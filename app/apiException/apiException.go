package apiException

import "net/http"

type Error struct {
	StatusCode int    `json:"-"`
	Code       int    `json:"code"`
	Msg        string `json:"msg"`
}

var (
	ServerError   = NewError(http.StatusInternalServerError, 200500, "系统异常，请稍后重试!")
	ParamError    = NewError(http.StatusInternalServerError, 200501, "参数错误")
	EmailExist    = NewError(http.StatusInternalServerError, 200502, "邮箱已存在")
	EmailValid    = NewError(http.StatusInternalServerError, 200503, "邮箱格式不正确")
	PhoneExist    = NewError(http.StatusInternalServerError, 200504, "电话已存在")
	PhoneValid    = NewError(http.StatusInternalServerError, 200505, "电话格式不正确")
	NicknameExist = NewError(http.StatusInternalServerError, 200506, "昵称已存在")
	NicknameValid = NewError(http.StatusInternalServerError, 200507, "昵称中不得包含特殊符号")
	PasswordValid = NewError(http.StatusInternalServerError, 200508, "密码必须包含字母，数字，特殊符号且长度在8位以上")
	UserExist     = NewError(http.StatusInternalServerError, 200509, "用户不存在")
	PasswordWrong = NewError(http.StatusInternalServerError, 200510, "密码不正确")
	CodeWrong     = NewError(http.StatusInternalServerError, 200511, "验证码不正确")
	PictureError  = NewError(http.StatusInternalServerError, 200512, "仅允许上传图片文件")
	NotInit       = NewError(http.StatusNotFound, 200404, http.StatusText(http.StatusNotFound))
	NotFound      = NewError(http.StatusNotFound, 200404, http.StatusText(http.StatusNotFound))
	Unknown       = NewError(http.StatusInternalServerError, 300500, "系统异常，请稍后重试!")
)

func OtherError(message string) *Error {
	return NewError(http.StatusForbidden, 100403, message)
}

func (e *Error) Error() string {
	return e.Msg
}

func NewError(statusCode, Code int, msg string) *Error {
	return &Error{
		StatusCode: statusCode,
		Code:       Code,
		Msg:        msg,
	}
}
