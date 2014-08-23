
var LightView = Backbone.View.extend({
  initialize: function() {
    this.render();
  },

  render: function() {
    new RangeView({
      el: this.$("#hue"),
      model: this.model,
      attr: "hue"
    });

    new RangeView({
      el: this.$("#sat"),
      model: this.model,
      attr: "sat"
    });

    new RangeView({
      el: this.$("#bri"),
      model: this.model,
      attr: "bri"
    });

    new CheckView({
      el: this.$("#on"),
      model: this.model,
      attr: "on"
    });
  }
});

