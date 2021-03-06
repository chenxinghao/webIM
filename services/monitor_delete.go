package services

import (
	"time"
)

type MonitorDelete struct {
	finishFlag int
	RoomName   string
}

func (this *MonitorDelete) monitorDeleteStatus(roomName string) bool {
	chatRoom := SchedulerService.FindChatRoom(roomName)
	if chatRoom == nil {
		return false
	}
	if chatRoom.GetSubscribersLength() != 0 {
		return false
	}
	return true
}

func (this *MonitorDelete) deleteCounter(countNumber int) {

	for i := 0; i < countNumber; i++ {
		if !this.monitorDeleteStatus(this.RoomName) {
			this.finishFlag = 0 //没有结束
			return
		}
		if this.finishFlag == -1 { //被中断
			return
		}
		time.Sleep(time.Duration(1) * time.Second)
	}
	this.finishFlag = 1 //可以结束
	return
}

func (this *MonitorDelete) cancel() {
	this.finishFlag = -1
}

func MonitorDeleteRun(roomName string) {

	if mdp, ok := SchedulerService.ClearList[roomName]; ok {
		mdp.cancel()

	}
	monitorDelete := MonitorDelete{finishFlag: 0, RoomName: roomName}
	SchedulerService.ClearList[roomName] = &monitorDelete

	monitorDelete.deleteCounter(120)
	if monitorDelete.finishFlag == 1 {
		chatRoom := SchedulerService.FindChatRoom(monitorDelete.RoomName)
		SchedulerService.DeleteChatRoom(monitorDelete.RoomName)
		chatRoom.Interrupt()
	}
	if monitorDelete.finishFlag != -1 {
		delete(SchedulerService.ClearList, roomName)
	}

}
