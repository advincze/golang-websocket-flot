'use strict';

/* Services */


// Demonstrate how to register services
// In this case it is a simple value service.
var myModule = angular.module('myApp.services', [])

myModule.factory('webSocketFactory', function($rootScope) {
    var wsSrv = {}

    wsSrv.closeHandler = function () {
        console.log("websocket closed")
    };

    wsSrv.messageHandler = function(msg){
        console.log("message received: "+msg)
    }

    var ws = new WebSocket("ws://"+location.hostname+":"+location.port + "/ws");

    ws.onopen = function () {
        console.log("websocket opened");
    };

    ws.onmessage = function (msg) {
        //console.log(msg)
        $rootScope.$apply(function(){
            wsSrv.messageHandler(msg.data);
        })
    };

    ws.onclose = wsSrv.closeHandler();

    return wsSrv;

});


myModule.value('version', '0.1');




