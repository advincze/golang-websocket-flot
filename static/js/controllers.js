'use strict';

/* Controllers */

function ChartCtrl($scope, wsSrv) {
    console.log("hello world")
    $scope.error = {msg:""}
    $scope.messages = [];

    $scope.counter = 3;
    $scope.chart = {}
    $scope.chart.data = [[]];

    wsSrv.messageHandler = function(msg){
        var obj = JSON.parse(msg);
        $scope.chart.data[0].push([$scope.counter++,obj.IntValue])
        if($scope.chart.data[0].length > 100){
            $scope.chart.data[0].shift();
        }
    }

}
ChartCtrl.$inject = ['$scope', 'webSocketFactory'];
