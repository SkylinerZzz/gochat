function newRoom(){
    let roomName=document.getElementById("room_name").value;
    let userId=document.getElementById("userId").value;
    let username=document.getElementById("username").value;
    let xmlHttp=new XMLHttpRequest();
    xmlHttp.onreadystatechange=function (){
        if(xmlHttp.readyState==4&&xmlHttp.status==200){
            console.log(this.responseText);
            let data=JSON.parse(this.responseText);
            if (data["msg"]=="success"){
                console.log("good");
                window.location.href="/index";
            }else{
                processData(data);
            }
        }
        if (xmlHttp.readyState==4&&xmlHttp.status==503){
            alertError();
        }
    }
    xmlHttp.open("POST","/index/new");
    xmlHttp.setRequestHeader("Content-Type","application/json");
    xmlHttp.send(JSON.stringify({
        "room_name":roomName,
        "username":username,
        "user_id":userId
    }));
}

function processData(data){
    if (data["msg"]!=""){
        let msg=data["msg"];
        $("#btnCloseAlert").click();
        let btn='<button type="button" class="btn-close" data-bs-dismiss="alert" id="btnCloseAlert"></button>';
        let text='<strong>'+msg+'</strong>';
        let div0=$('<div class="alert alert-danger alert-dismissible fade show">').append(btn,text);
        $("#alertDock").append(div0);
        window.setTimeout("closeAlert()",3000)
    }
}

function closeAlert(){
    $("#btnCloseAlert").click();
}

function alertError(){
    $("#btnCloseAlert").click();
    let btn='<button type="button" class="btn-close" data-bs-dismiss="alert" id="btnCloseAlert"></button>';
    let text='<strong>Something wrong with server</strong>';
    let div0=$('<div class="alert alert-danger alert-dismissible fade show">').append(btn,text);
    $("#alertDock").append(div0);
    window.setTimeout("closeAlert()",3000)
}