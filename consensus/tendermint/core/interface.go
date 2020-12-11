package core

//Engine abstract the core's functions
//Note that backend and other packages doesn't care about core's internal logic.
//It only requires core to start receiving/handling messages
//The sending of events/message from core to backend will be done by calling accessing Backend.EventMux()
type Engine interface {
	Start() error
	Stop() error
}
