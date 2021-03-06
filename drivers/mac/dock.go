// +build darwin,amd64

package mac

import (
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/murlokswarm/app"
	"github.com/murlokswarm/app/html"
)

// DockTile implements the app.DockTile interface.
type DockTile struct {
	Menu
}

func newDockTile(c app.MenuConfig) (app.DockTile, error) {
	var markup app.Markup = html.NewMarkup(driver.factory)
	markup = app.ConcurrentMarkup(markup)

	dock := &DockTile{
		Menu: Menu{
			id:             uuid.New().String(),
			markup:         markup,
			lastFocus:      time.Now(),
			keepWhenClosed: true,
		},
	}

	if err := driver.macRPC.Call("menus.New", nil, struct {
		ID string
	}{
		ID: dock.ID(),
	}); err != nil {
		return nil, err
	}

	driver.elems.Put(dock)

	if len(c.DefaultURL) != 0 {
		if err := dock.Load(c.DefaultURL); err != nil {
			return nil, err
		}
	}

	return dock, nil
}

// Load satisfies the app.DockTile interface.
func (d *DockTile) Load(url string, v ...interface{}) error {
	if err := d.Menu.Load(url, v...); err != nil {
		return err
	}

	return driver.macRPC.Call("docks.SetMenu", nil, struct {
		ID string
	}{
		ID: d.ID(),
	})
}

// SetIcon satisfies the app.DockTile interface.
func (d *DockTile) SetIcon(name string) error {
	if _, err := os.Stat(name); err != nil && len(name) != 0 {
		return err
	}

	return driver.macRPC.Call("docks.SetIcon", nil, struct {
		Icon string
	}{
		Icon: name,
	})
}

// SetBadge satisfies the app.DockTile interface.
func (d *DockTile) SetBadge(v interface{}) error {
	var badge string
	if v != nil {
		badge = fmt.Sprint(v)
	}

	return driver.macRPC.Call("docks.SetBadge", nil, struct {
		Badge string
	}{
		Badge: badge,
	})
}

// WhenMenu override the embedded Menu.WhenMenu to do nothing.
// It satisfies the app.ElemWithCompo interface.
func (d *DockTile) WhenMenu(f func(app.Menu)) {
}

// WhenDockTile calls the given handler.
// It satisfies the app.ElemWithCompo interface.
func (d *DockTile) WhenDockTile(f func(app.DockTile)) {
	f(d)
}
