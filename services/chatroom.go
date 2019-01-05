package services

import (
	"WebIM/models"
	"container/list"
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/gorilla/websocket"
	"strconv"
	"time"
)

type ChatRoom struct {
	Name string //ChatRoom name
	// 此值还没用上
	archive *list.List // All the events from the archive.

	//TODO 感觉没用
	New       <-chan models.Event    // New events coming in.
	subscribe chan models.Subscriber // Channel for new join users.

	unsubscribe chan string //Channel for exit users.

	publish chan models.Event // Send events here to publish them.

	send chan models.Event

	subscribers *list.List

	interrupt chan bool
}

func (this *ChatRoom) NewEvent(ep models.EventType, user, msg string) models.Event {
	return models.Event{ep, user, int(time.Now().Unix()), msg, nil}
}

func (this *ChatRoom) Join(user string, ws *websocket.Conn) {
	this.subscribe <- models.Subscriber{Name: user, Conn: ws}
}

func (this *ChatRoom) Leave(user string) {
	this.unsubscribe <- user
}
func (this *ChatRoom) Interrupt() {
	this.interrupt <- true
}

func (this *ChatRoom) Create() {
	this.archive = list.New()
	this.subscribe = make(chan models.Subscriber, 10)
	// Channel for exit users.
	this.unsubscribe = make(chan string, 10)
	// Send events here to publish them.
	this.publish = make(chan models.Event, 10)
	this.send = make(chan models.Event, 10)

	this.subscribers = list.New()
	this.interrupt = make(chan bool, 1)
}

func (this *ChatRoom) broadcastWebSocket(event models.Event) {
	message := &models.Message{}
	message.FromChatRoom = this.Name
	message.Text = event.Content
	event.Message = message
	data, err := json.Marshal(event)
	if err != nil {
		beego.Error("Fail to marshal event:", err)
		return
	}

	for sub := this.subscribers.Front(); sub != nil; sub = sub.Next() {
		// Immediately send event to WebSocket users.
		ws := sub.Value.(models.Subscriber).Conn
		if ws != nil {
			if ws.WriteMessage(websocket.TextMessage, data) != nil {
				// User disconnected.
				this.unsubscribe <- sub.Value.(models.Subscriber).Name
			}
		}
	}
}

func (this *ChatRoom) sendOneWebSocket(event models.Event) {
	data, err := json.Marshal(event)
	if err != nil {
		beego.Error("Fail to marshal event:", err)
		return
	}
	//TODO 重复 查找效率低
	for sub := this.subscribers.Front(); sub != nil; sub = sub.Next() {
		if sub.Value.(models.Subscriber).Name == event.User {
			ws := sub.Value.(models.Subscriber).Conn
			if ws != nil {
				if ws.WriteMessage(websocket.TextMessage, data) != nil {
					// User disconnected.
					this.unsubscribe <- sub.Value.(models.Subscriber).Name
				}
			}
		}

	}
}

func (this *ChatRoom) GetArchive() *list.List {
	return this.archive
}
func (this *ChatRoom) GetSubscribersLength() int {
	return this.subscribers.Len()
}

//有就删除，没有就不删
func (this *ChatRoom) ExitRoom(user string) {
	this.unsubscribe <- user
}

func (this *ChatRoom) ChangeRoom(user, ChangeName string) {
	this.send <- this.NewEvent(models.EVENT_CHANGE, user, ChangeName)
	this.ExitRoom(user)

}

//TODO 不合理 单聊 还是群发
func (this *ChatRoom) Send(event models.Event) {
	//TODO 有延时 可能会导致依然会接受部分消息
	this.publish <- event
}

// This function handles all incoming chan messages.
func (this *ChatRoom) Run() {
	flag := false
	for {
		select {
		case sub := <-this.subscribe:
			if !this.IsUserExist(sub.Name) {
				this.subscribers.PushBack(sub) // Add user to the end of list.
				// Publish a JOIN event.
				this.publish <- this.NewEvent(models.EVENT_JOIN, sub.Name, strconv.Itoa(this.subscribers.Len()))
			} else {
				beego.Info("Old user:", sub.Name, ";WebSocket:", sub.Conn != nil)
			}
		case event := <-this.publish:
			this.broadcastWebSocket(event)
			this.NewArchive(event)
		case unsub := <-this.unsubscribe:
			for sub := this.subscribers.Front(); sub != nil; sub = sub.Next() {
				if sub.Value.(models.Subscriber).Name == unsub {
					this.subscribers.Remove(sub)
					// Clone connection.
					ws := sub.Value.(models.Subscriber).Conn
					if ws != nil {
						ws.Close()
						beego.Error("WebSocket closed:", unsub)
					}
					this.publish <- this.NewEvent(models.EVENT_LEAVE, unsub, strconv.Itoa(this.subscribers.Len())) // Publish a LEAVE event.
					break
				}
			}
			go MonitorDeleteRun(this.Name)
		case event := <-this.send:
			for sub := this.subscribers.Front(); sub != nil; sub = sub.Next() {
				if sub.Value.(models.Subscriber).Name == event.User {
					//TODO 怎么区分单发还是群发
					this.sendOneWebSocket(event)
					this.NewArchive(event)
				}
			}
		case flag = <-this.interrupt:
		}
		if flag {
			break
		}
	}

	beego.Info(this.Name + " Run End")
	//TODO  清空所有的chan  返回所有打开错误的内容
}

func (this *ChatRoom) IsUserExist(user string) bool {
	for sub := this.subscribers.Front(); sub != nil; sub = sub.Next() {
		if sub.Value.(models.Subscriber).Name == user {
			return true
		}
	}
	return false
}

func (this *ChatRoom) NewArchive(event models.Event) {
	if this.archive.Len() >= 20 {
		this.archive.Remove(this.archive.Front())
	}
	this.archive.PushBack(event)
}
