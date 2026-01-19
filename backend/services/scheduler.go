package services

import (
	"log"
	"time"

	"remote-debug-platform/models"
)

type Scheduler struct {
	stop chan struct{}
}

func NewScheduler() *Scheduler {
	return &Scheduler{
		stop: make(chan struct{}),
	}
}

func (s *Scheduler) Start() {
	go s.checkExpiredPorts()
	log.Println("定时任务已启动")
}

func (s *Scheduler) Stop() {
	close(s.stop)
}

// checkExpiredPorts 检查过期端口
func (s *Scheduler) checkExpiredPorts() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-s.stop:
			return
		case <-ticker.C:
			s.doCheckExpiredPorts()
		}
	}
}

func (s *Scheduler) doCheckExpiredPorts() {
	ports, err := models.GetExpiredPorts()
	if err != nil {
		log.Printf("获取过期端口失败: %v", err)
		return
	}

	pm := GetPortManager()
	for _, port := range ports {
		log.Printf("端口 %d 已过期，正在关闭", port.Port)

		// 停止监听
		pm.StopPortListener(port.Port)

		// 更新状态
		port.Status = "expired"
		port.Update()

		// 设备下线
		models.SetPortDevicesOffline(port.ID)
	}
}
