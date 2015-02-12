(function(root){
    // Get app
    var app = root.app;

    // Project view
    var Projects = app.Projects = React.createClass({
        displayName: "Projects",

        getInitialState: function(){
            return {
                projects: null,
                adding: false,

                // new project
                title: "",
                description: "",
                name: "",
            };
        },

        // Fetch projects
        fetch: function() {
            return $.ajax({
                url: this.props.source,
                method: "GET"
            });
        },

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
                        title: this.state.title
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

        componentDidMount: function() {
            this.fetch().success(function(result) {
                if (this.isMounted()) {
                    this.setState({
                        projects: result
                    });
                }
            }.bind(this));
        },

        cancelAdding: function(){
            this.setState({
                adding: false
            });
        },

        startAdding: function(){
            this.setState({
                adding: true
            }, function(){
                this.refs.theTitle.getDOMNode().focus();
            });
        },

        titleChange: function(e){
            this.setState({
                title: (e.target.value || "").trim()
            }, function(){
                this.titleCheck();
            });
        },

        descChange: function(e){
            this.setState({
                description: (e.target.value || "").trim()
            });
        },

        render: function() {
            var projects = null;

            if (this.state.projects && this.state.projects.length) {
                projects =
                    <div className="list-group">{
                        this.state.projects.map(function(project, key) {
                            return (<a href="/" className="list-group-item text-capitalize" key={key}>
                                <h4 className="list-group-item-heading">{project.title}</h4>
                                <p className="list-group-item-text">{project.description}</p>
                            </a>);
                        })}
                    </div>
            } else {
                projects =
                    <div className="no-list text-center">
                        No projects
                    </div>
            }

            var newButton =
                <div className="clearfix padding-bottom-10">
                    <button className="pull-right btn btn-info" onClick={this.startAdding}>+ Create New</button>
                </div>

            var helpCN = "help-block " + (!!this.state.name ? "" : "hide");
            var btnDisabled = (!!this.state.title && !!this.state.name) ? "" : "disabled";
            var newForm =
                <div className="clearfix">
                    <form role="form">
                        <div className="form-group">
                            <label className="text-muted">Title</label>
                            <input ref="theTitle" type="text" className="form-control" onChange={this.titleChange}/>
                            <p className={helpCN}>This project will be created as <b className="text-info">{this.state.name}</b></p>
                        </div>
                        <div className="form-group">
                            <label className="text-muted">Description (optional)</label>
                            <input type="text" className="form-control" onChange={this.descChange}/>
                        </div>
                        <div className="form-group pull-right clearfix">
                            <button className="btn btn-default" type="button" onClick={this.cancelAdding}>Cancel</button>&nbsp;
                            <button className="btn btn-info" disabled={btnDisabled}>Create project</button>
                        </div>
                    </form>
                </div>

            return (
                <div>
                    {!this.state.adding?newButton:null}
                    {this.state.adding?newForm:null}
                    {projects}
                </div>
            );
        }
    });

    // Load and show all projects
    React.render(
        <Projects source="/projects/all" post="/projects" check="/projects_check"/>,
        document.getElementById('projects')
    );
})(this);
