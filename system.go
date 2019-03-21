package admin_plugin

import (
	"github.com/ecletus/admin"
	"github.com/ecletus/db"
)

func initResources(e *AdminEvent) {
	e.Admin.NewResource(&admin.AdminSetting{}, &admin.Config{Invisible: true, NotMount: true})
}

func migrate(e *db.DBEvent) error {
	return e.AutoMigrate(&admin.AdminSetting{}).Error
}
