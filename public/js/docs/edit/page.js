(function(root){
    // Get app
    var app = root.app;

    // Page
    var Page = app.Page = React.createClass({
        displayName: "Page",

        getInitialState: function() {
            return {
                active: false
            };
        },

        componentWillMount: function(){
            this.dispatchToken = app.dispatcher.register(function(payload){
                switch(payload.actionType) {
                    case "change:url":
                        this.state.active = payload.currentPage == this.props.page.name;
                        this.setState(payload);
                        break;
                }
            }.bind(this));
        },

        componentWillUnmount: function(){
            app.dispatcher.unregister(this.dispatchToken);
        },

        render: function() {
            var page = this.props.page;
            var source = this.props.source;

            return <div className="page" data-page={page.id} data-name={page.name}>
                    <div className="page-text" data-page={page.id} data-name={page.name}>
                        <a href={source} data-noreload="true">
                            {page.title}
                        </a>
                    </div>
                </div>
        }
    });
})(this);
