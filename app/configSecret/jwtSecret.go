package configSecret

const jwtKey = "jwtKey"

func SetJwtKey(value string) error {
	return setConfig(jwtKey, value)
}

func GetJwtKey() (string, error) {
	return getConfig(jwtKey)
}

func IsSetJwtKey() bool {
	return checkConfig(jwtKey)
}

func DelJwtKey() error {
	return delConfig(jwtKey)
}
