package overwatch

// AppManager is the default implementation of overwatch service management
type AppManager struct {
	services  map[string]Service
	errorChan chan error
}

// NewAppManager creates a new overwatch manager
func NewAppManager(errorChan chan error) Manager {
	return &AppManager{services: make(map[string]Service), errorChan: errorChan}
}

// Add takes in a new service to manage.
// It stops the service if it already exist in the manager and is running
// It then starts the newly added service
func (m *AppManager) Add(service Service) {
	// check for existing service
	if currentService, ok := m.services[service.Name()]; ok {
		if currentService.Hash() == service.Hash() {
			return // the exact same service, no changes, so move along
		}
		currentService.Shutdown() //shutdown the listener since a new one is starting
	}
	m.services[service.Name()] = service

	//start the service!
	go m.serviceRun(service)
}

// Remove shutdowns the service by name and removes it from its current management list
func (m *AppManager) Remove(name string) {
	if currentService, ok := m.services[name]; ok {
		currentService.Shutdown()
	}
	delete(m.services, name)
}

// Services returns all the current Services being managed
func (m *AppManager) Services() []Service {
	values := []Service{}
	for _, value := range m.services {
		values = append(values, value)
	}
	return values
}

func (m *AppManager) serviceRun(service Service) {
	err := service.Run()
	if err != nil && m.errorChan != nil {
		m.errorChan <- err
	}
}
