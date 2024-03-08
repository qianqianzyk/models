package configSecret

const aesKey = "aesKey"

func SetAesKey(value string) error {
	return setConfig(aesKey, value)
}

func GetAesKey() (string, error) {
	return getConfig(aesKey)
}

func IsSetAesKey() bool {
	return checkConfig(aesKey)
}

func DelAesKey() error {
	return delConfig(aesKey)
}
