package lifecycle

import (
	"github.com/downflux/game/engine/id/id"
)

// Component implements a subset of the Entity interface concerned with tracking
// the lifecycle of the Entity. Entities such as tanks are created inside a
// factory, and are destroyed at the end of the game or when attacked by another
// Entity.
type Component struct {
	start id.Tick
	end   id.Tick
}

// Start returns the tick at which the Entity is spawned. The tick is set in the
// constructor (delegated to each concrete impementation).
func (e Component) Start() id.Tick { return e.start }

// End returns the tick at which the Entity was destroyed. Since the game state
// is append-only, the instance itself is not removed from the internal list,
// hence the need for this marker.
func (e Component) End() id.Tick { return e.end }

// Delete marks the target Entity as having been destroyed.
func (e Component) Delete(tick id.Tick) { e.end = tick }
