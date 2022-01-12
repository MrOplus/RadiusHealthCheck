package models

type BasicCredentials struct {
	Username string `json:"Username"`
	Password string `json:"Password"`
}

func IsBasicCredentialsEmpty(basicCredentials BasicCredentials) bool {
	return basicCredentials.Username == "" && basicCredentials.Password == ""
}
