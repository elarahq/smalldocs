(function(root){
    // Get app
    var app = root.app;

    // Docs view
    var Docs = app.Docs = React.createClass({
        displayName: "Docs",

        getInitialState: function(){
            return {
                currentProject: null,
                currentTopic: null,
                currentPage: null,
            };
        },

        componentWillMount: function(){
            this.dispatchToken = app.dispatcher.register(function(payload){
                switch(payload.actionType) {
                    case "change:url":
                        this.setState(payload);
                        break;
                }
            }.bind(this));
        },

        componentWillUnmount: function(){
            app.dispatcher.unregister(this.dispatchToken);
        },

        render: function() {
            var projectId = $("body").data("id");
            var projectName = $("body").data("name");

            if (!this.state.currentProject) {
                return <div></div>
            }

            // source
            var topics = ["/projects", this.state.currentProject, "topics"].join("/");

            var p = <app.TopicList
                    source={topics}
                    post={topics}
                    projectName={projectName}
                    projectId={projectId}/>;

            if (this.state.currentTopic) {
                if (this.state.currentPage) {
                    p = <app.Page projectName={projectName} projectId={projectId}/>;
                } else {
                    var source = [topics, this.state.currentTopic].join("/");
                    var pages = [topics, this.state.currentTopic, "pages"].join("/");

                    p = <app.Topic
                        source={source}
                        put={source}
                        remove={source}
                        post={pages}
                        pages={pages}
                        projectName={projectName}
                        projectId={projectId}
                        topicId={this.state.currentTopic}/>;
                }
            }

            return <div className="padding-top-10">
                    <div className="col-sm-6 col-sm-offset-3 padding-top-10">
                        {p}
                    </div>
                </div>;
        }
    });

    // Load docs
    React.render(
        <Docs/>,
        document.getElementById('docs-edit')
    );
})(this);
