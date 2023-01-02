function login(){
    let username=document.getElementById("username").value;
    let password=document.getElementById("password").value;
    let postRequest=new XMLHttpRequest();
    postRequest.onreadystatechange=function (){
        if(postRequest.readyState==4&&postRequest.status==200){
            console.log(this.responseText);
            let data=JSON.parse(this.responseText);
            if (data["msg"]=="success"){
                console.log("good");
                window.location.href="/index";
            }else{
                processData(data);
            }
        }
        if (postRequest.readyState==4&&postRequest.status==503){
            alertError();
        }
    }
    postRequest.open("POST","/login");
    postRequest.setRequestHeader("Content-Type","application/json");
    postRequest.send(JSON.stringify({
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