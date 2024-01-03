let socket = new WebSocket("ws://localhost:8080/ws");

socket.onmessage = (event) => {
    updateList(event.data);
};

function updateList(data){
    var ul = document.getElementById("the-list");
    var li = document.createElement("li");
    li.appendChild(document.createTextNode(data));
    ul.appendChild(li);
};

function send(){
    socket.send(document.getElementById('msg-input').value);}