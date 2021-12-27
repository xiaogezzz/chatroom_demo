<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title></title>
    <meta name="viewport" content="width=device-width,initial-scale=1,user-scalable=0">
    <link href="../static/stylesheets/weui-1.1.3.min.css" rel="stylesheet">
</head>
    <style>
        .weui-cells:after{border: 1px solid #fff;}
    </style>
<body>

<div style="max-width:640px; margin:0 auto;">
    <div class="weui-panel__hd" style="background-color: #1aad19; height:2.5em">
        <img style="float: left; width: 32px; height:32px;" src="../static/img/chatroom.png">
        <a style="float: right; width: 32px; height:32px;" href="/">
            <img src="../static/img/log-out.png">
        </a>
    </div>
    <div id="chat-column">
        <div id="chat-messages" style='overflow:auto;'></div>
    </div>
    <div id="input-box">
        <form id="msg_form" class="weui-cells weui-cells_form">
            <div class="weui-cell weui-cell_vcode">
                <div class="weui-cell__bd">
                    <textarea class="weui-textarea" placeholder="请输入文本" rows="3" id="appendedPrependedInput"></textarea>
                </div>
            </div>
            <div class="weui-cells weui-cells_form">
                <div class="weui-cell">
                    <button class="weui-btn weui-btn_primary" type="submit">发送</button>
                </div>
            </div>
        </form>
    </div>
</div>

<div id="msg-template" class="weui-panel__bd" style="display:none;">
    <div class="weui-media-box_appmsg" style="padding: 0.8em 0.8em 0 0.8em">
        <div class="weui-media-box__hd userpic" style="width: 32px; height: 32px; line-height: 60px;vertical-align:top">
            <img class="weui-media-box__thumb" style="border-radius:0.2em" src="" alt="">
        </div>
        <div class="weui-media-box__bd" style="font-size: 12px;">
            <p class="uid" style="color: #999;"></p>
            <p class="content" style="color:#444"></p>
            <p class="msg-time" style=" color: #999; text-align: right;"></p>
        </div>
    </div>
</div>

</body>
</html>
<script type="text/javascript" src="../static/js/jquery-1.11.1.min.js"></script>
<script>
    $(function(){
        var conn;

        $.getDocHeight = function(){
            return Math.max(
                $(document).height(),
                $(window).height(),
                /* For opera: */
                document.documentElement.clientHeight
            );
        };

        function chatResized() {
            var chat_width = document.body.clientWidth - 405;
            var chat_height = $.getDocHeight() - 230;
            if (chat_width > 0) {
                $('#chat-column, #input-box').css('max-width', 640);
                $("#appendedPrependedInput").css('max-width', 640);
            }
            $('#chat-column').css('height',chat_height);
            $('#chat-messages').css('height',chat_height);

        };

        chatResized();
        window.onresize = chatResized;

        function addMessage(textMessage) {

            $("#msg-template .userpic").html("<img class=\"weui-media-box__thumb\" style='border-radius:0.2em' src='" + textMessage.Gravatar + "'>")
            $("#msg-template .msg-time").html(textMessage.Time);
            $("#msg-template .uid").html(textMessage.Uid);
            $("#msg-template .content").html(textMessage.Content);
            $("#chat-messages").append($("#msg-template").html());
            $('#chat-column')[0].scrollTop = $('#chat-column')[0].scrollHeight;
            var div_chat_message = document.getElementById("chat-messages");
            div_chat_message.scrollTop = div_chat_message.scrollHeight
        };

        function updateUsers(textMessage) {
            $("#msg-template .userpic").html("")
            $("#msg-template .msg-time").html("");
            $("#msg-template .uid").html("");
            $("#msg-template .content").html(textMessage.Content);
            $("#chat-messages").append($("#msg-template").html());
            $('#chat-column')[0].scrollTop = $('#chat-column')[0].scrollHeight;
            var div_chat_message = document.getElementById("chat-messages");
            div_chat_message.scrollTop = div_chat_message.scrollHeight
        };

        function errorMessage(msg) {
            $("#msg-template .content").html(msg);
            $("#chat-messages").append($("#msg-template").html());
            $('#chat-column')[0].scrollTop = $('#chat-column')[0].scrollHeight;
        };

        $("#msg_form").submit(function() {
            var msg = $("#appendedPrependedInput");
            if (!conn) {
                return false;
            }
            if (!msg.val()) {
                alert("不能发送空白消息");
                return false;
            }
            conn.send(msg.val());
            msg.val("");
            return false
        });


        if (window["WebSocket"]) {
            conn = new WebSocket("ws://localhost:8080/chat?uid={{.uid}}");
            conn.onopen = function() {
                console.log('Connection open ...');
            };

            conn.onmessage = function(evt) {
                console.log(evt.data)

                var data = JSON.parse(evt.data);
                switch (data.Type) {
                    case "text_type":
                        addMessage(data)
                        break;
                    case "status_type":
                        updateUsers(data)
                        break;
                    default:
                        console.log(data)
                }
            };

            conn.onerror = function() {
                errorMessage("<strong> An error just occured.<strong>")
            };

            conn.onclose = function() {
                errorMessage("<strong>Connection closed.<strong>")
            };
        } else {
            errorMessage("Your browser does not support WebSockets.");
        }
    });
</script>