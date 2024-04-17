package response

type MyCode int64

const (
	CodeDefault MyCode = 0

	CodeSuccess       MyCode = 200
	CodeInvalidParams MyCode = 201
	CodeNoData        MyCode = 202
	CodeServerBusy    MyCode = 500

	CodeInvalidToken      MyCode = 301
	CodeInvalidAuthFormat MyCode = 302
	CodeNotLogin          MyCode = 303

	CodeTypeError MyCode = 415

	CodePasswordMistake MyCode = 1001
)

var msgFlags = map[MyCode]string{
	CodeDefault:       "Request successful",
	CodeSuccess:       "success",
	CodeInvalidParams: "Request parameter error",
	CodeServerBusy:    "bustling",
	CodeNoData:        "No data",

	CodeInvalidToken:      "Invalid Token",
	CodeInvalidAuthFormat: "Authentication format is incorrect",
	CodeNotLogin:          "not logged in",

	CodeTypeError: "type error",

	CodePasswordMistake: "incorrect password",
}

func (c MyCode) Msg() string {
	msg, ok := msgFlags[c]
	if ok {
		return msg
	}
	return msgFlags[CodeServerBusy]
}
