package userService

import (
	"context"
	"errors"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/gomail.v2"
	"live-chat/app/configSecret"
	"live-chat/app/models"
	"live-chat/app/utils"
	"live-chat/config/database"
	"live-chat/config/redis"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"
)

var lastEmailTime time.Time

func IsPhoneExist(phone string) error {
	var info models.User
	phone = utils.AesEncrypt(phone)
	result := database.DB.Where(models.User{
		Phone: phone,
	}).First(&info)
	return result.Error
}

func IsEmailExist(email string) error {
	var info models.User
	email = utils.AesEncrypt(email)
	result := database.DB.Where(models.User{
		Email: email,
	}).First(&info)
	return result.Error
}

func IsNicknameExist(nickname string) error {
	var info models.User
	result := database.DB.Where(models.User{
		Nickname: nickname,
	}).First(&info)
	return result.Error
}

func IsPasswordValid(password string) bool {
	hasLetter := false
	hasNumber := false
	hasSpecialChar := false
	specialChars := "@$!%*?&"
	for _, char := range password {
		if unicode.IsDigit(char) {
			hasNumber = true
		} else if unicode.IsLetter(char) {
			hasLetter = true
		} else if strings.ContainsRune(specialChars, char) {
			hasSpecialChar = true
		}
		if hasLetter && hasNumber && hasSpecialChar {
			return true
		}
	}
	return false
}

func CreateUser(info models.User) error {
	info.Email = utils.AesEncrypt(info.Email)
	info.Phone = utils.AesEncrypt(info.Phone)
	info.Name = utils.AesEncrypt(info.Name)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(info.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	info.Password = string(hashedPassword)
	result := database.DB.Model(models.User{}).Create(&info)
	return result.Error
}

func IsUserExist(account string) (*models.User, error) {
	var field string
	if match, _ := regexp.MatchString(`^\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*$`, account); match {
		field = "email"
		account = utils.AesEncrypt(account)
	} else {
		field = "nickname"
	}
	var user models.User
	if err := database.DB.Model(&models.User{}).Where(field+" = ?", account).First(&user).Error; err != nil {
		return nil, err
	}
	aesDecryptUserInfo(&user)
	return &user, nil
}

func GetUserByUserID(userId int) (*models.User, error) {
	var user models.User
	if err := database.DB.Model(&models.User{}).Where(models.User{
		ID: userId,
	}).First(&user).Error; err != nil {
		return nil, err
	}
	aesDecryptUserInfo(&user)
	return &user, nil
}

func UpdateUserInfoByUserID(id int, user models.User) error {
	user.Email = utils.AesEncrypt(user.Email)
	user.Phone = utils.AesEncrypt(user.Phone)
	user.Name = utils.AesEncrypt(user.Name)
	if err := database.DB.Model(models.User{}).Select("*").Omit("id", "type", "password", "create_time").
		Where(&models.User{ID: id}).Updates(&user).Error; err != nil {
		return err
	}
	return nil
}

func aesDecryptUserInfo(user *models.User) {
	user.Email = utils.AesDecrypt(user.Email)
	user.Phone = utils.AesDecrypt(user.Phone)
	user.Name = utils.AesDecrypt(user.Name)
}

func storeVerificationCode(key string, value int) error {
	redis.RedisClient.Set(context.TODO(), "user:"+key, value, 10*time.Minute)
	return nil
}

func GetVerificationCode(key string) (int, error) {
	val, err := redis.RedisClient.Get(context.Background(), "user:"+key).Result()
	code, err := strconv.Atoi(val)
	if err != nil {
		return 0, err
	}
	return code, nil
}

func sendVerificationCode(email string) error {
	key, err := configSecret.GetEmailSendKey()
	if err != nil {
		return err
	}
	m := gomail.NewMessage()
	m.SetHeader("From", "qianqianzyk@qq.com")
	m.SetHeader("To", email)
	m.SetHeader("Subject", "[邮箱验证]")
	code := rand.Intn(900000) + 100000
	body := "验证码：" + strconv.Itoa(code) + "。有效期10分钟，请勿泄露。"
	m.SetBody("text/plain", body)
	d := gomail.NewDialer("smtp.qq.com", 465, "qianqianzyk@qq.com", key)
	if err := storeVerificationCode(email, code); err != nil {
		return err
	}
	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}

func SendMsgToEmail(email string) error {
	elapsedTime := time.Since(lastEmailTime)
	if elapsedTime < 60*time.Second {
		return errors.New("请等待60秒后再尝试发送邮件")
	}
	err := sendVerificationCode(email)
	if err != nil {
		return err
	}
	lastEmailTime = time.Now()
	return nil
}

func VerifyEmail(id int) error {
	result := database.DB.Model(models.User{}).Where(&models.User{ID: id}).Update("email_type", 2)
	return result.Error
}

func FindPassword(id int, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	result := database.DB.Model(models.User{}).Where(&models.User{ID: id}).Update("password", hashedPassword)
	return result.Error
}

func SendMsgToPhone(phone string) error {
	code := rand.Intn(900000) + 100000
	accessKey, err := configSecret.GetAccessKey()
	accessSecret, err := configSecret.GetAccessSecret()
	templateCode, err := configSecret.GetTemplateCode()
	if err != nil {
		return err
	}
	client, err := dysmsapi.NewClientWithAccessKey("cn-hangzhou", accessKey, accessSecret)
	request := dysmsapi.CreateSendSmsRequest()
	request.Scheme = "https"
	request.PhoneNumbers = phone
	request.SignName = "qianqianzyk"
	request.TemplateCode = templateCode
	request.TemplateParam = "{\"code\":\"" + strconv.Itoa(code) + "\"}"
	response, err := client.SendSms(request)
	if response.Code == "isv.BUSINESS_LIMIT_CONTROL" {
		return errors.New("frequency_limit")
	}
	if err != nil {
		fmt.Print(err.Error())
		return errors.New("failed")
	}
	if err := storeVerificationCode(phone, code); err != nil {
		return err
	}
	return nil
}
