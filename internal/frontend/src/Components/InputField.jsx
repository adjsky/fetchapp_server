const React = require("react")
const PropTypes = require("prop-types")

function InputField(props) {
  const { value, onChange, onEnterPress } = props
  const [error, setError] = React.useState(false)
  return (
    <input
      className={`user_input ${error ? "input_error" : ""}`}
      type="text"
      value={value}
      onChange={onChange}
      onKeyDown={(event) => {
        if (event.key === "Enter") {
          if (value) {
            onEnterPress()
            setError(false)
          } else {
            setError(true)
          }
        }
      }}
    />
  )
}

InputField.propTypes = {
  value: PropTypes.string.isRequired,
  onChange: PropTypes.func.isRequired,
  onEnterPress: PropTypes.func.isRequired,
}

module.exports = InputField
