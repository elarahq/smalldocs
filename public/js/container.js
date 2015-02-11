(function(root){
    var views = root.app.views;
    var Directory = views.Directory,
        Search = views.Search,
        Breadcrumb = views.Breadcrumb,
        Blob = views.Blob,
        BlobEdit = views.BlobEdit;

    var Container = root.app.views.Container = React.createClass({
        displayName: "Container",

        getInitialState: function(){
            return {
                editing: false,
                adding: false,
                currentNode: null
            }
        },

        componentWillMount: function(){
            this.dispatchToken = app.dispatcher.register(function(payload){
                switch(payload.actionType) {
                    case "change:currentNode":
                    case "change:currentNode.blob":
                        this.setState({currentNode: payload.currentNode, editing: false, adding: false});
                        break;
                    case "edit:node":
                        this.setState({editing: true, currentNode: payload.node});
                        break;
                    case "new:node":
                        this.setState({adding: true});
                        break;
                }
            }.bind(this));
        },

        componentWillUnmount: function(){
            app.dispatcher.unregister(this.dispatchToken);
        },

        onEditCancel: function(e){
            this.setState({adding: false, editing: false});
            this.editNode = null;
            e && e.preventDefault();
        },

        onSave: function(e){
            var currentNode = this.state.currentNode;

            var node = this.editNode;
            if (this.state.adding) {
                node.path = ((currentNode && currentNode.path) || "")
                node.path += node.name.trim().toLowerCase().replace(/ +/g, "_") + ".md";
            }

            // Save blob
            return $.ajax({
                url: "/blob/" + (node.path || ""),
                method: "post",
                data: {
                    name: node.name,
                    markdown: node._markdown
                },
                success: function(result) {
                    var cn = this.state.currentNode;
                    var isNew = null, parentNode = null;
                    var ncn = null;

                    if (cn) {
                        if (!cn.isTree) {
                            cn.name = result.name;
                            cn.path = result.path;
                            cn.isTree = result.isTree;
                            cn._fetched = false;
                            cn._blob = false;
                            ncn = cn;
                        } else {
                            cn.children = cn.children || [];
                            ch.children.push(result);
                            ncn = result;
                        }
                    } else {
                        app.nodeStore.push(result)
                        ncn = result;
                    }

                    ncn && app.dispatcher.dispatch({
                        actionType: "change:currentNode",
                        currentNode: ncn
                    });
                }.bind(this)
            });
        },

        onDelete: function(e) {
            e && e.preventDefault();
            var node = this.editNode;
            // delete blob
            return $.ajax({
                url: "/blob/" + (node.path || ""),
                method: "delete",
                success: function() {
                    app.dispatcher.dispatch({
                        actionType: "change:currentNode",
                        currentNode: null
                    });
                }
            });
        },

        render: function() {
            var node = this.state.currentNode;

            var breadcrumb = <Breadcrumb node={this.state.currentNode} />;
                edit = null,
                blob =
                    <div className="blob-placeholder">
                        No Topic Selected
                    </div>;

            if (node && node._blob) {
                blob = <Blob node={this.state.currentNode} />;
            }

            // add/edit mode
            if (this.state.editing) {
                var node = this.state.currentNode || {};
                this.editNode = {
                    name: node.name,
                    path: node.path
                };
                edit = <BlobEdit
                        node={this.editNode}
                        onCancel={this.onEditCancel}
                        onSave={this.onSave}
                        onDelete={this.onDelete}/>

            } else if (this.state.adding) {
                this.editNode = {};
                edit = <BlobEdit
                        node={this.editNode}
                        onCancel={this.onEditCancel}
                        onSave={this.onSave}
                        adding={true}/>
            }

            return (
                <div className="container u-full-width">
                    <div className="row">
                        <div className="three columns">
                            <Search />
                            <Directory />
                        </div>
                        <div className="nine columns">
                            {breadcrumb}
                            {(this.state.editing || this.state.adding) ? edit : blob}
                        </div>
                    </div>
                </div>
            )
        }
    });

    React.render(
        <Container />,
        document.getElementById('container')
    );
})(this);
