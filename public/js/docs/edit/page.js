(function(root){
    // Get app
    var app = root.app;

    // Page
    var Page = app.Page = React.createClass({
        displayName: "Page",

        getInitialState: function() {
            return {
            }
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

        componentDidMount: function(){
            this.fetch();
        },

        componentDidUpdate: function(){
            if (!this.editor) {
                this.editor = CodeMirror.fromTextArea(this.refs.markdownEditor.getDOMNode(), {
                    mode: 'markdown',
                    lineWrapping: true,
                    matchBrackets: true,
                    lineNumbers: true,
                    theme: 'solarized light'
                });
                this.editor.on("change", this.onContentChange);
            }
        },

        onContentChange: function(){
            this.setState({
                content: this.editor.getDoc().getValue()
            });
        },

        titleChange: function(e){
            var value = e.target.value || "";
            this.setState({
                title: value.replace(new RegExp("[^a-zA-Z0-9 \.-]", "gi"), "")
            })
        },

        // Save page
        savePage: function(){
            var xhr = this.savePage.xhr;
            xhr && xhr.abort && xhr.abort();

            this.setState({saving: true});
            this.savePage.xhr = $.ajax({
                url: this.props.put,
                method: "PUT",
                contentType: "application/json",
                dataType: "json",
                data: JSON.stringify({
                    title: this.state.title,
                    content: this.state.content,
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

        deletePage: function(e){
            e.preventDefault();
            return $.ajax({
                url: this.props.remove,
                method: "DELETE",
                success: function() {
                    app.router.setRoute("/edit/" + this.props.projectId + "/" + this.props.topicId);
                }.bind(this)
            });
        },

        render: function() {
            // Get Page component
            var Page = app.Page;
            var projectName = this.props.projectName;
            var projectId = this.props.projectId;
            var topicId = this.props.topicId;

            var btnDisabled = (!!this.state.title && !this.state.saving) ? "" : "disabled";
            var editForm =
                <div>
                    <form role="form" onSubmit={this.savePage}>
                        <div className="form-group">
                            <label className="text-muted">Title</label>
                            <input ref="theTitle"
                                type="text"
                                className="form-control"
                                value={this.state.title}
                                onChange={this.titleChange}/>
                        </div>
                        <div className="form-group editor">
                            <textarea
                                ref="markdownEditor"
                                id="editNodeTextarea"
                                value={this.state.content}
                                className="u-full-width"
                                placeholder="Write here"
                                onChange={function(){}}></textarea>
                        </div>
                        <div className="form-group clearfix">
                            <div className="pull-left">
                                <button className="btn btn-danger pull-right" onClick={this.deletePage}>
                                    Delete page
                                </button>
                            </div>
                            <div className="pull-right">
                                <button className="btn btn-info" disabled={btnDisabled} onClick={this.savePage}>
                                    Save page
                                </button>
                            </div>
                        </div>
                    </form>
                </div>

            return <div>
                    <div className="clearfix">
                        <a data-noreload="true"
                            href={"/edit/"+this.props.projectId+"/"+this.props.topicId}
                            className="pull-right">
                            &#x276e; &nbsp; All pages
                        </a>
                    </div>
                    {editForm}
                </div>
        }
    });
})(this);
