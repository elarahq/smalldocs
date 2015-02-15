(function(root){
    // get app
    var app = root.app;

    // dispatch route
    var dispatchURL = function(projectName, topicName, pageName) {
        app.dispatcher.dispatch({
            actionType: "changeURL",
            projectName: projectName,
            topicName: topicName,
            pageName: pageName
        });
    }

    // routes
    var routes = {
        "/docs/:projectName": dispatchURL,
        "/docs/:projectName/:topicName/:pageName": dispatchURL,
        "/edit/:projectName/:topicName/:": dispatchURL,
        "/edit/:projectName/:topicName/:pageName": dispatchURL,
        "/edit/:projectName/:topicName/:pageName": dispatchURL,
    };

    // add router
    app.router = Router(routes);
    app.router.configure({
        html5history: true,
    });
    app.router.init();

    // No reload for links
    $('#docs').on("click", "a[data-noreload]", function(event) {
        // Prevent default click action
        event.preventDefault();
        app.router.setRoute($(this).attr('href'));
        return false;
    });
})(this);
