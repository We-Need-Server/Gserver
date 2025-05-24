package mediator

//
//import "fmt"
//
//internal_types Component interface {
//	Register(m *Mediator)
//	Send(receiverName string, message interface{})
//	Receive(senderName string, message interface{})
//}
//
//internal_types Interface interface {
//	Register()
//	Notify()
//}
//
//internal_types Mediator struct {
//	components map[string]Component
//}
//
//func NewMediator() *Mediator {
//	return &Mediator{make(map[string]Component)}
//}
//
//func (m *Mediator) Register(name string, component Component) (string, error) {
//	if _, exists := m.components[name]; exists {
//		return name, fmt.Errorf("already existed component")
//	}
//	m.components[name] = component
//	component.Register(m)
//	return name, nil
//}
//
//func (m *Mediator) Notify(senderName string, receiverName string, message interface{}) {
//	m.components[receiverName].Receive(senderName, message)
//}
