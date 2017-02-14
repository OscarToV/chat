new Vue({
    el: '#app',

    data: {
        ws: null, // Nuestro websocket
        newMsg: '', // mantiene nuevos mensajes para ser enviados
        chatContent: '', // lista de los mensajes desplegados en el chat
        avatar: null, // avatar
        username: null, //username
        joined: false // se pone en verdadero si se ingresa un username valido
    },

    created: function() {
        var self = this;
        this.ws = new WebSocket('ws://' + window.location.host + '/ws');
        this.ws.addEventListener('message', function(e) {
            var msg = JSON.parse(e.data);
            self.chatContent += '<div class="chip">'
                    + '<img src="' + self.gravatarURL(msg.username) + '">' // Avatar
                    + msg.username
                + '</div>'
                + emojione.toImage(msg.message) + '<br/>'; // Parse emojis

            var element = document.getElementById('chat-messages');
            element.scrollTop = element.scrollHeight; // Auto scroll
        });
    },

    methods: {
        send: function () {
            if (this.newMsg != '') {
                this.ws.send(
                    JSON.stringify({
                        username: this.username,
                        message: $('<p>').html(this.newMsg).text()
                    }
                ));
                this.newMsg = '';
            }
        },

        join: function () {
            if (!this.username) {
                Materialize.toast('Debes seleccionar un username', 2000);
                return
            }
            this.username = $('<p>').html(this.username).text();
            this.joined = true;
        },

        gravatarURL: function(avatar) {
            return 'http://www.gravatar.com/avatar/' + CryptoJS.MD5(avatar);
        }
    }
});
