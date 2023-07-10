package fsm

import (
	tele "gopkg.in/telebot.v3"
)

// Handler is object for handling  updates with FSM FSMContext
type Handler func(c tele.Context, state Context) error

// ContextMakerFunc alias for function for create new context.
// You can use custom Context implementation.
type ContextMakerFunc = func(ctx tele.Context, storage Storage) Context // TODO: add error to return values

// Manager is object for managing FSM, binding handlers.
type Manager struct {
	bot          *tele.Bot
	group        *tele.Group // handlers will add to group
	store        Storage
	handlers     handlerStorage
	contextMaker ContextMakerFunc
}

// NewManager returns new Manger.
func NewManager(
	bot *tele.Bot,
	group *tele.Group,
	storage Storage,
	ctxMaker ContextMakerFunc,
) *Manager {
	if group == nil {
		group = bot.Group()
	}
	if ctxMaker == nil {
		ctxMaker = NewFSMContext
	}
	return &Manager{
		bot:          bot,
		group:        group,
		store:        storage,
		contextMaker: ctxMaker,
		handlers:     make(handlerStorage),
	}
}

// Group handlers for manger.
func (m *Manager) Group() *tele.Group {
	return m.group
}

func (m *Manager) With(g *tele.Group) *Manager {
	return &Manager{
		bot:          m.bot,
		group:        g,
		store:        m.store,
		handlers:     m.handlers,
		contextMaker: m.contextMaker,
	}
}

// SetContextMaker sets new context maker to current Manager instance.
func (m *Manager) SetContextMaker(contextMaker ContextMakerFunc) {
	m.contextMaker = contextMaker
}

func (m *Manager) NewGroup() *Manager {
	return m.With(m.bot.Group())
}

// Use add middlewares to group.
func (m *Manager) Use(middlewares ...tele.MiddlewareFunc) {
	m.group.Use(middlewares...)
}

// Bind adds handler (with FSMContext) with filter on state.
//
// Difference between Bind and Handle methods what Handle require Filter objects.
// And this method can work with only one state.
func (m *Manager) Bind(end interface{}, state State, h Handler, middlewares ...tele.MiddlewareFunc) {
	m.Handle(F(end, state), h, middlewares...)
}

// Handle adds handler to group chain with filter on states.
// Allowed use more handler for one endpoint.
func (m *Manager) Handle(f Filter, h Handler, middlewares ...tele.MiddlewareFunc) {
	endpoint := f.CallbackUnique()
	m.handlers.add(endpoint, h, f.States)
	m.group.Handle(
		endpoint,
		m.HandlerAdapter(m.handlers.forEndpoint(endpoint)),
		middlewares...,
	)
}

// HandlerAdapter create telebot.HandlerFunc object for Handler with FSM FSMContext.
func (m *Manager) HandlerAdapter(handler Handler) tele.HandlerFunc {
	return func(c tele.Context) error {
		return handler(c, m.contextMaker(c, m.store))
	}
}

// Storage returns manger storage instance.
func (m *Manager) Storage() Storage {
	return m.store
}
