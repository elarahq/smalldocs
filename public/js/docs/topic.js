(function(root){
    // Get app
    var app = root.app;

    // Topic
    var Topic = app.Topic = React.createClass({
        displayName: "Topic",

        getInitialState: function() {
            return {
                collapsed: this.props.collapsed
            };
        },

        toggleCollapsed: function() {
            this.setState({
                collapsed: !this.state.collapsed
            });
        },

        render: function() {
            // Get Page component
            var Page = app.Page;
            var topic = this.props.topic;
            var cViews = (topic.children || []).map(function(page, key){
                var source = ["/docs", topic.name, page.name].join("/");
                return <Page key={key} source={source} topic={topic} page={page}/>
            });
            var cn = ["topic", this.state.collapsed ? "collapsed" : ""].join(' ');
            return <div className={cn} data-topic={topic.id} data-name={topic.name}>
                    <div className="topic-text" data-topic={topic.id} data-name={topic.name}
                        onClick={this.toggleCollapsed}>
                        {topic.title}
                    </div>
                    <div className="page-list">{cViews}</div>
                </div>
        }
    });
})(this);
