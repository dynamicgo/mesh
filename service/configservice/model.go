package configservice

import "github.com/dynamicgo/orm"

// ServiceConfig .
type ServiceConfig struct {
	ID      string `xorm:"pk"`
	Path    string `xorm:"index"`
	Content string `xorm:""`
}

// TableName .
func (table *ServiceConfig) TableName() string {
	return "mesh_serviceconfig"
}

func init() {
	orm.RegisterWithName("configservice", func() []interface{} {
		return []interface{}{
			&ServiceConfig{},
		}
	})
}
