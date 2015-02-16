(function(root){
    // Get app
    var app = root.app;

    // Page
    var Page = app.Page = React.createClass({
        displayName: "Page",

        getInitialState: function() {
            return {
                active: this.props.active
            };
        },

        componentWillMount: function(){
            this.dispatchToken = app.dispatcher.register(function(payload){
                switch(payload.actionType) {
                    case "change:url":
                        this.state.active = payload.currentPage == this.props.page.name && payload.currentTopic == this.props.topic.name;
                        this.setState(payload);
                        break;
                }
            }.bind(this));
        },

        componentWillUnmount: function(){
            app.dispatcher.unregister(this.dispatchToken);
        },

        componentDidMount: function(){
            if (this.state.active) {
                app.dispatcher.dispatch({
                    actionType: "change:page",
                    pageId: this.props.page.id,
                });
            }
        },

        changeMarkdown: function(e){
            e.preventDefault();

            app.router.setRoute($(e.target).attr("href"));
            app.dispatcher.dispatch({
                actionType: "change:page",
                pageId: this.props.page.id,
            });
        },

        render: function() {
            var page = this.props.page;
            var source = this.props.source;

            var cn = ["page", this.state.active ? "active" : ""].join(' ');
            return <div className={cn} data-page={page.id} data-name={page.name}>
                    <div className="page-text" data-page={page.id} data-name={page.name}>
                        <a href={source} onClick={this.changeMarkdown}>
                            {page.title}
                        </a>
                    </div>
                </div>
        }
    });
})(this);
