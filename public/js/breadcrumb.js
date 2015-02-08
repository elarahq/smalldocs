(function(root){
    var app = root.app;
    var Breadcrumb = app.views.Breadcrumb = React.createClass({
        displayName: "Breadcrumb",

        getInitialState: function(){
            return {
            }
        },

        pathTokens: function() {
            var node = this.props.node;
            var result = [];
            if (node && node.path) {
                result = node.path.split("/").map(function(token){
                    token = token.toLowerCase().replace(/_+/g, " ");
                    return token.substr(0, token.lastIndexOf('.')) || token;
                }).filter(function(token){
                    return !!token;
                });
            }
            return result;
        },

        navigateToRoot: function(e) {
            app.dispatcher.dispatch({
                actionType: "change:currentNode",
                currentNode: null
            });
            if (e) {
                e.preventDefault();
                e.stopPropagation();
            }
        },

        onNewAdd: function(e) {
            e && e.preventDefault();
            app.dispatcher.dispatch({
                actionType: "new:node",
                currentNode: this.props.node
            });
        },

        onEdit: function(e) {
            e && e.preventDefault();
            app.dispatcher.dispatch({
                actionType: "edit:node",
                node: this.props.node
            });
        },

        render: function() {
            var node = this.props.node;
            var newFile = null,
                editFile = null;

            if (!node || node.isTree) {
                newFile =
                    <li className="new-add">
                        <a href="" onClick={this.onNewAdd}>+ Add New</a>
                    </li>
            }

            if (node && !node.isTree) {
                editFile =
                    <li className="edit-add">
                        <a href="" onClick={this.onEdit}>Edit Topic</a>
                    </li>
            }

            return (
                <div className="breadcrumb">
                    <ul>
                        <li><a href="" onClick={this.navigateToRoot}>Jump to Root</a></li>
                        {this.pathTokens().map(function(token, i){
                            return <li key={i}>{token}</li>
                        })}
                        {newFile}
                        {editFile}
                    </ul>
                </div>
            )
        }
    });
})(this);
