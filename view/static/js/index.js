function search(){
    let roomName=document.getElementById("room_name").value;
    let xmlHttp=new XMLHttpRequest();
    xmlHttp.onreadystatechange=function (){
        if(xmlHttp.readyState==4&&xmlHttp.status==200){
            let data=JSON.parse(this.responseText);
            processData(data);
        }
    }
    xmlHttp.open("POST","/index/search");
    xmlHttp.setRequestHeader("Content-Type","application/x-www-form-urlencoded");
    xmlHttp.send("room_name="+roomName);
}

function processData(data){
    let len=data.length;
    for (let i=0;i<len;i++){
        let roomId=data[i]["ID"];
        let roomName=data[i]["room_name"];
        let username=data[i]["username"];
        let createdAt=new Date(data[i]["CreatedAt"]).toLocaleDateString();
        console.log(roomId,roomName,username,createdAt);

        let href="room/"+roomId;
        let title='<h5 class="mb-0">'+roomName+'</h5>';
        let owner='<p class="mb-0 opacity-75">'+username+'</p>';
        let div0=$('<div>').append(title,owner);
        let date='<small class="opacity-50 text-nowrap">'+createdAt+'</small>'
        let div1=$('<div class="d-flex gap-2 w-100 justify-content-between">').append(div0,date);
        let div2=$('<a href='+href+' class="list-group-item list-group-item-action d-flex gap-3 py-3">').append(div1)
        $("#searchResult").append(div2)
    }
}