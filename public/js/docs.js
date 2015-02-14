(function(root){
    // Get app
    var app = root.app;

    // Docs view
    var Docs = app.Docs = React.createClass({
        displayName: "Docs",

        getInitialState: function(){
            return {
            };
        },

        componentDidMount: function() {
        },

        render: function() {
            var topicUrl = this.props.source.replace("ID", $("body").data("id"));
            // Get TopicList component
            var TopicList = app.TopicList;

            return <div>
                    <div className="col-sm-4 col-md-3">
                        <TopicList source={topicUrl}/>
                    </div>
                    <div className="col-sm-8 col-md-9">
                    </div>
                </div>;
        }
    });

    // Load docs
    React.render(
        <Docs source="/docs/ID"/>,
        document.getElementById('docs')
    );
})(this);
