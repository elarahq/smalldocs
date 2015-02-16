(function(root){
    // Get app
    var app = root.app;

    // Topic
    var Topic = app.Topic = React.createClass({
        displayName: "Topic",

        getInitialState: function() {
            return {
                pages: null
            };
        },

        // fetch topic
        fetch: function(){
            return $.ajax({
                url: this.props.source,
                method: "GET"
            }).success(function(result) {
                if (this.isMounted()) {
                    this.setState(result);
                }
            }.bind(this));
        },

        // Fetch pages
        fetchPages: function() {
            return $.ajax({
                url: this.props.pages,
                method: "GET"
            }).success(function(result) {
                if (this.isMounted()) {
                    this.setState({
                        pages: result
                    });
                }
            }.bind(this));
        },

        componentDidMount: function(){
            this.fetch();
            this.fetchPages();
        },

        titleChange: function(e){
            var value = e.target.value || "";
            this.setState({
                title: value.replace(new RegExp("[^a-zA-Z0-9 \.-]", "gi"), "")
            })
        },

        // Save topic
        saveTopic: function(){
            var xhr = this.saveTopic.xhr;
            xhr && xhr.abort && xhr.abort();

            this.setState({saving: true});
            this.saveTopic.xhr = $.ajax({
                url: this.props.put,
                method: "PUT",
                contentType: "application/json",
                dataType: "json",
                data: JSON.stringify({
                    title: this.state.title,
                    name: this.state.name,
                }),
                success: function(result){
                    result.saving = false;
                    this.setState(result);
                }.bind(this),
                error: function() {
                    this.setState({
                        saving: false
                    });
                }.bind(this)
            });
        },

        deleteTopic: function(e){
            e.preventDefault();
            return $.ajax({
                url: this.props.remove,
                method: "DELETE",
                success: function() {
                    app.router.setRoute("/edit/" + this.props.projectId);
                }.bind(this)
            });
        },

        createPage: function(e){
            e.preventDefault();
            return $.ajax({
                url: this.props.post,
                method: "POST",
                data: JSON.stringify({
                    title: this.refs.theNewTitle.getDOMNode().value
                }),
                success: function(result) {
                    this.state.pages = this.state.pages || [];
                    this.state.pages.push(result);
                    this.setState({
                        pages: this.state.pages
                    });
                    this.refs.theNewTitle.getDOMNode().value = "";
                }.bind(this)
            });
        },

        render: function() {
            // Get Page component
            var Page = app.Page;
            var projectName = this.props.projectName;
            var projectId = this.props.projectId;
            var topicId = this.props.topicId;
            var cViews = (this.state.pages || []).map(function(page, key){
                var source = ["/edit", projectId, page.topic, page.id].join("/");
                return <div className="page-edit" data-page={page.id} data-name={page.name}>
                        <div className="page-text" data-page={page.id} data-name={page.name}>
                            <a href={source} data-noreload="true">
                                {page.title}
                            </a>
                        </div>
                    </div>
            });

            var pages = null;
            if (this.state.pages && this.state.pages.length) {
                pages =
                    <div className="list-group">{
                        this.state.pages.map(function(page, key) {
                            var url = ["/edit", projectId, topicId, page.id].join("/");
                            return (<a href={url} data-noreload="true" className="list-group-item" key={key}>
                                <h4 className="list-group-item-heading text-capitalize">{page.title}</h4>
                            </a>);
                        }.bind(this))}
                    </div>
            } else {
                pages =
                    <div className="no-list text-center">
                        No pages
                    </div>
            }

            var btnDisabled = (!!this.state.title && !this.state.saving) ? "" : "disabled";
            var editForm =
                <div>
                    <form role="form" onSubmit={this.saveTopic}>
                        <div className="form-group">
                            <label className="text-muted">Title</label>
                            <input ref="theTitle"
                                type="text"
                                className="form-control"
                                value={this.state.title}
                                onChange={this.titleChange}/>
                        </div>
                        <div className="form-group clearfix">
                            <div className="pull-left">
                                <button className="btn btn-danger pull-right" onClick={this.deleteTopic}>
                                    Delete topic
                                </button>
                            </div>
                            <div className="pull-right">
                                <button className="btn btn-info" disabled={btnDisabled} onClick={this.saveTopic}>
                                    Save topic
                                </button>
                            </div>
                        </div>
                    </form>
                </div>

            var newForm =
                <div className="clearfix">
                    <form role="form" onSubmit={this.createPage}>
                        <div className="form-group">
                            <input ref="theNewTitle"
                                type="text"
                                className="form-control"
                                placeholder="Topic Title"/>
                        </div>
                        <div className="form-group pull-right clearfix">
                            <button className="btn btn-info" onClick={this.createPage}>
                                Create Page
                            </button>
                        </div>
                    </form>
                </div>

            return <div>
                    {editForm}
                    <hr/>
                    {newForm}
                    <hr/>
                    <div className="page-list">{pages}</div>
                </div>
        }
    });
})(this);
