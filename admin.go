package admin_plugin

import (
	"path/filepath"

	"github.com/moisespsena-go/aorm"

	"github.com/ecletus/db"

	"strings"

	"github.com/ecletus/admin"
	"github.com/ecletus/plug"
	"github.com/ecletus/router"
	errwrap "github.com/moisespsena-go/error-wrap"
	path_helpers "github.com/moisespsena-go/path-helpers"
	"github.com/moisespsena-go/xroute"
)

const DEFAULT_ADMIN = "default"

var (
	PKG = path_helpers.GetCalledDir()
)

type AdminGetter struct {
	admins      *Admins
	initialized bool
	dispatcher  plug.PluginEventDispatcherInterface
}

func (this *AdminGetter) GetInitialized(name string) (Admin *admin.Admin) {
	this.Initialize()
	return this.admins.ByName[name]
}

func (this *AdminGetter) Initialize() {
	if !this.initialized {
		this.initialized = true
		err := this.admins.Trigger(this.dispatcher)
		if err != nil {
			panic(errwrap.Wrap(err, "Trigger Admins [%s]", strings.Join(this.admins.Names(), ", ")))
		}
	}
}

type Plugin struct {
	plug.EventDispatcher
	db.DBNames

	AdminsKey, AdminGetterKey string
	SystemDBDialectKey        string
	initialized               bool
}

func (p *Plugin) ProvideOptions() []string {
	return []string{p.AdminGetterKey}
}

func (p *Plugin) RequireOptions() []string {
	return []string{p.AdminsKey, p.SystemDBDialectKey}
}

func (p *Plugin) NameSpace() string {
	return filepath.Join("github.com", "ecletus", "admin")
}

func (p *Plugin) AssetsRootPath() (pth string) {
	pth = filepath.Join(path_helpers.GetCalledDir(true), "..", "..", "ecletus", "admin")
	return
}

func (p *Plugin) ProvidesOptions(options *plug.Options) {
	dis := plug.Dis(options)
	admins := options.GetInterface(p.AdminsKey).(*Admins)
	getter := &AdminGetter{
		admins:     admins,
		dispatcher: dis,
	}
	options.Set(p.AdminGetterKey, getter)
}

func (p *Plugin) Init(options *plug.Options) {
	options.GetInterface(p.AdminGetterKey).(*AdminGetter).Initialize()
}

func (p *Plugin) OnRegister(options *plug.Options) {
	plug.OnFS(p, func(e *plug.FSEvent) {
		e.RegisterAssetPath(e.PathOf(&admin.Admin{}))
	})

	_ = plug.OnPostInit(p, func(e plug.PluginEventInterface) {
		Events(p).InitResources(initResources)
		db.Events(p).DBOnMigrate(migrate)
	})

	adminsCalled := map[string]bool{}

	log := p.Logger()

	p.On(E_ADMIN, func(e plug.PluginEventInterface) (err error) {
		adminEvent := e.(*AdminEvent)
		if _, ok := adminsCalled[adminEvent.AdminName]; ok {
			return nil
		}
		adminsCalled[adminEvent.AdminName] = true
		adminName, Admin := adminEvent.AdminName, adminEvent.Admin

		if systemDialect := options.GetString(p.SystemDBDialectKey, Admin.Config.FakeDBDialect); systemDialect != Admin.Config.FakeDBDialect {
			Admin.FakeDB = aorm.FakeDB(systemDialect)
		}
		// github.com/ecletus-pkg/admin
		//
		log.Debugf("trigger AdminEvent")
		if err = e.PluginDispatcher().TriggerPlugins(&AdminEvent{plug.NewPluginEvent(EAdmin(adminName)), Admin, adminName, e}); err != nil {
			return errwrap.Wrap(err, "AdminEvent")
		}
		log.Debugf("trigger AdminInitResourcesEvent")
		if err = e.PluginDispatcher().TriggerPlugins(&AdminEvent{plug.NewPluginEvent(EInitResources(adminName)), Admin, adminName, e}); err != nil {
			return errwrap.Wrap(err, "AdminInitResourcesEvent")
		}
		log.Debugf("trigger AdminFuncMapEvent")
		if err = e.PluginDispatcher().TriggerPlugins(&AdminFuncMapEvent{&AdminEvent{plug.NewPluginEvent(EFuncMap(adminName)), Admin, adminName, e}}); err != nil {
			return errwrap.Wrap(err, "AdminFuncMapEvent")
		}
		log.Debugf("trigger AdminDone")
		if err = Admin.TriggerDone(&admin.AdminEvent{plug.NewEvent(admin.E_DONE), Admin}); err != nil {
			return errwrap.Wrap(err, "Admin.Done")
		}
		if err = e.PluginDispatcher().TriggerPlugins(&AdminEvent{plug.NewPluginEvent(EDone(adminName)), Admin, adminName, e}); err != nil {
			return errwrap.Wrap(err, "AdminDone")
		}
		return nil
	})

	err := router.OnRouteE(p, func(e *router.RouterEvent) (err error) {
		getter := options.GetInterface(p.AdminGetterKey).(*AdminGetter)
		admins := getter.admins
		getter.Initialize()
		return admins.Each(func(adminName string, Admin *admin.Admin) error {
			log.Debugf("[admin=%q] mounted on %v", adminName, Admin.Config.MountPath)
			var router xroute.Router = e.Router.Mux
			if Admin.Config.MountPath != "/" {
				router = xroute.NewMux(Admin.Name)
				Admin.InitRoutes(router)
				e.Router.Mux.Mount(Admin.Config.MountPath, router)
			} else {
				Admin.InitRoutes(router)
			}
			return errwrap.Wrap(e.PluginDispatcher().TriggerPlugins(&AdminRouterEvent{&AdminEvent{
				plug.NewPluginEvent(ERoute(adminName)), Admin,
				adminName, e}, router}),
				"AdminRouterEvent")
		})
	})
	if err != nil {
		panic(err)
	}
}
