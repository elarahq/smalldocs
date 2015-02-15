(function(root){
    root.app = root.app || {};

    // Get app
    var app = root.app;
    // Flux dispatcher
    if (typeof window.Flux != "undefined") {
        var dispatcher = app.dispatcher = new Flux.Dispatcher();
    }

    // Dispatch event for currentTopic and currentPage
    app.dispatchInfo = function(){
        if (!app.dispatcher) return;

        var pathname = window.location.pathname;
        var tokens = pathname.split("/").filter(function(p){
            return p != "";
        });

        var obj = {
            actionType: "changeURL",
            isEdit: tokens[0] == "edit",
            isView: tokens[0] == "docs",
            currentProject: tokens[1],
            currentTopic: tokens[2],
            currentPage: tokens[3]
        }

        // Dispatch URL information
        app.dispatcher && app.dispatcher.dispatch(obj);
    }

    // on pop state
    window.onpopstate = app.dispatchInfo;

    // No reload for links
    $('#docs').on("click", "a[data-noreload]", function(event) {
        // Prevent default click action
        event.preventDefault();

        // Detect if pushState is available
        if(history.pushState) {
            history.pushState(null, null, $(this).attr('href'));
        }
        root.app.dispatchInfo();
        return false;
    });

})(this);
