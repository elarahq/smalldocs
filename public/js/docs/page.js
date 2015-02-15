(function(root){
    // Get app
    var app = root.app;

    // Page
    var Page = app.Page = React.createClass({
        displayName: "Page",

        getInitialState: function() {
            return {
                currentPage: null,

                page: this.props.page,
                source: this.props.source,
                topic: this.props.topic,
            };
        },

        componentWillMount: function(){
        },

        render: function() {
            var page = this.state.page;
            var source = this.state.source;
            var active = page.name != this.state.currentPage;

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
