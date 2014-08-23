
var CheckView = Backbone.View.extend({
  initialize: function(options) {
    this.render(options.attr);
  },

  render: function(attr) {
    var self = this;

    self.model.on('change:' + attr, function() {
      self.$el.prop("checked", self.model.get(attr));
    });

    self.$el.on("change", function() {
      attrs = {}
      attrs[attr] = self.$el.prop("checked");
      self.model.save(attrs);
    });
  }
});

