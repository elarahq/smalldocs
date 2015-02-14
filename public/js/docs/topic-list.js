(function(root){
    // Get app
    var app = root.app;

    // Topic
    var TopicList = app.TopicList = React.createClass({
        displayName: "TopicList",

        getInitialState: function() {
            return {
                topics: null
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

            var topics = this.state.topics || [];
            var views = topics.map(function(topic, key) {
                var source = ["/docs", topic.name];
                var collapsed = true;
                return <Topic key={key} topic={topic} source={source} collapsed={collapsed}/>
            });
            return <div className="topic-list">{views}</div>
        }
    });
})(this);
