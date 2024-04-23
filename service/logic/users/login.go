package users

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"simple-video-net/consts"
	"simple-video-net/global"
	receive "simple-video-net/interaction/receive/users"
	response "simple-video-net/interaction/response/users"
	"simple-video-net/models/common"
	userModel "simple-video-net/models/users"
	"simple-video-net/utils/conversion"
	"simple-video-net/utils/email"
	"simple-video-net/utils/jwt"
	"time"

	"github.com/go-redis/redis"
)

func WxAuthorization(data *receive.WxAuthorizationReceiveStruct) (results interface{}, err error) {

	type WXLoginResp struct {
		OpenID     string `json:"openid"`
		SessionKey string `json:"session_key"`
		UnionID    string `json:"unionId"`
		ErrCode    int    `json:"errCode"`
		ErrMsg     string `json:"errMsg"`
		Token      string `json:"token"`
		Auth       uint32 `json:"auth"`
		Name       string `json:"name"`
		Phone      string `json:"phone"`
		Nickname   string `json:"nickname"`
		HeadImage  string `json:"headImage"`
	}
	url := "https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code"
	url = fmt.Sprintf(url, "wxfbd9d7966fc9796c", "92fbe8e2921e00fc3ba68e34d5d0b986", data.Code)
	resp, err := http.Get(url)
	if err != nil {
		//response.ResponseError(ctx, err.Error())
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)
	// Parses the body data from the http request into the structure we defined.
	wxResp := WXLoginResp{}
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&wxResp); err != nil {
		//response.ResponseError(ctx, err.Error())
		return nil, err
	}
	//Get openid for processing
	users := new(userModel.User)
	if !users.IsExistByField("openid", wxResp.OpenID) {
		//When this microsoft is not registered
		photo, _ := json.Marshal(common.Img{
			Src: data.AvatarUrl,
			Tp:  "wx",
		})
		users := userModel.User{
			Username: data.NickName,
			Openid:   wxResp.OpenID,
			Photo:    photo,
		}
		registerRes := users.Create()
		if !registerRes {
			return nil, fmt.Errorf("registration failure")
		}
		//Registering a token
		tokenString := jwt.NextToken(users.ID)
		src, _ := conversion.FormattingJsonSrc(users.Photo)
		userInfo := response.UserInfoResponseStruct{
			ID:       users.ID,
			UserName: users.Username,
			Photo:    src,
			Token:    tokenString,
		}
		return userInfo, nil
	}
	//Returns a token if you have already registered
	fmt.Printf("The queried user id is:%v", users.ID)
	src, _ := conversion.FormattingJsonSrc(users.Photo)
	tokenString := jwt.NextToken(users.ID)
	userInfo := response.UserInfoResponseStruct{
		ID:       users.ID,
		UserName: users.Username,
		Photo:    src,
		Token:    tokenString,
	}
	return userInfo, nil
}

func Register(data *receive.RegisterReceiveStruct) (results interface{}, err error) {
	//Determine if a mailbox is unique
	users := new(userModel.User)
	if users.IsExistByField("email", data.Email) {
		return nil, fmt.Errorf("Email already registered")
	}
	//Determine if the CAPTCHA is correct
	verCode, err := global.RedisDb.Get(fmt.Sprintf("%s%s", consts.RegEmailVerCode, data.Email)).Result()
	if err == redis.Nil {
		return nil, fmt.Errorf("Captcha expired!")
	}

	if verCode != data.VerificationCode {
		return nil, fmt.Errorf("CAPTCHA error")
	}
	//Generate password salt 8 bits
	salt := make([]byte, 6)
	for i := range salt {
		salt[i] = jwt.SaltStr[rand.Int63()%int64(len(jwt.SaltStr))]
	}
	password := []byte(fmt.Sprintf("%s%s%s", salt, data.Password, salt))
	passwordMd5 := fmt.Sprintf("%x", md5.Sum(password))
	photo, _ := json.Marshal(common.Img{
		Src: "",
		Tp:  "local",
	})
	registerData := &userModel.User{
		Email:     data.Email,
		Username:  data.UserName,
		Salt:      string(salt),
		Password:  passwordMd5,
		Photo:     photo,
		BirthDate: time.Now(),
	}
	registerRes := registerData.Create()
	if !registerRes {
		return nil, fmt.Errorf("registration failure")
	}
	//Registering a token
	tokenString := jwt.NextToken(registerData.ID)
	results = response.UserInfoResponse(registerData, tokenString)

	return results, nil

}

