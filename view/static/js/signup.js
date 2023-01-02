function signup(){
    let username=document.getElementById("username").value;
    let password=document.getElementById("password").value;
    let xmlHttp=new XMLHttpRequest();
    xmlHttp.onreadystatechange=function (){
        if(xmlHttp.readyState==4&&xmlHttp.status==200){
            let data=JSON.parse(this.responseText);
            processData(data);
        }
        if (xmlHttp.readyState==4&&xmlHttp.status==503){
            alertError();
        }
    }
    xmlHttp.open("POST","/signup");
    xmlHttp.setRequestHeader("Content-Type","application/json");
    xmlHttp.send(JSON.stringify({
        "username":username,
        "password":password
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