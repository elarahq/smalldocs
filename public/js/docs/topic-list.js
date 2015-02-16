(function(root){
    // Get app
    var app = root.app;

    // Topic
    var TopicList = app.TopicList = React.createClass({
        displayName: "TopicList",

        getInitialState: function() {
            return {
                topics: null,
            };
        },

        fetch: function(){
            return $.ajax({
                url: this.props.source,
                method: "GET"
            });
        },

        componentDidMount: function() {
            if (this.props.source) {
                this.fetch().success(function(result){
                    if (this.isMounted()) {
                        this.setState({
                            topics: result
                        });
                    }
                }.bind(this));
            }
        },

        render: function() {
            // Get Topic component
            var Topic = app.Topic;

            var projectName = this.props.projectName;
            var projectId = this.props.projectId;

            var topics = this.state.topics || [];
            var views = topics.map(function(topic, key) {
                var source = ["/docs", projectName, topic.name].join("/");
                var pages = ["/projects", projectId, "topics", topic.id, "pages"].join("/");
                return <Topic key={key}
                        projectName={this.props.projectName}
                        topicName={this.props.topicName}
                        pageName={this.props.pageName}
                        pages={pages}
                        topic={topic}
                        collapsed={this.props.topicName != topic.name}
                        source={source}/>
            }.bind(this));
            return <div className="topic-list">{views}</div>
        }
    });
})(this);
