package service

import (
	"database/sql"
	"fmt"
	"gin/dao"
	"gin/model"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	tencentsms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
	"math/rand"
	"time"
)

// Register 注册
func Register(mobile string) error {
	err := dao.Register(mobile)
	if err != nil {
		return err
	}
	return nil
}

// IsRegister 查询是否注册
func IsRegister(mobile string) (bool, error) {
	err := dao.IsRegister(mobile)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil //未被注册返回false
		}
		return true, err
	}
	return true, nil
}

// SendSms 发送短信
func SendSms(mobile string) error {
	//生成验证码
	code := fmt.Sprintf("%06v", rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(1000000))
	//判读该电话号码是否被注册
	err := dao.IsRegister(mobile)
	if err != nil {
		if err != sql.ErrNoRows {
			return err
		}
		//未被注册
		err := dao.Register(mobile)
		if err != nil {
			return err
		}
	}
	//储存验证码
	sendTime := time.Now()
	user := model.User{
		Mobile:     mobile,
		VerifyCode: code,
		SendTime:   sendTime,
	}
	err = dao.StoreVerifyCode(user)
	if err != nil {
		return err
	}
	//发送短信
	message := model.Message{
		SecretId:   "123",
		SecretKey:  "123",
		AppId:      "123",
		AppKey:     "123",
		SignId:     "123",
		TemplateId: "123",
		Sign:       "123",
	}

	credential := common.NewCredential(
		message.SecretId,
		message.SecretKey,
	)
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "sms.tencentcloudapi.com"
	client, err := tencentsms.NewClient(credential, "ap-guangzhou", cpf)
	if err != nil {
		return err
	}

	request := tencentsms.NewSendSmsRequest()
	request.SmsSdkAppId = common.StringPtr(message.AppId)
	request.SignName = common.StringPtr(message.Sign)
	request.SenderId = common.StringPtr("")
	request.ExtendCode = common.StringPtr("")
	request.TemplateParamSet = common.StringPtrs([]string{code})
	request.TemplateId = common.StringPtr(message.TemplateId)
	request.PhoneNumberSet = common.StringPtrs([]string{"+86" + mobile})

	_, err = client.SendSms(request)
	if err != nil {
		return err
	}
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		return err
	}
	return nil
}

// SelectVerifyCodeAndSendTime 查询验证码以及发送时间
func SelectVerifyCodeAndSendTime(mobile string) (string, time.Time, error) {
	code, time, err := dao.SelectVerifyCodeAndSendTime(mobile)
	if err != nil {
		return "", time, err
	}
	return code, time, nil
}

// ImprovePersonalInfo 完善个人信息
func ImprovePersonalInfo(user model.User) error {
	err := dao.ImprovePersonalInfo(user)
	if err != nil {
		return err
	}
	return nil
}

// SelectUsernameByMobile 查询用户名
func SelectUsernameByMobile(mobile string) (error, bool) {
	err := dao.SelectUsernameByMobile(mobile)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, false
		}
		return err, true
	}
	return nil, true
}

// SelectInfoByMobile 查询个人信息
func SelectInfoByMobile(mobile string) (error, model.User) {
	err, user := dao.SelectInfoByMobile(mobile)
	if err != nil {
		return err, user
	}
	return nil, user
}

// SelectPasswordByMobile 查询密码
func SelectPasswordByMobile(user model.User) (error, string) {
	err, pwd := dao.SelectPasswordByMobile(user)
	if err != nil {
		return err, pwd
	}
	return err, pwd
}

// ChangePassword 修改密码
func ChangePassword(user model.User) error {
	err := dao.ChangePassword(user)
	if err != nil {
		return err
	}
	return nil
}
