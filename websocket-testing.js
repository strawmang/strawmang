// I need some kind of comment here because otherwise the line below looks odd to me
// sorry..

function onmessage(event) {
  console.log(event.data);
  var data = JSON.parse(event.data);
  switch (data.type) {
    case "message":
      var area = document.getElementById("chat");
      area.innerHTML = area.innerHTML + '<span style="color: #'+data.color+';">'+data.username+'</span>:  '+data.text+'<br>';
      break;
    case "status":
      break;
  }
}

function onerror(event) {
  console.log(event.data);
}

var socket = new WebSocket("ws://localhost:8080/ws");

socket.onmessage = onmessage;
socket.onerror = onerror;

$(window).on('beforeunload', function(){
  socket.close();
});


function connect() {
  socket = new WebSocket("ws://localhost:8080/ws");
  socket.onmessage = onmessage;
  socket.onerror = onerror;
}

function login() {
  var name = document.getElementById("username").value;
  var data = JSON.stringify({type:"login", username:name});
  console.log("Sending: "+data);
  socket.send(data);
}

function sendMessage() {
  switch (socket.readyState) {
  case 0:
    console.log("That connetion isn't open yet!");
    return;
  case 1:
    break;
  case 2:
    console.log("That socket is closing!");
    return;
  case 3:
    console.log("That socket is closed!");
    return;
  }
  var text = document.getElementById("text").value;
  var data = JSON.stringify({type:"message", text:text});
  socket.send(data);
}

function leave() {
  var data = JSON.stringify({type:"leave"});
  socket.send(data);
}

function getTopics() {
  var req = new XMLHttpRequest();
  req.open("GET", "http://localhost:8080/status", true);
  req.send(null);

  req.onload = function() {
    var data = JSON.parse(req.responseText);
    console.log(req.responseText);
    var topics = document.getElementById("topics");
    topics.innerHTML = "";
    for (var i = 0; i < data.topics.length; i++) {
      topics.innerHTML = topics.innerHTML + '<p>ID: ' + data.topics[i].id + ' ' + data.topics[i].option-a + ' vs. ' + data.topics[i].option-b;
    }
  }
}
