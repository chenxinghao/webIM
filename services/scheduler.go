package services

var SchedulerService *Scheduler

type Scheduler struct {
	RoomList  map[string]*ChatRoom
	ClearList chan string
}

func init() {
	SchedulerService = &Scheduler{}
	SchedulerService.RoomList = make(map[string]*ChatRoom)
	defaultARoom := &ChatRoom{Name: "defaultRoom"}
	defaultARoom.Create()
	go defaultARoom.Run()
	SchedulerService.RoomList["defaultRoom"] = defaultARoom
	//需要限制长度
	SchedulerService.ClearList = make(chan string)
}
func (this *Scheduler) CreateChatRoom(name string) *ChatRoom {
	room := &ChatRoom{Name: name}
	room.Create()
	go room.Run()
	SchedulerService.RoomList[name] = room
	return room
}

//TODO 设置默认值
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

//TODO 没用到这个方法
func (this *Scheduler) ChangeChatRoom(userName, fromRoom, toRoom string) {
	fRoom := this.FindChatRoom(fromRoom)
	if fRoom == nil {

	}
	if fRoom.IsUserExist(userName) {
		fRoom.ExitRoom(userName)
	}
	fRoom.ChangeRoom(userName, toRoom)

}
