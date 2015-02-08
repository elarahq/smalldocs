(function(root){
    var Search = root.app.views.Search = React.createClass({
        displayName: "Search",
        render: function() {
            return (
                <div className="search-bar">
                    <input type="text" name="search" placeholder="Search" className="u-full-width"/>
                </div>
            )
        }
    });
})(this);
