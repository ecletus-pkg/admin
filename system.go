package admin_plugin

import (
	"github.com/aghape/admin"
	"github.com/aghape/db"
)

func initResources(e *AdminEvent) {
	e.Admin.NewResource(&admin.AdminSetting{}, &admin.Config{Invisible: true, NotMount: true})
}

func migrate(e *db.DBEvent) error {
	return e.AutoMigrate(&admin.AdminSetting{}).Error
}
