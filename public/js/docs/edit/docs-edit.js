(function(root){
    // Get app
    var app = root.app;

    // Docs view
    var Docs = app.Docs = React.createClass({
        displayName: "Docs",

        render: function() {
            var topicUrl = this.props.source.replace("ID", $("body").data("id"));
            var checkUrl = topicUrl+"_check";

            var projectName = $("body").data("name")
            // Get TopicList component
            var TopicList = app.TopicList;

            return <div className="padding-top-10">
                    <div className="col-sm-6 col-sm-offset-3 padding-top-10">
                        <TopicList source={topicUrl} post={topicUrl} check={checkUrl} projectName={projectName}/>
                    </div>
                </div>;
        }
    });

    // Load docs
    React.render(
        <Docs source="/projects/ID/topics"/>,
        document.getElementById('docs-edit')
    );
})(this);
