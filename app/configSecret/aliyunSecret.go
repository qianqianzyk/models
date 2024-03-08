package configSecret

const accessKey = "accessKey"
const accessSecret = "accessSecret"
const templateCode = "templateCode"

func SetAccessKey(value string) error {
	return setConfig(accessKey, value)
}

func GetAccessKey() (string, error) {
	return getConfig(accessKey)
}

func IsSetAccessKey() bool {
	return checkConfig(accessKey)
}

func DelAccessKey() error {
	return delConfig(accessKey)
}

func SetAccessSecret(value string) error {
	return setConfig(accessSecret, value)
}

func GetAccessSecret() (string, error) {
	return getConfig(accessSecret)
}

func IsSetAccessSecret() bool {
	return checkConfig(accessSecret)
}

func DelAccessSecret() error {
	return delConfig(accessSecret)
}

func SetTemplateCode(value string) error {
	return setConfig(templateCode, value)
}

func GetTemplateCode() (string, error) {
	return getConfig(templateCode)
}

func IsSetTemplateCode() bool {
	return checkConfig(templateCode)
}

func DelTemplateCode() error {
	return delConfig(templateCode)
}
