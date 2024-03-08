package configSecret

const emailSendKey = "emailSendKey"

func SetEmailSendKey(value string) error {
	return setConfig(emailSendKey, value)
}

func GetEmailSendKey() (string, error) {
	return getConfig(emailSendKey)
}

func IsSetEmailSendKey() bool {
	return checkConfig(emailSendKey)
}

func DelEmailSendKey() error {
	return delConfig(emailSendKey)
}
