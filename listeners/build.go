package listeners

import (
	"github.com/keyvie/modmail/internal/dis"
	"github.com/keyvie/modmail/internal/state"
)

func Build() map[dis.ListenerKey]*dis.Listener[state.Public] {
	listeners := []*dis.Listener[state.Public]{
		NewReady(),
	}

	all := make(map[dis.ListenerKey]*dis.Listener[state.Public], len(listeners))
	for _, listener := range listeners {
		all[listener.Key()] = listener
	}

	return all
}
