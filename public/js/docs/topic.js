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

        fetchPages: function() {
            return $.ajax({
                url: this.props.pages,
                method: "GET"
            }).success(function(result){
                this.props.topic.children = result;
                this.forceUpdate();
            }.bind(this));
        },

        toggleCollapsed: function() {
            if (!this.props.topic.children) {
                this.fetchPages();
            }

            this.setState({
                collapsed: !this.state.collapsed
            });
        },

        componentWillMount: function(){
            if (!this.state.collapsed) {
                this.fetchPages();
            }

            this.dispatchToken = app.dispatcher.register(function(payload){
                switch(payload.actionType) {
                    case "change:url":
                        this.state.collapsed = payload.currentTopic != this.props.topicName;
                        this.setState(payload);
                        break;
                }
            }.bind(this));
        },

        componentWillUnmount: function(){
            app.dispatcher.unregister(this.dispatchToken);
        },

        render: function() {
            // Get Page component
            var Page = app.Page;
            var topic = this.props.topic;
            var projectName = this.props.projectName;
            var cViews = (topic.children || []).map(function(page, key){
                var source = ["/docs", projectName, topic.name, page.name].join("/");
                return <Page key={key}
                    projectName={this.props.projectName}
                    topicName={this.props.topicName}
                    pageName={this.props.pageName}
                    source={source}
                    active={this.props.pageName == page.name && this.props.topicName == topic.name}
                    topic={topic}
                    page={page}/>
            }.bind(this));

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
