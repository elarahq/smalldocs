(function(root){
    var TreeView = root.app.views.TreeView = React.createClass({
        displayName: "TreeView",

        propTypes: {
            collapsed: React.PropTypes.bool,
            defaultCollapsed: React.PropTypes.bool,
            nodeLabel: React.PropTypes.string,
            isTree: React.PropTypes.bool
        },

        getInitialState: function(){
            return {
                collapsed: this.props.defaultCollapsed
            }
        },

        handleClick: function(e) {
            if (this.props.node.isTree) {
                this.setState({
                    collapsed: !this.state.collapsed
                });
            }
            var props = this.props;
            props.onClick && props.onClick(e, props.node, this.state.collapsed);
        },

        render: function() {
            var props = this.props;
            var node = props.node;

            var collapsed = props.collapsed != null ?
                props.collapsed :
                this.state.collapsed;

            var icon =
                <div
                {...props}
                className='tree-view-icon'
                onClick={this.handleClick}>
                </div>;

            var treeClass = 'tree-view ' + (node.isTree ? 'tree-view-tree' : 'tree-view-blob') + (collapsed ? ' collapsed' : '');
            var children = node.children && node.children.map(function(n, i) {
                return (
                    <TreeView
                        onClick={props.onClick}
                        node={n}
                        isTree={n.isTree}
                        nodeLabel={n.name}
                        key={n.name}
                        defaultCollapsed={true}>
                    </TreeView>
                );
            }.bind(this));
            return (
                <div className={treeClass}>
                    <div className="tree-view-item">
                        {icon}
                        <div onClick={this.handleClick} className="tree-view-label">{props.nodeLabel}</div>
                    </div>
                    <div className="tree-view-children">
                        {children}
                    </div>
                </div>
            );
        }
    });
})(this);
