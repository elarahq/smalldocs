(function(root){
    root.app = root.app || {};

    // Get app
    var app = root.app;
    // Flux dispatcher
    if (typeof window.Flux != "undefined") {
        var dispatcher = app.dispatcher = new Flux.Dispatcher();
    }
})(this);
