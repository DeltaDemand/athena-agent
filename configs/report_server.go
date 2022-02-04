package appConfigs

type ReportServer struct {
	Ip   string `json:"ip"`
	Port string `json:"port"`
}

func (r ReportServer) GetAddress() string {
	return r.Ip + ":" + r.Port
}
