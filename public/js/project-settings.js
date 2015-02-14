(function(root){
    // Get app
    var app = root.app;

    // Project view
    var ProjectSettings = app.ProjectSettings = React.createClass({
        displayName: "ProjectSettings",

        getInitialState: function(){
            return {
                saving: false,

                // project
                title: "",
                description: "",
                name: "",
            };
        },

        // Fetch projects
        fetch: function() {
            return $.ajax({
                url: this.props.source.replace("ID", $('body').data("id")),
                method: "GET"
            });
        },

        // Check duplicate title
        titleCheck: function() {
            var timeoutId = this.titleCheck.timeoutId;
            timeoutId && clearTimeout(timeoutId);
            var xhr = this.titleCheck.xhr;
            xhr && xhr.abort && xhr.abort();

            this.titleCheck.timeoutId = setTimeout(function(){
                this.titleCheck.xhr = $.ajax({
                    url: this.props.check,
                    method: "POST",
                    contentType: "application/json",
                    dataType: "json",
                    data: JSON.stringify({
                        title: this.state.title,
                        id: this.state.id
                    }),
                    success: function(result){
                        this.setState(result);
                    }.bind(this),
                    error: function() {
                        this.setState({
                            name: ""
                        });
                    }.bind(this)
                });
            }.bind(this), 500);
        },

        // Save project
        saveProject: function(){
            var xhr = this.saveProject.xhr;
            xhr && xhr.abort && xhr.abort();

            this.setState({saving: true});
            this.saveProject.xhr = $.ajax({
                url: this.props.put.replace("ID", $('body').data('id')),
                method: "PUT",
                contentType: "application/json",
                dataType: "json",
                data: JSON.stringify({
                    title: this.state.title,
                    description: this.state.description.trim(),
                    name: this.state.name,
                }),
                success: function(result){
                    this.state.projects = this.state.projects || [];
                    this.state.projects.unshift(result);
                    this.setState({
                        saving: false,
                    });
                }.bind(this),
                error: function() {
                    this.setState({
                        saving: false
                    });
                }.bind(this)
            });
        },

        componentDidMount: function() {
            this.fetch().success(function(result) {
                this.currentName = result.name;
                if (this.isMounted()) {
                    this.setState(result);
                }
            }.bind(this));
        },

        titleChange: function(e){
            var value = e.target.value || "";
            this.setState({
                title: value.replace(new RegExp("[^a-zA-Z0-9 \.-]", "gi"), "")
            }, function(){
                this.titleCheck();
            });
        },

        descChange: function(e){
            this.setState({
                description: (e.target.value || "")
            });
        },

        render: function() {
            var isCurrent = (this.currentName && (this.currentName == this.state.name));
            var helpCN = "help-block " + ((!!this.state.name && !isCurrent) ? "" : "hide");
            var btnDisabled = (!!this.state.title && !!this.state.name && !this.state.saving) ? "" : "disabled";
            var form =
                <form role="form">
                    <div className="form-group">
                        <label className="text-muted">Title</label>
                        <input ref="theTitle" type="text" className="form-control" value={this.state.title} onChange={this.titleChange}/>
                        <p className={helpCN}>This project will be saved as <b className="text-info">{this.state.name}</b></p>
                    </div>
                    <div className="form-group">
                        <label className="text-muted">Description (optional)</label>
                        <input type="text" className="form-control" value={this.state.description} onChange={this.descChange}/>
                    </div>
                    <div className="form-group pull-right clearfix">
                        <a href="/" className="btn btn-default">Cancel</a>&nbsp;
                        <button className="btn btn-info" disabled={btnDisabled} onClick={this.saveProject}>Save project</button>
                    </div>
                </form>

            return (<div>{form}</div>);
        }
    });

    // Load project settings
    React.render(
        <ProjectSettings source="/projects/ID" put="/projects/ID" check="/projects_check"/>,
        document.getElementById('settings')
    );
})(this);
