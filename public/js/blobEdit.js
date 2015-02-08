(function(root){
    var BlobEdit = root.app.views.BlobEdit = React.createClass({
        displayName: "BlobEdit",

        getInitialState: function(){
            return {
                markdown: "",
                saving: false,
            }
        },

        onNameChange: function(e){
            this.props.node.name = e.target.value;
        },

        onContentChange: function(e) {},

        onSave: function(e) {
            if (e) {
                e.preventDefault();
                e.stopPropagation();
            }
            this.props.node._markdown = this.editor.getDoc().getValue();
            if (this.props.onSave) {
                this.setState({saving: true});
                this.props.onSave(e).always(function(){
                    if (this.isMounted()) {
                        this.setState({saving: false});
                    }
                }.bind(this));
            }
        },

        componentWillMount: function(){
            $.ajax({
                url: "/blob/" + this.props.node.path,
                type: "get",
                data: {type: 'markdown'},
                success: function(result) {
                    this.setState({markdown: result});
                }.bind(this)
            });
        },

        render: function() {
            var node = this.props.node;
            var adding = this.props.adding;
            return (
                <div className="blob-view-edit">
                    <form>
                        <div className="row">
                            <div className="twelve columns">
                                <input
                                    className="u-full-width"
                                    type="text"
                                    defaultValue={node.name}
                                    placeholder="Name"
                                    onChange={this.onNameChange}/>
                            </div>
                        </div>
                        <div className="row">
                            <div className="twelve columns">
                                <textarea
                                    ref="markdownEditor"
                                    id="editNodeTextarea"
                                    value={this.state.markdown}
                                    className="u-full-width"
                                    placeholder="Write here"
                                    onChange={this.onContentChange}></textarea>
                            </div>
                        </div>
                        <div className="row">
                            <button value="Cancel" onClick={this.props.onCancel}>
                                Cancel
                            </button>
                            &nbsp;
                            <button className="button-primary" value="Save" onClick={this.onSave}>
                                {this.state.saving ? 'Saving...' : 'Save' }
                            </button>
                            {
                                !adding ?
                                (<button className="u-pull-right" value="Delete Topic" onClick={this.props.onDelete}>
                                    Delete Topic
                                </button>): null
                            }
                        </div>
                    </form>
                </div>
            );
        },

        componentDidUpdate: function(){
            if (!this.editor) {
                this.editor = CodeMirror.fromTextArea(this.refs.markdownEditor.getDOMNode(), {
                    mode: 'markdown',
                    lineNumbers: true,
                    lineWrapping: true,
                    matchBrackets: true
                });
            }
        }
    });
})(this);
