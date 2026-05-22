// Package watcher provides file-system polling for Terraform state files.
//
// It computes SHA-256 hashes of watched files at a configurable interval and
// emits an Event on the Changes channel whenever a file's content changes.
//
// Typical usage:
//
//	w := watcher.New([]string{"base.tfstate", "target.tfstate"}, 2*time.Second)
//	go func() {
//		for ev := range w.Changes {
//			fmt.Printf("changed: %s\n", ev.Path)
//		}
//	}()
//	if err := w.Start(ctx); err != nil && !errors.Is(err, context.Canceled) {
//		log.Fatal(err)
//	}
package watcher
