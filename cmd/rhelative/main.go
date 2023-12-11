package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/crozzy/rhelative/rhel"
	"github.com/quay/claircore/libvuln/driver"
	"github.com/quay/claircore/libvuln/jsonblob"
	"github.com/quay/claircore/libvuln/updates"
)

func main() {
	ctx := context.Background()
	cl := http.DefaultClient

	out, err := os.Create("test_db.json")

	store, err := jsonblob.New()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func() {
		if err := store.Store(out); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}()
	fac, err := rhel.NewFactory(ctx, rhel.DefaultManifest)
	if err != nil {
		fmt.Println(err)
		return
	}

	mgr, err := updates.NewManager(ctx, store, updates.NewLocalLockSource(), cl,
		updates.WithEnabled([]string{}),
		updates.WithFactories(map[string]driver.UpdaterSetFactory{
			"rhelish": fac,
		}),
	)
	if err = mgr.Run(ctx); err != nil {
		fmt.Println(err)
		return
	}
}
