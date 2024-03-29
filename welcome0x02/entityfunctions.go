package main

import "github.com/Lealen/engi"

var entitiesChange bool
var entititesFunctions = make(map[engi.Scene][]*Entity)
var entitiesUpdateMap [engi.HighestGround + 1][]*Entity //= make([][]*Entity, engi.HighestGround)

var lastwindowwidth, lastwindowheight float32

var mouseX, mouseY float32
var prevMouseX, prevMouseY float32

func IsMouseOn(e *Entity) bool {
	if mouseX > e.Space.Position.X && mouseX < (e.Space.Position.X+e.Space.Width) &&
		mouseY > e.Space.Position.Y && mouseY < (e.Space.Position.Y+e.Space.Height) {
		return true
	}
	return false
}

func UpdateEntities(dt float32) {
	if entitiesChange {
		for k := range entitiesUpdateMap {
			entitiesUpdateMap[k] = nil
		}
		for _, v := range entititesFunctions[engi.CurrentScene()] {
			if v.priority > engi.HighestGround || v.priority < 0 {
				continue
			}
			entitiesUpdateMap[v.priority] = append(entitiesUpdateMap[v.priority], v)
		}
	}

	//game coordinates
	mouseX = engi.Mouse.X*CameraGetZ()*(engi.Width()/engi.WindowWidth()) + CameraGetX() - (engi.Width()/2)*CameraGetZ()
	mouseY = engi.Mouse.Y*CameraGetZ()*(engi.Height()/engi.WindowHeight()) + CameraGetY() - (engi.Height()/2)*CameraGetZ()

	for _, v := range entititesFunctions[engi.CurrentScene()] {
		if !v.Initialized {
			if v.OnFirstUpdate != nil {
				v.OnFirstUpdate(v)
			}
			v.Initialized = true
		}
		if v.OnUpdate != nil {
			v.OnUpdate(v, dt)
		}
	}

	if lastwindowwidth == 0 && lastwindowheight == 0 {
		lastwindowwidth = engi.WindowWidth()
		lastwindowheight = engi.WindowHeight()
	} else if lastwindowwidth != engi.WindowWidth() || lastwindowheight != engi.WindowHeight() {
		for _, v := range entititesFunctions[engi.CurrentScene()] {
			if v.OnWindowResize != nil {
				v.OnWindowResize(v)
			}
		}
		lastwindowwidth = engi.WindowWidth()
		lastwindowheight = engi.WindowHeight()
	}

	ontop := true
	ontopchanged := false
	for prio := engi.HighestGround; prio >= engi.Background; prio-- {
		for _, v := range entitiesUpdateMap[prio] {
			if v.Mouse == nil {
				continue
			}
			ison := IsMouseOn(v)
			if ison && (ontop || v.IgnoreWhatIsOnTop) && engi.Keys.Get(engi.MouseButtonLeft).Down() {
				if !v.IgnoreWhatIsOnTop {
					ontopchanged = true
				}
				if !v.IsClicked {
					if v.OnPress != nil {
						v.OnPress(v)
					}
					v.IsClicked = true
				}
				if v.OnClicked != nil {
					v.OnClicked(v)
				}
			} else if v.IsClicked {
				if v.OnRelease != nil {
					v.OnRelease(v)
				}
				v.IsClicked = false
			}
			if (ontop || v.IgnoreWhatIsOnTop) && v.OnDragged != nil && v.Mouse.Dragged && engi.Keys.Get(engi.MouseButtonLeft).Down() {
				if !v.IgnoreWhatIsOnTop {
					ontopchanged = true
				}
				v.OnDragged(v)
			}
			if ison && (ontop || v.IgnoreWhatIsOnTop) && engi.Keys.Get(engi.MouseButtonRight).Down() {
				if !v.IgnoreWhatIsOnTop {
					ontopchanged = true
				}
				if !v.IsRightClicked {
					if v.OnRightPress != nil {
						v.OnRightPress(v)
					}
					v.IsRightClicked = true
				}
				if v.OnRightClicked != nil {
					v.OnRightClicked(v)
				}
			} else if v.IsRightClicked {
				if v.OnRightRelease != nil {
					v.OnRightRelease(v)
				}
				v.IsRightClicked = false
			}
			if v.OnEnter != nil && v.Mouse.Enter {
				v.OnEnter(v)
			}
			if v.OnLeave != nil && v.Mouse.Leave {
				v.OnLeave(v)
			}

			if ontopchanged {
				ontop = false
				ontopchanged = false
			}
		}
	}

	prevMouseX = mouseX
	prevMouseY = mouseY
}
