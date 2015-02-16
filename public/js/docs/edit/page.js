(function(root){
    // Get app
    var app = root.app;

    // Page
    var Page = app.Page = React.createClass({
        displayName: "Page",

        render: function() {
            var page = this.props.page;
            var source = this.props.source;

            return <div className="page-edit" data-page={page.id} data-name={page.name}>
                    <div className="page-text" data-page={page.id} data-name={page.name}>
                        <a href={source} data-noreload="true">
                            {page.title}
                        </a>
                    </div>
                </div>
        }
    });
})(this);
