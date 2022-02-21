package model

type Message struct {
	SecretId   string `json:"secretId"`
	SecretKey  string `json:"secretKey"`
	AppId      string `json:"appId"`
	AppKey     string `json:"appKey"`
	SignId     string `json:"signId"`
	TemplateId string `json:"templateId"`
	Sign       string `json:"sign"`
}
