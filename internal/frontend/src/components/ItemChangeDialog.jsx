const React = require("react")
const ReactDOM = require("react-dom")
const PropTypes = require("prop-types")

class ItemChangeDialog extends React.Component {
  constructor(props) {
    super(props)
    this.state = {
      inputValue: "",
      error: false,
    }

    this.handleKeyboardListener = this.handleKeyboardListener.bind(this)
    this.onInputKeyPress = this.onInputKeyPress.bind(this)
    this.onConfirmClick = this.onConfirmClick.bind(this)
    this.close = this.close.bind(this)
  }

  componentDidMount() {
    window.addEventListener("keydown", this.handleKeyboardListener)
  }

  componentWillUnmount() {
    window.removeEventListener("keydown", this.handleKeyboardListener)
  }

  handleKeyboardListener(event) {
    const { show } = this.props
    if (show && event.key === "Escape") {
      this.close()
    }
  }

  onInputKeyPress(e) {
    const { inputValue } = this.state
    if (e.key === "Enter") {
      if (inputValue) {
        this.confirm()
      } else {
        this.setState({ error: true })
      }
    }
  }

  onConfirmClick() {
    const { inputValue } = this.state
    if (inputValue) {
      this.confirm()
    } else {
      this.setState({ error: true })
    }
  }

  confirm() {
    const { onConfirm } = this.props
    const { inputValue } = this.state
    onConfirm(inputValue)
    this.setState({ inputValue: "", error: false })
  }

  close() {
    const { onClose } = this.props
    onClose()
    this.setState({ inputValue: "", error: false })
  }

  render() {
    const { inputValue, error } = this.state
    const { show } = this.props

    if (!show) {
      return null
    }

    return ReactDOM.createPortal(
      <div className="dialog_bg">
        <div className="dialog">
          <h3 className="dialog_title">Enter a new todo</h3>
          <input
            className={`dialog_input ${error ? "input_error" : ""}`}
            type="text"
            value={inputValue}
            onChange={(e) => this.setState({ inputValue: e.target.value })}
            onKeyPress={this.onInputKeyPress}
          />
          <div className="dialog_controls">
            <button
              type="button"
              className="dialog_confirm_btn"
              onClick={this.onConfirmClick}
            >
              Confirm
            </button>
            <button
              type="button"
              className="dialog_close_btn"
              onClick={this.close}
            >
              Close
            </button>
          </div>
        </div>
      </div>,
      document.getElementById("root"),
    )
  }
}

ItemChangeDialog.propTypes = {
  show: PropTypes.bool.isRequired,
  onConfirm: PropTypes.func.isRequired,
  onClose: PropTypes.func.isRequired,
}

module.exports = ItemChangeDialog
