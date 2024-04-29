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
	// 构造 beego server 生命周期管理器并返回
	m := new(Manager)
	m.wg = new(sync.WaitGroup)

	return m
}

// start 启动 beego server
func (m *Manager) start() error {
	// 开一个 goroutine 去启动 beego server 若main()程序退出 beego server 也会跟着退出，
	// 所以为了保证程序接收到退出信号后阻塞main()直至
	go beego.Run()

	return nil
}

// stop 程序接收到退出请求后 等待所有任务执行完成后退出
func (m *Manager) stop() error {

	if m == nil {
		return fmt.Errorf("lifeManager 未初始化")
	}

	// 程序接收到退出请求后，若还有任务未完成则阻塞住，直至所有任务完成
	m.wg.Wait()

	// 确保所有接收方返回数据后退出程序
	time.Sleep(5 * time.Second)

	return nil
}

// WaitAdd 执行请求任务时 生命
func (m *Manager) WaitAdd() {

	if m == nil {
		panic("lifeManager 未初始化")
	}

	//
	m.wg.Add(1)
}

// WaitDone beego server 生命周期管理器计数器加一
func (m *Manager) WaitDone() {

	if m == nil {
		panic("lifeManager 未初始化")
	}

	// 计数器加一
	m.wg.Done()
}

// Run 控制beego server的生命周器，负责启动 beego server，
// 还有在程序接收到退出信号的时候保证任务都完成后beego server才退出
func Run(life Lifecycle) error {

	// 启动beego server
	if err := life.start(); err != nil {
		return err
	}

	stop := make(chan os.Signal, 1)
	// 监听发送给该程序的退出信号，若接收到退出信号后传入到stop中
	signal.Notify(stop, syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	// 在程序接收退出信号前阻塞住Run()
	<-stop

	// 程序接收到退出请求后 等待所有任务执行完成后退出
	return life.stop()
}
