var ws=null;
var uri="ws://localhost:8080/room/ws"
var userId,username,roomId;
function wsConnect(){
    userId=document.getElementById("userId").value;
    username=document.getElementById("username").value;
    roomId=document.getElementById("roomId").value;
    ws=new WebSocket(uri);
    ws.onopen=function (){
        console.log("connected to "+uri+" at "+Date());
        let msg=JSON.stringify({
            "msgType":1,
            "data":{
                "username":username,
                "roomId":roomId
            }
        });
        ws.send(msg);
    }
    ws.onclose=function (e){
        console.log("connection closed ("+e.code+"): "+e.reason);
    }
    ws.onmessage=function (e){
        console.log("message received: "+e.data);
        let msg=JSON.parse(e.data);
        console.log("message type",msg.msgType);
        let data=msg.data;
        switch (msg.msgType){
            case 1: // msgTypeOnline
                let v=data.username+" enters into the room "+data.roomId;
                $("#messageArea").append('<li class="public">'+v+'</li>');
                break;
            case 2:
                let msgHead='<span class="head">'+data.username+'(room '+data.roomId+'):'+'</span><br>';
                let msgBody='<span class="body">'+data.content+'</span>';
                let msgItem;
                if(username==data.username){ // sender
                    msgItem=$('<li class="send"/>').append(msgHead,msgBody);
                }else{ // receiver
                    msgItem=$('<li class="receive"/>').append(msgHead,msgBody);
                }
                $("#messageArea").append(msgItem);
                break;
        }
        let m=document.getElementById("messageArea");
        m.scrollTop=m.scrollHeight;
    }
}
function wsClose(){
    let msg=JSON.stringify({
        "msgType":3,
        "data":{
            "username":username,
            "roomId":roomId
        }
    });
    ws.send(msg);
    ws.close();
}
function sendMsg(){
    let content=$("#inputArea").val();
    $("#inputArea").val("");
    $("#sendBtn").attr("disabled",true);
    let msg=JSON.stringify({
        "msgType":2,
        "data":{
            "userId":userId,
            "roomId":roomId,
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
