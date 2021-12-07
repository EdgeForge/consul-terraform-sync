package driver

import (
	"errors"
	"fmt"
	"sync"
)

// Drivers wraps the map of task-name to associated driver so that the map
// can be accessed concurrently
type Drivers struct {
	mu *sync.RWMutex

	// Map of task name to driver
	drivers map[string]Driver

	// Map of template ID to task name
	driverTemplates map[string]string

	// Tracks which driver is currently active
	active sync.Map
}

// NewDrivers returns a new drivers object
func NewDrivers() *Drivers {
	return &Drivers{
		mu:              &sync.RWMutex{},
		drivers:         make(map[string]Driver),
		driverTemplates: make(map[string]string),
	}
}

// Add adds a new driver
func (d *Drivers) Add(taskName string, driver Driver) error {
	if taskName == "" {
		return errors.New("error adding driver: task name cannot be empty")
	}
	if driver == nil {
		return fmt.Errorf("error adding driver: '%s' driver cannot be nil", taskName)
	}

	d.mu.Lock()
	defer d.mu.Unlock()

	if _, ok := d.drivers[taskName]; ok {
		return fmt.Errorf("error adding driver: a driver already exists for '%s'",
			taskName)
	}

	d.drivers[taskName] = driver
	for _, id := range driver.TemplateIDs() {
		d.driverTemplates[id] = taskName
	}
	return nil
}

// Get retrieves the driver for a task
func (d *Drivers) Get(taskName string) (Driver, bool) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	driver, ok := d.drivers[taskName]
	if !ok {
		return nil, false
	}

	return driver, true
}

func (d *Drivers) GetTask(tmplID string) (string, bool) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	taskName, ok := d.driverTemplates[tmplID]
	return taskName, ok
}

func (d *Drivers) Reset() {
	d.mu.Lock()
	defer d.mu.Unlock()

	for k := range d.drivers {
		delete(d.drivers, k)
		d.active.Delete(k)
	}
}

func (d *Drivers) Len() int {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return len(d.drivers)
}

// Map returns a copy of the map containing the drivers
func (d *Drivers) Map() map[string]Driver {
	d.mu.RLock()
	defer d.mu.RUnlock()

	copy := make(map[string]Driver, len(d.drivers))
	for k, v := range d.drivers {
		copy[k] = v
	}
	return copy
}

func (d *Drivers) SetBufferPeriod() {
	d.mu.Lock()
	defer d.mu.Unlock()
	for _, driver := range d.drivers {
		driver.SetBufferPeriod()
	}
}

func (d *Drivers) SetActive(name string) bool {
	d.active.Store(name, struct{}{})
	return true
}

func (d *Drivers) SetInactive(name string) bool {
	_, ok := d.active.Load(name)
	if ok {
		d.active.Delete(name)
	}
	return ok
}

func (d *Drivers) IsActive(name string) bool {
	_, ok := d.active.Load(name)
	return ok
}

// Delete removes the driver for the given task name from
// the map of drivers.
func (d *Drivers) Delete(taskName string) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	if taskName == "" {
		return errors.New("task name cannot be empty")
	}

	delete(d.drivers, taskName)
	return nil
}
