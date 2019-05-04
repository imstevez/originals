package conf

type database struct {
	UserName     string `json:"user_name"`
	Password     string `json:"password"`
	Protocol     string `json:"protocol"`
	Host         string `json:"host"`
	Port         int    `json:"port"`
	Database     string `json:"database"`
	MaxLifeTime  int    `json:"max_life_time"`
	MaxOpenConns int    `json:"max_open_conns"`
	MaxIdleConns int    `json:"max_idle_conns"`
}

type email struct {
	UserName string `json:"user_name"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Sender   string `json:"sender"`
}

type service struct {
	Name    string            `json:"name"`
	Version string            `json:"version"`
	Port    int               `json:"port"`
	Extra   map[string]string `json:"extra"`
}
