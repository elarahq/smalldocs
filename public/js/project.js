(function(root){
    // Get app
    var app = root.app;

    // Project view
    var Projects = app.Projects = React.createClass({
        displayName: "Projects",

        getInitialState: function(){
            return {
                projects: null
            };
        },

        // Fetch projects
        fetch: function() {
            return $.ajax({
                url: this.props.source,
                method: "GET"
            });
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

        render: function() {
            var projects = (this.state.projects || []).map(function(project, key) {
                return (<a href="/" className="list-group-item text-capitalize" key={key}>
                      <h4 className="list-group-item-heading">{project.title}</h4>
                      <p className="list-group-item-text">{project.description || "hello"}</p>
                    </a>);
            });
            return (<div className="list-group">{projects}</div>);
        }
    });

    // Load and show all projects
    React.render(
        <Projects source="/projects/all" post="/projects"/>,
        document.getElementById('projects')
    );
})(this);
