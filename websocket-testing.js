// I need some kind of comment here because otherwise the line below looks odd to me
// sorry..

function onmessage(event) {
  console.log(event.data);
  var data = JSON.parse(event.data);
  switch (data.type) {
    case "error":
      handleError(data.error);
      break;
    case "status":
      handleStatus(data.text);
      getTopics();
      break;
    case "message":
      var chat = $("#chat");
      var topic = null;
      $.each(topics, function(key, value) {
        if (value.id == data["topic-id"]) {
          console.log("match!");
          topic = value;
        }
      });
      
      var topicName = "";
      if (topic === null) {
        topicName = "GLOBAL";
      } else {
        console.log(topic);
        topicName = topic["option-a"] + ' vs. ' + topic["option-b"];
      }
      chat.append('[' + topicName + ']' +
        '<span style="color: #'+data.color+';">'+data.username+'</span>:  '+data.text+'<br>');
      break;
  }
}

function handleError(error) {
  $("#chat").append('[<span style="color: red;">ERROR</span>] ' + error + '<br>'); 
}
function handleStatus(text) {
  $("#chat").append('[<span style="color: cyan;">SERVER</span>] ' + text + '<br>');
}

function onerror(err) {
  console.log(err);
}

var socket = new WebSocket("ws://localhost:8080/ws");
var topics = {};

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
  var sid = $("#topics-sel").val();
  var id = parseInt(sid);
  var data = JSON.stringify({type:"message", text:text, "topic-id": id});
  socket.send(data);
}

function leave() {
  var data = JSON.stringify({type:"leave"});
  socket.send(data);
}

function newTopic() {
  var a = $("#option-a").val();
  var b = $("#option-b").val();

  var data = {type:"newtopic"};
  data["option-a"] = a;
  data["option-b"] = b;
  socket.send(JSON.stringify(data));
}

function getTopics() {
  var req = new XMLHttpRequest();
  req.open("GET", "http://localhost:8080/status", true);
  req.send(null);

  req.onload = function() { 
    var data = JSON.parse(req.responseText);
    topics = data.topics;
    console.log(req.responseText);
    topicsElem = document.getElementById("topics");
    topicsElem.innerHTML = "";
    for (var i = 0; i < data.topics.length; i++) {
      topicsElem.innerHTML = topicsElem.innerHTML + '<p>ID: ' + data.topics[i].id + ' ' + data.topics[i]["option-a"] + ' vs. ' + data.topics[i]["option-b"];
    }
    $("#topics-sel").find("option").remove();
    $("#topics-sel").append('<option value="" disabled selected>Select your option</option>');
    $.each(data.topics, function(key, value) {
      $("#topics-sel").append('<option value="'+value.id+'">'+
          value["option-a"] +" vs. " + value["option-b"])
    })
  }
}
