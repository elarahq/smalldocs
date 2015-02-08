(function(root){
    var app = root.app;
    var TreeView = app.views.TreeView;
    var Directory = app.views.Directory = React.createClass({
        displayName: "Directory",

        getInitialState: function(){
            return  {
                currentNode: null,
                data: app.nodeStore
            };
        },

        componentWillMount: function(){
            this.dispatchToken = app.dispatcher.register(function(payload){
                switch(payload.actionType) {
                    case "reset:nodes":
                        this.setState({data: payload.nodes});
                        break;
                    case "change:currentNode":
                        this.setState({currentNode: payload.currentNode}, this.fetchData);
                        break;
                }
            }.bind(this));
        },

        componentWillUnmount: function(){
            app.dispatcher.unregister(this.dispatchToken);
        },

        retrieveDirectory: function(){
            var node = this.state.currentNode;
            var rootDir = (node && node.path) || "";
            if (node && node._xhr) return;

            var xhr = $.get("/tree/" + rootDir, function(result) {
                if (node) {
                    node._fetched = true;
                    node._xhr = null;
                    node.children = result;
                    this.setState({data: this.state.data});
                } else {
                    app.dispatcher.dispatch({
                        actionType: "reset:nodes",
                        nodes: result
                    });
                }
            }.bind(this))
            .fail(function() {
                node._fetched = false;
                node._xhr = null;
            });

            if (node) {
                node._xhr = xhr;
            }
        },

        showBlob: function() {
            var node = this.state.currentNode;
            var path = (node && node.path) || "";
            if (node && node._xhr) return;

            var xhr = $.get("/blob/" + path + "?type=html", function(result) {
                if (node) {
                    node._fetched = true;
                    node._blob = result;
                    node._xhr = null;
                    app.dispatcher.dispatch({
                        actionType: "change:currentNode.blob",
                        currentNode: this.state.currentNode
                    });
                }
            }.bind(this))
            .fail(function() {
                node._fetched = false;
                node._blob = null;
                node._xhr = null;
            });

            if (node) {
                node._xhr = xhr;
            }
        },

        componentDidMount: function() {
            this.retrieveDirectory();
        },

        handleClick: function(e, node, collapsed) {
            app.dispatcher.dispatch({
                actionType: "change:currentNode",
                currentNode: node
            });
        },

        fetchData: function() {
            var node = this.state.currentNode;
            if (node && !node._fetched) {
                if (node.isTree) {
                    this.retrieveDirectory();
                } else {
                    this.showBlob();
                }
            }
        },

        render: function() {
            var data = this.state.data;
            var result =  data.map(function(node, key){
                return (
                    <TreeView node={node}
                        onClick={this.handleClick}
                        isTree={node.isTree}
                        key={key}
                        nodeLabel={node.name}
                        defaultCollapsed={true}>
                    </TreeView>
                );
            }.bind(this));
            return <div className="directory">{result}</div>
        }
    });
})(this);
