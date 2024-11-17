package observer

type Subject interface {
	Register(observer Observer)
	Unregister(observer Observer)
	Notify(chatID string, message string)
}

// ConcreteSubject - конкретний субʼєкт, який зберігає список спостерігачів та повідомляє їх
type ConcreteSubject struct {
	observers []Observer
}

// Register - додає спостерігача до субʼєкта
func (s *ConcreteSubject) Register(observer Observer) {
	s.observers = append(s.observers, observer)
}

// Unregister - видаляє спостерігача із субʼєкта
func (s *ConcreteSubject) Unregister(observer Observer) {
	for i, obs := range s.observers {
		if obs == observer {
			s.observers = append(s.observers[:i], s.observers[i+1:]...)
			break
		}
	}
}

// Notify - повідомляє усіх спостерігачів про подію
func (s *ConcreteSubject) Notify(chatID string, message string) {
	for _, observer := range s.observers {
		observer.Notify(chatID, message)
	}
}
