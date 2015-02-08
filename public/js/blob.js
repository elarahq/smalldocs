(function(root){
    var Blob = root.app.views.Blob = React.createClass({
        displayName: "Blob",
        render: function() {
            return (
                <div className="blob-view">
                    <div className="blob-view-content"></div>
                </div>
            );
        },
        componentDidMount: function() {
            var $el = $(this.getDOMNode());
            var node = this.props.node;
            if (node) {
                $el
                    .find(".blob-view-content")
                    .html(node._blob)
                    .find("pre")
                    .addClass('prettyprint');
                prettyPrint();
            }
        }
    });
})(this);
