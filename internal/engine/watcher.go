package engine

import (
	"github.com/jademperor/common/configs"
	"time"

	"github.com/jademperor/api-proxier/internal/logger"
	"github.com/jademperor/api-proxier/internal/plugin/cache"
	"github.com/jademperor/common/etcdutils"
)

var (
	clusterWatcher  *etcdutils.Watcher // cluster watcher
	cacheWatcher    *etcdutils.Watcher // cache watcher
	apisWatcher     *etcdutils.Watcher
	routingsWatcher *etcdutils.Watcher
	// rbacWatcher     *etcdutils.Watcher // rabc plugin watcher
	// and etc

	defaultDuration = 2 * time.Second
)

// initialWatchers ...
func (e *Engine) initialWatchers() {
	clusterWatcher = etcdutils.NewWatcher(e.kapi, defaultDuration, configs.ClustersKey)
	apisWatcher = etcdutils.NewWatcher(e.kapi, defaultDuration, configs.APIsKey)
	routingsWatcher = etcdutils.NewWatcher(e.kapi, defaultDuration, configs.RoutingsKey)
	// rbacWatcher = etcdutils.NewWatcher(e.kapi, defaultDuration, configs.RbacKey)
	cacheWatcher = etcdutils.NewWatcher(e.kapi, defaultDuration, configs.CacheKey)

	go clusterWatcher.Watch(e.clusterCallback)
	go apisWatcher.Watch(e.apisCallback)
	go routingsWatcher.Watch(e.routingsCallback)
	go cacheWatcher.Watch(e.cacheCallback)
	// go rbacWatcher.Watch(e.rbacCallback)
}

func (e *Engine) clusterCallback(op etcdutils.OpCode, k, v string) {
	logger.Logger.Infof("clusters Op: %d, key: %s, value: %s", op, k, v)
	e.prepareClusters()
}

func (e *Engine) apisCallback(op etcdutils.OpCode, k, v string) {
	logger.Logger.Infof("apis Op: %d, key: %s, value: %s", op, k, v)
	e.prepareAPIs()
}

func (e *Engine) routingsCallback(op etcdutils.OpCode, k, v string) {
	logger.Logger.Infof("routings Op: %d, key: %s, value: %s", op, k, v)
	e.prepareRoutings()
}

func (e *Engine) cacheCallback(op etcdutils.OpCode, k, v string) {
	logger.Logger.Infof("cache Op: %d, key: %s, value: %s", op, k, v)
	e.prepareCache(e.allPlugins[1].(*cache.Cache))
}

// func (e *Engine) rbacCallback(op etcdutils.OpCode, k, v string) {
// 	// TODO
// 	logger.Logger.Infof("RBAC Op: %d, key: %s, value: %s", op, k, v)
// 	// notify reload
// }