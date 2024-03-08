package configSecret

const urlKey = "urlKey"

func SetUrlKey(value string) error {
	return setConfig(urlKey, value)
}

func GetUrlKey() (string, error) {
	return getConfig(urlKey)
}

func IsSetUrlKey() bool {
	return checkConfig(urlKey)
}

func DelUrlKey() error {
	return delConfig(urlKey)
}
