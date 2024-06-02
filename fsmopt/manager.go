package fsmopt

import "github.com/TheonAegor/ta-fsm-telebot"

func Strategy(strategy fsm.Strategy) fsm.ManagerOption {
	return func(config *fsm.Settings) {
		config.Strategy = strategy
	}
}

func ContextFactory(fn fsm.ContextFactory) fsm.ManagerOption {
	return func(config *fsm.Settings) {
		config.ContextFactory = fn
	}
}

func FilterProcessor(processor fsm.StateFilterProcessor) fsm.ManagerOption {
	return func(config *fsm.Settings) {
		config.StateFilterProcessor = processor
	}
}
