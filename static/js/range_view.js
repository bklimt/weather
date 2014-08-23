
var RangeView = Backbone.View.extend({
  initialize: function(options) {
    this.render(options.attr);
  },

  render: function(attr) {
    var self = this;

    self.$el.on("change", function() {
      attrs = {}
      attrs[attr] = parseInt(self.$el.val());
      if (self.model.saveThrottled) {
        self.model.saveThrottled(attrs);
      } else {
        self.model.save(attrs);
      }
    });
  }
});

