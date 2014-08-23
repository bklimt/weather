
var Light = Backbone.Model.extend({
  initialize: function() {
    this.previousSaves = $.Deferred().resolve().promise();
    this.cancellationToken = {};
  },

  saveThrottled: function(attrs) {
    var self = this;

    // Cancel any previous save that hasn't started.
    self.cancellationToken.cancelled = true;

    var token = {};
    self.cancellationToken = token;

    self.previousSaves = self.previousSaves.then(function() {
      if (token.cancelled) {
        return;
      }

      return self.save(attrs);
    });
  }
});

var Lights = Backbone.Collection.extend({
  model: Light,
  url: "/light",
  comparator: "id"
});

var lights = new Lights();

