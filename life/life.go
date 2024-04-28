package life

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	beego "github.com/beego/beego/v2/server/web"
	_ "github.com/go-sql-driver/mysql"
)

type Lifecycle interface {
	start() error
	stop() error
	WaitAdd()
	WaitDone()
}

type Manager struct {
	wg *sync.WaitGroup
}

func NewLifeManager() *Manager {

	m := new(Manager)
	m.wg = new(sync.WaitGroup)

	return m
}

func (m *Manager) start() error {

	go beego.Run()

	return nil
}

func (m *Manager) stop() error {

	if m == nil {
		return fmt.Errorf("lifeManager is nil")
	}

	m.wg.Wait()

	// 确保所有请求返回数据后退出
	time.Sleep(5 * time.Second)

	return nil
}

func (m *Manager) WaitAdd() {

	if m == nil {
		panic("lifeManager is nil")
	}

	m.wg.Add(1)
}

func (m *Manager) WaitDone() {

	if m == nil {
		panic("lifeManager is nil")
	}

	m.wg.Done()

}

func Run(life Lifecycle) error {

	if err := life.start(); err != nil {
		return err
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	<-stop

	return life.stop()
}
