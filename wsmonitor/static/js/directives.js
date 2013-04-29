'use strict';

/* Directives */


angular.module('myApp.directives', []).
directive('chart', function() {
    return {
        restrict: 'E',
        link: function(scope, elem, attrs){

            var chart = null,
                opts  = { };

            scope.$watch(attrs.ngModel, function(v){
                if(!chart){
                    chart = $.plot(elem, v , opts);
                    elem.show();
                }else{
                    chart.setData(v);
                    chart.setupGrid();
                    chart.draw();
                }
            },true);
        }
    }


});
