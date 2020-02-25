/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package mdns

import (
	// Frameworks
	"context"
	"sync"
	"time"

	gopi "github.com/djthorpe/gopi/v2"
	base "github.com/djthorpe/gopi/v2/base"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type ServiceDB struct {
	Discovery gopi.RPCServiceDiscovery
	Listener  ListenerIface
	Bus       gopi.Bus
}

type servicedb struct {
	discovery  gopi.RPCServiceDiscovery
	listener   ListenerIface
	bus        gopi.Bus
	services   map[string]map[string]instance
	queue      []string
	stopName   gopi.Channel
	stopRecord gopi.Channel
	stopLookup chan struct{}

	base.Unit
	sync.RWMutex
	sync.WaitGroup
}

type instance struct {
	service gopi.RPCServiceRecord
	expires time.Time
	ttl     time.Duration
}

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION gopi.Unit

func (ServiceDB) Name() string { return "gopi/mdns/servicedb" }

func (config ServiceDB) New(log gopi.Logger) (gopi.Unit, error) {
	this := new(servicedb)

	// Init
	if err := this.Unit.Init(log); err != nil {
		return nil, err
	} else if err := this.Init(config); err != nil {
		return nil, err
	}

	// Success
	return this, nil
}

func (this *servicedb) Init(config ServiceDB) error {

	// Check Discovery
	if config.Discovery == nil {
		return gopi.ErrBadParameter.WithPrefix("Discovery")
	} else {
		this.discovery = config.Discovery
	}

	// Check Listener
	if config.Listener == nil {
		return gopi.ErrBadParameter.WithPrefix("Listener")
	} else {
		this.listener = config.Listener
	}

	// Check Bus
	if config.Bus == nil {
		return gopi.ErrBadParameter.WithPrefix("Bus")
	} else {
		this.bus = config.Bus
	}

	// Empty service array
	this.services = make(map[string]map[string]instance)
	this.queue = make([]string, 0, 10)

	// Start background message processor
	if this.stopName = this.listener.Subscribe(QUEUE_NAME, this.EventHandler); this.stopName == 0 {
		return gopi.ErrInternalAppError
	}
	if this.stopRecord = this.listener.Subscribe(QUEUE_RECORD, this.EventHandler); this.stopRecord == 0 {
		this.listener.Unsubscribe(this.stopName)
		return gopi.ErrInternalAppError
	}

	// Background lookup
	this.stopLookup = make(chan struct{})
	this.WaitGroup.Add(1)
	go this.LookupHandler(this.stopLookup)

	// Return success
	return nil
}

func (this *servicedb) Close() error {

	// Wait for lookup to end
	close(this.stopLookup)
	this.WaitGroup.Wait()

	// Wait for EventHandler to end
	this.listener.Unsubscribe(this.stopName)
	this.listener.Unsubscribe(this.stopRecord)

	// Release resources
	this.stopName = 0
	this.stopRecord = 0
	this.listener = nil
	this.services = nil
	this.queue = nil
	this.discovery = nil
	this.bus = nil

	// Return success
	return this.Unit.Close()
}

func (this *servicedb) String() string {
	return "<" + this.Log.Name() + ">"
}

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION gopi.RPCServiceDiscovery

// Lookup service instances by name
func (this *servicedb) Lookup(ctx context.Context, service string) ([]gopi.RPCServiceRecord, error) {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()

	if len(this.services) == 0 {
		return this.discovery.Lookup(ctx, service)
	} else if instances, exists := this.services[service]; exists == false {
		return this.discovery.Lookup(ctx, service)
	} else {
		records := make([]gopi.RPCServiceRecord, 0, len(instances))
		for _, instance := range instances {
			if instance.expires.After(time.Now()) {
				records = append(records, instance.service)
			}
		}
		return records, nil
	}
}

// Return list of service names
func (this *servicedb) EnumerateServices(ctx context.Context) ([]string, error) {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()

	return this.discovery.EnumerateServices(ctx)
	/*
		if len(this.services) == 0 {
			return this.discovery.EnumerateServices(ctx)
		} else {
			services := make([]string, 0, len(this.services))
			for key := range this.services {
				services = append(services, key)
			}
			return services, nil
		}*/

}

////////////////////////////////////////////////////////////////////////////////
// BACKGROUND PROCESSOR

func (this *servicedb) EventHandler(value interface{}) {
	if evt, ok := value.(gopi.RPCEvent); ok == false || evt == nil {
		// No nothing
	} else if record := evt.Service(); record.Name == "" {
		// No nothing
	} else if evt.Type() == gopi.RPC_EVENT_SERVICE_NAME && evt.Service().Name != "" {
		this.AddServiceName(evt.Service().Name)
	} else if evt.Type() == gopi.RPC_EVENT_SERVICE_RECORD && evt.Service().Name != "" {
		this.AddServiceName(evt.Service().Service)
		this.AddServiceInstance(evt.Service(), evt.TTL())
	}
}

func (this *servicedb) LookupHandler(<-chan struct{}) {
	lookupTimer := time.NewTimer(2 * time.Second)
	expireTimer := time.NewTimer(DISCOVERY_LOOKUP_DELTA * 2)
FOR_LOOP:
	for {
		select {
		case <-lookupTimer.C:
			if name := this.popServiceName(); name != "" {
				this.lookupServices(name)
			}
			lookupTimer.Reset(DISCOVERY_LOOKUP_DELTA)
		case <-expireTimer.C:
			this.expireInstances()
			expireTimer.Reset(DISCOVERY_LOOKUP_DELTA)
		case <-this.stopLookup:
			lookupTimer.Stop()
			expireTimer.Stop()
			break FOR_LOOP
		}
	}
	this.WaitGroup.Done()
}

func (this *servicedb) AddServiceName(name string) {
	if this.shouldServiceName(name) {
		this.addServiceName(name)
	}
}

func (this *servicedb) AddServiceInstance(service gopi.RPCServiceRecord, ttl time.Duration) {
	name, key := nameKeyForService(service)
	if ttl > 0 {
		this.addServiceInstance(name, key, instance{service, time.Now().Add(ttl), ttl})
	} else {
		this.removeServiceInstanceForKey(name, key, gopi.RPC_EVENT_SERVICE_REMOVED)
	}
}

func (this *servicedb) addServiceInstance(name, key string, value instance) {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()

	// Create the map of service names
	if _, exists := this.services[name]; exists == false {
		this.services[name] = make(map[string]instance, 1)
	}

	// If an instance already exists then update
	if instance_, exists := this.services[name][key]; exists {
		this.services[name][key] = value
		if serviceEquals(instance_.service, value.service) == false {
			this.bus.Emit(NewEvent(this, gopi.RPC_EVENT_SERVICE_UPDATED, value.service, value.ttl))
		}
	} else {
		this.bus.Emit(NewEvent(this, gopi.RPC_EVENT_SERVICE_ADDED, value.service, value.ttl))
		this.services[name][key] = value
	}
}

func (this *servicedb) removeServiceInstanceForKey(name, key string, eventType gopi.RPCEventType) {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()

	// Check for name already existing
	if _, exists := this.services[name]; exists == false {
		return
	}
	// Remove by key
	if instance, exists := this.services[name][key]; exists == false {
		return
	} else {
		// Remove the record
		delete(this.services[name], key)

		// Remove service
		if len(this.services[name]) == 0 {
			delete(this.services, name)
		}

		// Emit the event
		this.bus.Emit(NewEvent(this, eventType, instance.service, 0))
	}
}

// Expire instances
func (this *servicedb) expireInstances() {
	for {
		if name, key := this.instanceToExpire(); name != "" && key != "" {
			this.removeServiceInstanceForKey(name, key, gopi.RPC_EVENT_SERVICE_EXPIRED)
		} else {
			return
		}
	}
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

// Return first service name in the queue or empty string
func (this *servicedb) popServiceName() string {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()
	if len(this.queue) == 0 {
		return ""
	} else {
		name := this.queue[0]
		this.queue = this.queue[1:]
		return name
	}
}

// Lookup services
func (this *servicedb) lookupServices(name string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	// No need to actually record the records here
	if _, err := this.discovery.Lookup(ctx, name); err != nil {
		return err
	} else {
		return nil
	}
}

// Return an instance which needs to be expired
func (this *servicedb) instanceToExpire() (string, string) {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()

	for name, instances := range this.services {
		for key, instance := range instances {
			if instance.expires.Before(time.Now()) {
				return name, key
			}
		}
	}

	return "", ""
}

func (this *servicedb) addServiceName(name string) {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()
	this.queue = append(this.queue, name)
}

func (this *servicedb) shouldServiceName(name string) bool {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()

	// Ignore if service name is already registered
	if _, exists := this.services[name]; exists {
		return false
	}
	// Ignore if already in the queue
	if inStringSlice(name, this.queue) {
		return false
	}
	// Append to the end of the queue
	return true
}

func nameKeyForService(record gopi.RPCServiceRecord) (string, string) {
	if record.Name == "" {
		return "", ""
	} else {
		return record.Service, Quote(record.Name) + "." + record.Service
	}
}

func inStringSlice(value string, slice []string) bool {
	for _, elem := range slice {
		if elem == value {
			return true
		}
	}
	return false
}
