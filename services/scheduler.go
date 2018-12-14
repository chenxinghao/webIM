package services

var SchedulerService *Scheduler

type Scheduler struct {
	RoomList map[string]*ChatRoom
}

func init() {
	SchedulerService = &Scheduler{}
	SchedulerService.RoomList = make(map[string]*ChatRoom)
	defaultARoom := &ChatRoom{Name: "defaultARoom"}
	defaultARoom.Create()
	go defaultARoom.Run()
	SchedulerService.RoomList["defaultARoom"] = defaultARoom
}
func (this *Scheduler) CreateChatRoom() {
	//TODO 不要写死
	room := &ChatRoom{Name: "Room"}
	room.Create()
	go room.Run()
	SchedulerService.RoomList["Room"] = room
}

//TODO 设置默认值
func (this *Scheduler) FindChatRoom(roomName string) *ChatRoom {
	if room, ok := this.RoomList[roomName]; ok {
		return room
	}
	return nil
}

func (this *Scheduler) DeleteChatRoom(roomName string) {
	if _, ok := this.RoomList[roomName]; ok {
		delete(this.RoomList, roomName)
	}
}

func (this *Scheduler) ChangeChatRoom(userName, fromRoom, toRoom string) {
	fRoom := this.FindChatRoom(fromRoom)
	if fRoom == nil {
		//TODO
	}
	if fRoom.IsUserExist(userName) {
		fRoom.ExitRoom(userName)
	}
	fRoom.ChangeRoom(userName, toRoom)

}
