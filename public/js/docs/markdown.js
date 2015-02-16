(function(root){
    // Get app
    var app = root.app;

    // Markdown
    var Markdown = app.Markdown = React.createClass({
        displayName: "Markdown",

        getInitialState: function(){
            return {
                content: ""
            };
        },

        fetch: function(url){
            return $.ajax({
                url: url,
                method: "GET",
                success: function(result){
                    this.setState(result);
                }.bind(this)
            });
        },

        componentWillMount: function(){
            this.dispatchToken = app.dispatcher.register(function(payload){
                switch(payload.actionType) {
                    case "change:page":
                        this.fetch("/pages/" + payload.pageId);
                        break;
                }
            }.bind(this));
        },

        componentWillUnmount: function(){
            app.dispatcher.unregister(this.dispatchToken);
        },

        componentDidUpdate: function(){
            var $el = $(this.getDOMNode());
            $el
                .find("pre")
                .addClass('prettyprint');
            prettyPrint();
        },

        render: function(){
            return <div className="markdown">
                <div dangerouslySetInnerHTML={{__html: marked(this.state.content || "")}}></div>
            </div>;
        }
    });

})(this);