func Login(data *receive.LoginReceiveStruct) (results interface{}, err error) {
	users := new(userModel.User)
	if !users.IsExistByField("username", data.Username) {
		return nil, fmt.Errorf("Account does not exist")
	}
	if !users.IfPasswordCorrect(data.Password) {
		return nil, fmt.Errorf("incorrect password")
	}
	//Registering a token
	tokenString := jwt.NextToken(users.ID)
	userInfo := response.UserInfoResponse(users, tokenString)
	return userInfo, nil
}

func SendEmailVerCode(data *receive.SendEmailVerCodeReceiveStruct) (results interface{}, err error) {
	users := new(userModel.User)
	if users.IsExistByField("email", data.Email) {
		return nil, fmt.Errorf("Email already registered")
	}
	//sender
	mailTo := []string{data.Email}
	// Email Subject
	code := fmt.Sprintf("%06v", rand.New(rand.NewSource(time.Now().UnixNano())).Int63n(1000000))
	subject := "CAPTCHA"
	// Body of the email
	body := fmt.Sprintf("You are registering with the verification code:%s,5 minutes. Please do not forward to others.", code)
	err = email.SendMail(mailTo, subject, body)
	if err != nil {
		global.Logger.Error("send to%dEmail verification code failed", data.Email)
		return nil, fmt.Errorf("Failed to send")
	}
	err = global.RedisDb.Set(fmt.Sprintf("%s%s", consts.RegEmailVerCode, data.Email), code, 5*time.Minute).Err()
	if err != nil {
		return nil, err
	}
	return "Sent successfully", nil
}

func SendEmailVerCodeByForget(data *receive.SendEmailVerCodeReceiveStruct) (results interface{}, err error) {
	//Determine if a user exists
	users := new(userModel.User)
	if !users.IsExistByField("email", data.Email) {
		return nil, fmt.Errorf("This email is not registered")
	}
	//sender
	mailTo := []string{data.Email}
	// Email Subject
	code := fmt.Sprintf("%06v", rand.New(rand.NewSource(time.Now().UnixNano())).Int63n(1000000))
	subject := "CAPTCHA"
	// email body
	body := fmt.Sprintf("You are retrieving your passwordYour verification code is:%s,5 minutes. Please do not forward to others.", code)
	err = email.SendMail(mailTo, subject, body)
	if err != nil {
		return nil, err
	}
	err = global.RedisDb.Set(fmt.Sprintf("%s%s", consts.RegEmailVerCodeByForget, data.Email), code, 5*time.Minute).Err()
	if err != nil {
		return nil, err
	}
	return "Sent successfully", nil
}

func Forget(data *receive.ForgetReceiveStruct) (results interface{}, err error) {
	//Determine if a mailbox is unique
	users := new(userModel.User)
	if !users.IsExistByField("email", data.Email) {
		return nil, fmt.Errorf("This account does not exist")
	}
	//Determine if the CAPTCHA is correct
	verCode, err := global.RedisDb.Get(fmt.Sprintf("%s%s", consts.RegEmailVerCodeByForget, data.Email)).Result()
	if err == redis.Nil {
		return nil, fmt.Errorf("Captcha expired!")
	}

	if verCode != data.VerificationCode {
		return nil, fmt.Errorf("CAPTCHA error")
	}
	//Generate password salt 8 bits
	salt := make([]byte, 6)
	for i := range salt {
		salt[i] = jwt.SaltStr[rand.Int63()%int64(len(jwt.SaltStr))]
	}
	password := []byte(fmt.Sprintf("%s%s%s", salt, data.Password, salt))
	passwordMd5 := fmt.Sprintf("%x", md5.Sum(password))

	registerData := userModel.User{
		Salt:     string(salt),
		Password: passwordMd5,
	}
	registerRes := registerData.Update()
	if !registerRes {
		return nil, fmt.Errorf("Modification Failure")
	}
	return "Modified successfully", nil
}
