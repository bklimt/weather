
var Volume = Backbone.Model.extend({
  isNew: function() { return false; },
  url: "/volume"
});

var volume = new Volume();

