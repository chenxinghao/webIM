// Copyright 2013 Beego Samples authors
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package controllers

import (
	"WebIM/models"
	"encoding/json"
	"net/http"
	"time"
	"webIM/services"

	"github.com/astaxie/beego"
	"github.com/gorilla/websocket"
)

// WebSocketController handles WebSocket requests.
type WebSocketController struct {
	baseController
}

// Get method handles GET requests for WebSocketController.
func (this *WebSocketController) Get() {
	// Safe check.
	uname := this.GetString("uname")
	if len(uname) == 0 {
		this.Redirect("/", 302)
		return
	}

	this.TplName = "websocket.html"
	this.Data["IsWebSocket"] = true
	this.Data["UserName"] = uname
}

// Join method handles WebSocket requests for WebSocketController.
func (this *WebSocketController) Join() {

	uname := this.GetString("uname")
	roomName := this.GetString("room")

	if len(uname) == 0 {
		this.Redirect("/", 302)
		return
	}

	// Upgrade from http request to WebSocket.
	ws, err := websocket.Upgrade(this.Ctx.ResponseWriter, this.Ctx.Request, nil, 1024, 1024)
	if _, ok := err.(websocket.HandshakeError); ok {
		http.Error(this.Ctx.ResponseWriter, "Not a websocket handshake", 400)
		return
	} else if err != nil {
		beego.Error("Cannot setup WebSocket connection:", err)
		return
	}

	var chatRoom *services.ChatRoom
	if roomName == "" {
		roomName = "defaultRoom"
	}
	if chatRoom = services.SchedulerService.FindChatRoom(roomName); chatRoom == nil {
		beego.Info("CreateChatRoom: " + roomName)
		chatRoom = services.SchedulerService.CreateChatRoom(roomName)
	}

	chatRoom.Join(uname, ws)
	defer chatRoom.Leave(uname)

	// Join chat room.
	//Join(uname, ws)
	//defer Leave(uname)

	time.Sleep(time.Duration(2) * time.Second)

	for event := chatRoom.GetArchive().Front(); event != nil; event = event.Next() {
		ev := event.Value.(models.Event)
		data, err := json.Marshal(ev)
		if err != nil {
			beego.Error("Fail to marshal event:", err)
			return
		}
		if ws.WriteMessage(websocket.TextMessage, data) != nil {
			// User disconnected.
			chatRoom.Leave(uname)
		}
	}

	// Message receive loop.
	for {
		_, p, err := ws.ReadMessage()
		if err != nil {
			return
		}
		chatRoom.Send(chatRoom.NewEvent(models.EVENT_MESSAGE, uname, string(p)))
	}
}
