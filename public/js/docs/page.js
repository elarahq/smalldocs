(function(root){
    // Get app
    var app = root.app;

    // Page
    var Page = app.Page = React.createClass({
        displayName: "Page",

        getInitialState: function() {
            return {
                page: this.props.page,
                source: this.props.source,
                topic: this.props.topic,
            };
        },

        componentDidMount: function() {
        },

        render: function() {
            var page = this.state.page;
            return <div className="page" data-page={page.id} data-name={page.name}>
                    <div className="page-text" data-page={page.id} data-name={page.name}>
                        <a href={this.props.source}>
                            {page.title}
                        </a>
                    </div>
                </div>
        }
    });
})(this);
