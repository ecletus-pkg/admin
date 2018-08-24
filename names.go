package admin_plugin

type AdminNames struct {
	Names []string
}

func (dn *AdminNames) EachE(cb func(adminName string) error) (err error) {
	for _, name := range dn.Names {
		err = cb(name)
		if err != nil {
			return
		}
	}
	return
}

func (dn *AdminNames) Each(cb func(adminName string)) {
	dn.EachE(func(adminName string) error {
		cb(adminName)
		return nil
	})
}

func (dn *AdminNames) EachOrDefaultE(cb func(adminName string) error) (err error) {
	if len(dn.Names) == 0 {
		dn.Names = []string{DEFAULT_ADMIN}
	}
	return dn.EachE(cb)
}

func (dn *AdminNames) EachOrDefault(cb func(adminName string)) {
	dn.EachOrDefaultE(func(adminName string) error {
		cb(adminName)
		return nil
	})
}

func (a *AdminNames) GetNames() []string {
	if len(a.Names) == 0 {
		a.Names = []string{DEFAULT_ADMIN}
	}
	return a.Names
}
