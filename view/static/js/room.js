var ws=null;
//var uri="ws://localhost/room/ws" // reverse proxy
// var uri="ws://localhost:8080/room/ws"
var userId,username,roomId;
function wsConnect(){
    userId=document.getElementById("userId").value;
    username=document.getElementById("username").value;
    roomId=document.getElementById("roomId").value;
    let uri="ws://localhost:8080/room/ws"+"?room_id="+roomId+"&user_id="+userId;
    ws=new WebSocket(uri);
    ws.onopen=function (){
        console.log("connected to "+uri+" at "+Date());
        let msg=JSON.stringify({
            "type":0,
            "data":{
                "room_id":roomId,
                "user_id":userId,
                "username":username
            }
        });
        console.log(msg)
        ws.send(msg);
    }
    ws.onclose=function (e){
        console.log("connection closed ("+e.code+"): "+e.reason);
    }
    ws.onmessage=function (e){
        console.log("message received: "+e.data);
        let msg=JSON.parse(e.data);
        console.log("message type",msg.type);
        let data=msg.data;
        let div0;
        let div1;
        switch (msg.type){
            case 0: // msgTypeOnline
                div0='<div class="public text-center">'+data.username+" enters into the room"+'</div>';
                div1=$('<div class="list-group-item border-0">').append(div0);
                $("#messageArea").append(div1);
                break;
            case 1:
                if (userId==data.user_id){
                    let head='<span class="head">'+data.username+':'+'</span>';
                    let body='<span class="body bg-primary text-light rounded ps-0">'+data.content+'</span>';
                    div0=$('<div class="row send">').append(head,body);
                }else{
                    //let a='<a class="text-dark" style="text-decoration: none" href="">'+data.username+':'+'</a>';
                    //'<a class="text-dark" style="text-decoration: none" href="" userId="123" username="developer" onclick="open(this)"'
                    //let head=$('<span class="head">').append(a);
                    let headlink=$('<a class="text-dark" style="text-decoration: none" href="javascript:void(0);" onclick="openProfile(this);return false;">').append(data.username+':');
                    headlink.attr("userId",data.user_id);
                    headlink.attr("username",data.username);
                    let head=$('<span class="head">').append(headlink);
                    let body='<span class="body ps-0">'+data.content+'</span>';
                    div0=$('<div class="row receive">').append(head,body);
                }
                div1=$('<div class="list-group-item border-0">').append(div0);
                $("#messageArea").append(div1);
                break;
        }
        let m=document.getElementById("messageArea");
        m.scrollTop=m.scrollHeight;
    }
}
function wsClose(){
    let msg=JSON.stringify({
        "type":3,
        "data":{
            "user_id":userId,
            "room_id":roomId,
            "username":username
        }
    });
    ws.send(msg);
    ws.close();
}
function send(){
    let content=$("#inputArea").val();
    $("#inputArea").val("");
    $("#sendBtn").attr("disabled",true);
    // if toUserId is not empty, it is a private chat message
    let url=window.location.href;
    let params=new URLSearchParams(url.split('?')[1]);
    let msg=JSON.stringify({
        "type":1,
        "data":{
            "user_id":userId,
            "username":username,
            "room_id":roomId,
            "to_user_id":params.get('to_user_id'),
            "content":content
        }
    });
    ws.send(msg);
}
function preprocess(){
    // preprocess
    $("#inputArea").on("input propertychange",function (){
        if($("#inputArea").val()!=""){
            $("#sendBtn").attr("disabled",false);
        }else{
            $("#sendBtn").attr("disabled",true);
        }
    });
    wsConnect();
}
function postprocess(){
    wsClose();
}
window.onload=preprocess;
window.onunload=postprocess;

function openProfile(obj){
    let toUserId=obj.getAttribute("userId");
    let toUsername=obj.getAttribute("username");
    let userId=document.getElementById("userId").value;
    console.log(toUserId,toUsername,userId);
    $('#profileModal').modal("show");
    let toUserIdElem=$('#toUserId');
    toUserIdElem.text(toUserId);
    $('#profileUsername').text(toUsername);
}

function privateChat(){
    let userId=document.getElementById("userId").value;
    let toUserId=$('#toUserId').text();
    let toUsername=$('#profileUsername').text();
    let postRequest=new XMLHttpRequest();
    postRequest.onreadystatechange=function (){
        if(postRequest.readyState==4&&postRequest.status==200){
            console.log(this.responseText);
            let data=JSON.parse(this.responseText);
            processData(data);
        }
        if (postRequest.readyState==4&&postRequest.status==503){
            console.log("server errors");
        }
    }
    postRequest.open("POST","/private-chat");
    postRequest.setRequestHeader("Content-Type","application/x-www-form-urlencoded");
    postRequest.send("user_id="+userId+"&to_user_id="+toUserId+"&to_username="+toUsername);
}

function processData(data){
    console.log(data.room_id);
    let roomId=data.room_id;
    let toUsername=$('#profileUsername').text();
    let toUserId=$('#toUserId').text();
    console.log(toUsername);
    window.open("/room/"+roomId+"?room_name="+toUsername+"&to_user_id="+toUserId);
}