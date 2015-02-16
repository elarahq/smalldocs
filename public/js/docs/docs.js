(function(root){
    // Get app
    var app = root.app;

    // Docs view
    var Docs = app.Docs = React.createClass({
        displayName: "Docs",

        getInitialState: function(){
            return {};
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
            var projectName = $("body").data("name");
            var projectId = $("body").data("id");

            var topicUrl = this.props.source.replace("ID", projectId);
            // Get TopicList component
            var TopicList = app.TopicList;

            if (!this.state.currentProject) {
                return <div></div>
            }

            return <div>
                    <div className="col-sm-4 col-md-3">
                        <TopicList source={topicUrl}
                            projectName={this.state.currentProject}
                            topicName={this.state.currentTopic}
                            pageName={this.state.currentPage}
                            projectId={projectId}/>
                    </div>
                    <div className="col-sm-8 col-md-9">
                    </div>
                </div>;
        }
    });

    // Load docs
    React.render(
        <Docs source="/projects/ID/topics"/>,
        document.getElementById('docs')
    );
})(this);
