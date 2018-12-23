var socket;

$(document).ready(function () {
    // alert(getUrlParam("room"));
    if  (getUrlParam("room")!=''){
        $('#room_name').text(getUrlParam("room"));
    }
    // Create a socket
    socket = new WebSocket('ws://' + window.location.host + '/ws/join?uname=' + $('#uname').text()+'&room='+getUrlParam("room"));
    // Message received on the socket
    socket.onmessage = function (event) {
        var data = JSON.parse(event.data);
        var li = document.createElement('li');

        console.log(data);

        switch (data.Type) {
        case 0: // JOIN
            $("#members").text('在线人数：'+data.Content);
            if (data.User == $('#uname').text()) {
                li.innerText = 'You joined the chat room.';
            } else {
                li.innerText = data.User + ' joined the chat room.';
            }
            break;
        case 1: // LEAVE
            $("#members").text('在线人数：'+data.Content);
            li.innerText = data.User + ' left the chat room.';
            break;
            case 2: // MESSAGE
            var username = document.createElement('strong');
            var content = document.createElement('span');


            var fromChatRoom;
            var text;
            for(var p in data){//遍历json对象的每个key/value对,p为key
                if (p =="Message"){
                    for(var v in data[p]){
                        // alert(v+ " " +data[p][v]);
                       if (v=="fromChatRoom"){
                           fromChatRoon=data[p][v];
                       }
                       if (v=="text"){
                            text=data[p][v];
                       }
                    }
                }
            }
            username.innerText = fromChatRoon+"."+data.User;
            content.innerText = text;

            li.appendChild(username);
            li.appendChild(document.createTextNode(': '));
            li.appendChild(content);

            break;
        }
        if ($('#chatbox li').length >20) {
            $('#chatbox li').last().remove();
        }
        $('#chatbox li').first().before(li);
    };

    // Send messages.
    var postConecnt = function () {
        var uname = $('#uname').text();
        var content = $('#sendbox').val();
        socket.send(content);
        $('#sendbox').val('');
    }


    function  ChangeRoom(url,json){
        $.ajax({
            type: 'POST',
            url: url,
            contentType: "application/json; charset=utf-8",
            data: JSON.stringify(json),
            dataType: "json",
            success: function (message) {
                if (message > 0) {
                    alert("请求已提交！我们会尽快与您取得联系");
                }
            },
            error: function (message) {
                alert("提交数据失败！");
            }

        });
    }
    $('#sendbtn').click(function () {
        postConecnt();
    });
    $('#sendbox').keydown(function(event){
        if(event.keyCode ==13){
            postConecnt();
        }
    });

    $('#change_room').click(function () {
        ChangeRoomName();
        var path =GetUrlPath();
        // var  json ={
        //     "name": $('#uname').text(),
        //     "leftRoom": $('#room_name').text(),
        //     "joinRoomName": "NEW",
        //     "time": 11111111
        // };
        // alert(path)
        // ChangeRoom(path,json)
        alert($('#room_name').text());
        alert(path);
        window.location.href=path;

    });

    function ChangeRoomName()
    {
        $('#room_name').text( $('#room_name_input').val());
        $('#room_name_input').val("");
    }

    function GetUrlPath()
    {
        var url = document.location.toString();
        if(url.indexOf("?") != -1){
            url = url.split("?")[0];
        }
        // url=url+"/ChanfRoom";
        url=AddPath(url,"uname",$('#uname').text());
        url=AddPath(url,"room",$('#room_name').text());

        return url
    }

    function AddPath(url,name,val)
    {
        if (url.indexOf("?") > 0) {
            url = url + "&" + name + "=" + val;
        }
        else {
            url = url + "?" + name + "=" + val;
        }
        return url

    }
    function getUrlParam(name) {
        var reg = new RegExp("(^|&)" + name + "=([^&]*)(&|$)"); //构造一个含有目标参数的正则表达式对象
        var r = window.location.search.substr(1).match(reg); //匹配目标参数
        if (r != null) return unescape(r[2]); return ''; //返回参数值
    }


});
