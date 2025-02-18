package test

import (
	config2 "github.com/doubletrey/crawlab-core/config"
	"github.com/doubletrey/crawlab-core/entity"
	"github.com/doubletrey/crawlab-core/interfaces"
	service2 "github.com/doubletrey/crawlab-core/models/service"
	"github.com/doubletrey/crawlab-core/node/service"
	"github.com/doubletrey/crawlab-core/utils"
	"go.uber.org/dig"
	"io/ioutil"
	"os"
	"path"
	"testing"
	"time"
)

func init() {
	var err error
	T, err = NewTest()
	if err != nil {
		panic(err)
	}
}

var T *Test

type Test struct {
	DefaultSvc       interfaces.NodeMasterService
	MasterSvc        interfaces.NodeMasterService
	WorkerSvc        interfaces.NodeWorkerService
	MasterSvcMonitor interfaces.NodeMasterService
	WorkerSvcMonitor interfaces.NodeWorkerService
	ModelSvc         service2.ModelService
}

func NewTest() (res *Test, err error) {
	// test
	t := &Test{}

	// recreate config directory path
	_ = os.RemoveAll(config2.DefaultConfigDirPath)
	_ = os.MkdirAll(config2.DefaultConfigDirPath, os.FileMode(0766))

	// master config and settings
	masterNodeConfigName := "config-master.json"
	masterNodeConfigPath := path.Join(config2.DefaultConfigDirPath, masterNodeConfigName)
	if err := ioutil.WriteFile(masterNodeConfigPath, []byte("{\"key\":\"master\",\"is_master\":true}"), os.FileMode(0766)); err != nil {
		return nil, err
	}
	masterHost := "0.0.0.0"
	masterPort := "9667"

	// worker config and settings
	workerNodeConfigName := "config-worker.json"
	workerNodeConfigPath := path.Join(config2.DefaultConfigDirPath, workerNodeConfigName)
	if err = ioutil.WriteFile(workerNodeConfigPath, []byte("{\"key\":\"worker\",\"is_master\":false}"), os.FileMode(0766)); err != nil {
		return nil, err
	}
	workerHost := "localhost"
	workerPort := masterPort

	// master for monitor config and settings
	masterNodeMonitorConfigName := "config-master-monitor.json"
	masterNodeMonitorConfigPath := path.Join(config2.DefaultConfigDirPath, masterNodeMonitorConfigName)
	if err := ioutil.WriteFile(masterNodeMonitorConfigPath, []byte("{\"key\":\"master-monitor\",\"is_master\":true}"), os.FileMode(0766)); err != nil {
		return nil, err
	}
	masterMonitorHost := masterHost
	masterMonitorPort := "9668"

	// worker for monitor config and settings
	workerNodeMonitorConfigName := "config-worker-monitor.json"
	workerNodeMonitorConfigPath := path.Join(config2.DefaultConfigDirPath, workerNodeMonitorConfigName)
	if err := ioutil.WriteFile(workerNodeMonitorConfigPath, []byte("{\"key\":\"worker-monitor\",\"is_master\":false}"), os.FileMode(0766)); err != nil {
		return nil, err
	}
	workerMonitorHost := workerHost
	workerMonitorPort := masterMonitorPort

	// dependency injection
	c := dig.New()
	if err := c.Provide(service.ProvideMasterService(
		masterNodeConfigPath,
		service.WithMonitorInterval(3*time.Second),
		service.WithAddress(entity.NewAddress(&entity.AddressOptions{
			Host: masterHost,
			Port: masterPort,
		})),
	)); err != nil {
		return nil, err
	}
	if err := c.Provide(service.ProvideWorkerService(
		workerNodeConfigPath,
		service.WithHeartbeatInterval(1*time.Second),
		service.WithAddress(entity.NewAddress(&entity.AddressOptions{
			Host: workerHost,
			Port: workerPort,
		})),
	)); err != nil {
		return nil, err
	}
	if err := c.Provide(service2.NewService); err != nil {
		return nil, err
	}
	if err := c.Invoke(func(masterSvc interfaces.NodeMasterService, workerSvc interfaces.NodeWorkerService, modelSvc service2.ModelService) {
		t.MasterSvc = masterSvc
		t.WorkerSvc = workerSvc
		t.ModelSvc = modelSvc
	}); err != nil {
		return nil, err
	}

	// default service
	t.DefaultSvc, err = service.NewMasterService()
	if err != nil {
		return nil, err
	}

	// master and worker for monitor
	t.MasterSvcMonitor, err = service.NewMasterService(
		service.WithConfigPath(masterNodeMonitorConfigPath),
		service.WithAddress(entity.NewAddress(&entity.AddressOptions{
			Host: masterMonitorHost,
			Port: masterMonitorPort,
		})),
		service.WithMonitorInterval(3*time.Second),
		service.WithStopOnError(),
	)
	if err != nil {
		return nil, err
	}
	t.WorkerSvcMonitor, err = service.NewWorkerService(
		service.WithConfigPath(workerNodeMonitorConfigPath),
		service.WithAddress(entity.NewAddress(&entity.AddressOptions{
			Host: workerMonitorHost,
			Port: workerMonitorPort,
		})),
		service.WithHeartbeatInterval(1*time.Second),
		service.WithStopOnError(),
	)
	if err != nil {
		return nil, err
	}

	// removed all data in db
	_ = t.ModelSvc.DropAll()

	// visualize dependencies
	if err := utils.VisualizeContainer(c); err != nil {
		return nil, err
	}

	return t, nil
}

func (t *Test) Setup(t2 *testing.T) {
	if err := t.ModelSvc.DropAll(); err != nil {
		panic(err)
	}
	_ = os.RemoveAll(config2.DefaultConfigDirPath)
	t2.Cleanup(t.Cleanup)
}

func (t *Test) Cleanup() {
	if err := t.ModelSvc.DropAll(); err != nil {
		panic(err)
	}
	_ = os.RemoveAll(config2.DefaultConfigDirPath)
}

func (t *Test) StartMasterWorker() {
	startMasterWorker()
}

func (t *Test) StopMasterWorker() {
	stopMasterWorker()
}

func startMasterWorker() {
	go T.MasterSvc.Start()
	time.Sleep(1 * time.Second)
	go T.WorkerSvc.Start()
	time.Sleep(1 * time.Second)
}

func stopMasterWorker() {
	go T.MasterSvc.Stop()
	time.Sleep(1 * time.Second)
	go T.WorkerSvc.Stop()
	time.Sleep(1 * time.Second)
}

func startMasterWorkerMonitor() {
	go T.MasterSvcMonitor.Start()
	time.Sleep(1 * time.Second)
	go T.WorkerSvcMonitor.Start()
	time.Sleep(1 * time.Second)
}

func stopMasterWorkerMonitor() {
	go T.MasterSvcMonitor.Stop()
	time.Sleep(1 * time.Second)
	go T.WorkerSvcMonitor.Stop()
	time.Sleep(1 * time.Second)
}
