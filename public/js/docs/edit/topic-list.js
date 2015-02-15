(function(root){
    // Get app
    var app = root.app;

    // Topic view
    var TopicList = app.TopicList = React.createClass({
        displayName: "TopicList",

        getInitialState: function(){
            return {
                topics: null,
                adding: false,
                saving: false,

                // new topic
                title: "",
                name: "",
            };
        },

        // Fetch topics
        fetch: function() {
            return $.ajax({
                url: this.props.source,
                method: "GET"
            });
        },

        // Check duplicate title
        titleCheck: function() {
            var timeoutId = this.titleCheck.timeoutId;
            timeoutId && clearTimeout(timeoutId);
            var xhr = this.titleCheck.xhr;
            xhr && xhr.abort && xhr.abort();
s
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

        // Save topic
        saveTopic: function(){
            var xhr = this.saveTopic.xhr;
            xhr && xhr.abort && xhr.abort();

            this.setState({saving: true});
            this.saveTopic.xhr = $.ajax({
                url: this.props.post,
                method: "POST",
                contentType: "application/json",
                dataType: "json",
                data: JSON.stringify({
                    title: this.state.title,
                    name: this.state.name,
                }),
                success: function(result){
                    this.state.topics = this.state.topics || [];
                    this.state.topics.push(result);
                    this.setState({
                        adding: false,
                        saving: false,
                        name: "",
                        title: ""
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
                if (this.isMounted()) {
                    this.setState({
                        topics: result
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
            var value = e.target.value || "";
            this.setState({
                title: value.replace(new RegExp("[^a-zA-Z0-9 \.-]", "gi"), "")
            }, function(){
                this.titleCheck();
            });
        },

        render: function() {
            var topics = null;
            var projectName = this.props.projectName;

            if (this.state.topics && this.state.topics.length) {
                topics =
                    <div className="list-group">{
                        this.state.topics.map(function(topic, key) {
                            var url = ["/edit", projectName, topic.name].join("/");
                            return (<a href={url} className="list-group-item" key={key}>
                                <h4 className="list-group-item-heading text-capitalize">{topic.title}</h4>
                                <p className="list-group-item-text">{topic.name}</p>
                            </a>);
                        })}
                    </div>
            } else {
                topics =
                    <div className="no-list text-center">
                        No topics
                    </div>
            }

            var newButton =
                <div className="clearfix padding-bottom-10">
                    <button className="pull-right btn btn-info" onClick={this.startAdding}>+ Create New</button>
                </div>

            var helpCN = "help-block " + (!!this.state.name ? "" : "hide");
            var btnDisabled = (!!this.state.title && !!this.state.name && !this.state.saving) ? "" : "disabled";
            var newForm =
                <div className="clearfix">
                    <form role="form">
                        <div className="form-group">
                            <label className="text-muted">Title</label>
                            <input ref="theTitle" type="text" className="form-control" value={this.state.title} onChange={this.titleChange}/>
                            <p className={helpCN}>This project will be created as <b className="text-info">{this.state.name}</b></p>
                        </div>
                        <div className="form-group pull-right clearfix">
                            <button className="btn btn-default" type="button" onClick={this.cancelAdding}>Cancel</button>&nbsp;
                            <button className="btn btn-info" disabled={btnDisabled} onClick={this.saveTopic}>Create project</button>
                        </div>
                    </form>
                </div>

            return (
                <div>
                    {!this.state.adding?newButton:null}
                    {this.state.adding?newForm:null}
                    {topics}
                </div>
            );
        }
    });
})(this);
