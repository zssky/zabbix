package zabbix

// Template - zabbix template interface
type Template struct {
	ID          string `json:"templateid"`
	Host        string `json:"host"`
	Description string `json:"description"`
	Name        string `json:"name"`
}

// Templates - the array of template
type Templates []Template
