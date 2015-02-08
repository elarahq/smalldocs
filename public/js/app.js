(function(root) {
    // App object
    var app = root.app = {
        // All data node store
        nodeStore: [],
        // Views
        views: {},
    };

    // Flux dispatcher
    var dispatcher = app.dispatcher = new Flux.Dispatcher();
    // add callback to dispatcher
    dispatcher.register(function(payload) {
        switch(payload.actionType){
            case "reset:nodes":
                app.nodeStore = payload.nodes;
                break;
        }
    });

})(this);
