package core

import (
	"fmt"
	"github.com/LinMAD/SnapEngine/engine/entity"

	"github.com/veandco/go-sdl2/sdl"

	"github.com/LinMAD/SnapEngine/engine/graphics"
	"github.com/LinMAD/SnapEngine/engine/graphics/data"
	"github.com/LinMAD/SnapEngine/engine/graphics/factory"
	"github.com/LinMAD/SnapEngine/engine/logger"
	"github.com/LinMAD/SnapEngine/engine/platform"
)

// snapEngine core object with dependencies
type snapEngine struct {
	//
	// Engine state

	isDebug       bool
	isRunning     bool
	isHasFocus    bool
	isLevelLoaded bool
	fps           uint16

	//
	// Graphics related

	screen       platform.ScreenConfiguration
	nativeWindow *sdl.Window
	renderer     *sdl.Renderer

	//
	// sceneObjects (actors)
	sceneObjects []entity.SceneObject

	//
	// Other dependencies

	dataTextures  *data.TextureContainer
	spriteFactory *factory.SpriteFactory
	log           *logger.Logger
}

// New creates new instance of engine
func New(sc platform.ScreenConfiguration, isDebug bool) *snapEngine {
	return &snapEngine{
		isDebug:   isDebug,
		isRunning: false,
		screen:    sc,
		log:       &logger.Logger{IsDebug: isDebug},
	}
}

// Init all subsystems, create window
func (eng *snapEngine) Init() (err error) {
	eng.log.LogDebug("Initializing graphics and creating window...")

	if eng.nativeWindow, err = graphics.CreateNativeWindow(&eng.screen); err != nil {
		return err
	}
	if eng.renderer, err = graphics.CreateRenderer(eng.nativeWindow); err != nil {
		return err
	}

	eng.dataTextures = data.NewTexturesContainer(eng.renderer)
	eng.spriteFactory = factory.NewSpriteFactory(eng.renderer, eng.dataTextures)
	eng.isRunning = true

	return nil
}

// LoadComponents to engine with external logic
func (eng *snapEngine) LoadComponents(sceneObjects []entity.SceneObject) error {
	eng.log.LogDebug("Loading components...")
	eng.sceneObjects = sceneObjects

	// TODO Load in async
	for _, actor := range eng.sceneObjects {
		if err := eng.dataTextures.LoadFromFile(actor.GetDrawableInformation().TextureData); err != nil {
			return err
		}
	}

	eng.isLevelLoaded = true

	return nil
}

// UnloadComponents cleans module dependencies
func (eng *snapEngine) UnloadComponents() {
	eng.log.LogDebug("Unload components...")

	for n, t := range eng.dataTextures.GetAll() {
		if err := t.Destroy(); err != nil {
			eng.log.LogDebugWithObject(fmt.Sprintf("Texture (%s) was destroyed with error", n), err.Error())
		}
	}

	eng.isLevelLoaded = false
}

// HandleEvents like input, audio, triggers etc
func (eng *snapEngine) HandleEvents() {
	// TODO Add window, keyboard, mouse event handler

	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch t := event.(type) {
		case *sdl.QuitEvent:
			eng.isRunning = false
		case *sdl.MouseMotionEvent:
			eng.log.LogDebug(fmt.Sprintf("MouseMotionEvent - Code: %v Pos: x:%v, y:%v", t.Type, t.X, t.Y))
		case *sdl.MouseButtonEvent:
			eng.log.LogDebug(fmt.Sprintf("MouseButtonEvent - Code: %v, Clicks: %v, Pos: x:%v, y:%v", t.Button, t.Clicks, t.X, t.Y))
		case *sdl.KeyboardEvent:
			eng.log.LogDebug(fmt.Sprintf("KeyboardEvent - Code: %v, IsPressed: %v", t.Keysym.Scancode, t.State))
		}
	}
}

// HandleUpdate of engine state, physics simulation etc
func (eng *snapEngine) HandleUpdate() {
	// TODO Add Engine state handler
	eng.nativeWindow.SetTitle(fmt.Sprintf("%s |FPS: %d|", eng.screen.Title, eng.fps))

	for _, actor := range eng.sceneObjects {
		actor.OnUpdate()
	}
}

// HandleRender window frame
func (eng *snapEngine) HandleRender() error {
	if err := eng.renderer.Clear(); err != nil {
		return fmt.Errorf("renderer failed to clear frame: %s", err.Error())
	}

	for _, actor := range eng.sceneObjects {
		flipMode := sdl.FLIP_NONE
		if actor.GetDrawableInformation().IsFlipped {
			flipMode = sdl.FLIP_HORIZONTAL
		}

		if err := eng.spriteFactory.Draw(actor, flipMode); err != nil {
			return err
		}
	}

	eng.renderer.Present()

	return nil
}

// HandleClean gracefully shutdown, save state data and clean dependencies
func (eng *snapEngine) HandleClean() {
	var err error
	eng.log.LogDebug("Cleaning resources...")

	if eng.isLevelLoaded {
		eng.UnloadComponents()
	}

	if err = eng.renderer.Destroy(); err != nil {
		eng.log.LogDebugWithObject("Renderer for native window was destroyed with error", err.Error())
	}
	if err = eng.nativeWindow.Destroy(); err != nil {
		eng.log.LogDebugWithObject("Native window was destroyed with error", err.Error())
	}
}

// IsRunning return flag if application still must be executed
func (eng *snapEngine) IsRunning() bool {
	if !eng.isRunning {
		return eng.isRunning
	}

	return eng.isRunning
}

// HasFocus native window focus, not minimised and active
func (eng *snapEngine) HasFocus() bool {
	return eng.isHasFocus
}

// DeltaTime tick time in milliseconds
func (eng *snapEngine) DeltaTime() uint32 {
	return sdl.GetTicks()
}

// SetDelay waits milliseconds before continuing
func (eng *snapEngine) SetDelay(milliSeconds uint32) {
	sdl.Delay(milliSeconds)
}

// SetFps counter
func (eng *snapEngine) SetFps(fps uint16) {
	eng.fps = fps
}
