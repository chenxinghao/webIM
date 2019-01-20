package services

var SchedulerService *Scheduler

type Scheduler struct {
	RoomList  map[string]*ChatRoom
	ClearList map[string]*MonitorDelete
}

func init() {
	SchedulerService = &Scheduler{}
	SchedulerService.RoomList = make(map[string]*ChatRoom)
	defaultARoom := &ChatRoom{Name: "defaultRoom"}
	defaultARoom.Create()
	go defaultARoom.Run()
	SchedulerService.RoomList["defaultRoom"] = defaultARoom
	//需要限制长度
	SchedulerService.ClearList = make(map[string]*MonitorDelete)
}
func (this *Scheduler) CreateChatRoom(name string) *ChatRoom {
	room := &ChatRoom{Name: name}
	room.Create()
	go room.Run()
	SchedulerService.RoomList[name] = room
	return room
}

func (this *Scheduler) FindChatRoom(roomName string) *ChatRoom {
	if room, ok := this.RoomList[roomName]; ok {
		return room
	}
	return nil
}

func (this *Scheduler) DeleteChatRoom(roomName string) {
	if room, ok := this.RoomList[roomName]; ok {
		room.Interrupt()
		delete(this.RoomList, roomName)
	}
}
