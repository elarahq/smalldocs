(function(root){
    // Get app
    var app = root.app;

    // Topic view
    var TopicList = app.TopicList = React.createClass({
        displayName: "TopicList",

        getInitialState: function(){
            return {
                topics: null,
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
            }).success(function(result) {
                if (this.isMounted()) {
                    this.setState({
                        topics: result
                    });
                }
            }.bind(this));
        },

        // Save topic
        createTopic: function(){
            var xhr = this.createTopic.xhr;
            xhr && xhr.abort && xhr.abort();

            this.setState({saving: true});
            this.createTopic.xhr = $.ajax({
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
            this.fetch();
        },

        titleChange: function(e){
            var value = e.target.value || "";
            this.setState({
                title: value.replace(new RegExp("[^a-zA-Z0-9 \.-]", "gi"), "")
            });
        },

        render: function() {
            var topics = null;
            var projectId = this.props.projectId;

            if (this.state.topics && this.state.topics.length) {
                topics =
                    <div className="list-group">{
                        this.state.topics.map(function(topic, key) {
                            var url = ["/edit", projectId, topic.id].join("/");
                            return (<a href={url} data-noreload="true" className="list-group-item" key={key}>
                                <h4 className="list-group-item-heading text-capitalize">{topic.title}</h4>
                            </a>);
                        }.bind(this))}
                    </div>
            } else {
                topics =
                    <div className="no-list text-center">
                        No topics
                    </div>
            }

            var btnDisabled = (!!this.state.title && !this.state.saving) ? "" : "disabled";
            var newForm =
                <div className="clearfix">
                    <form role="form">
                        <div className="form-group">
                            <input ref="theTitle"
                                type="text"
                                className="form-control"
                                placeholder="Topic Title"
                                value={this.state.title}
                                onChange={this.titleChange}/>
                        </div>
                        <div className="form-group pull-right clearfix">
                            <button className="btn btn-info" disabled={btnDisabled} onClick={this.createTopic}>
                                Create Topic
                            </button>
                        </div>
                    </form>
                </div>

            return (
                <div>
                    {newForm}
                    {topics}
                </div>
            );
        }
    });
})(this);
