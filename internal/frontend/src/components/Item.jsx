const React = require("react")
const PropTypes = require("prop-types")

function Item(props) {
  const { item, onChange, onDelete } = props
  return (
    <div className="item">
      <div className="item_row">
        <span className="item_row_name">{item.value}</span>
        <div className="item_row_controls">
          <button
            type="button"
            className="item_row_control"
            onClick={() => onChange(item.id)}
          >
            Change
          </button>
          <button
            type="button"
            className="item_row_control"
            onClick={() => onDelete(item.id)}
          >
            Delete
          </button>
        </div>
      </div>
    </div>
  )
}

Item.propTypes = {
  item: PropTypes.shape({
    id: PropTypes.string,
    value: PropTypes.string,
  }).isRequired,
  onChange: PropTypes.func.isRequired,
  onDelete: PropTypes.func.isRequired,
}

module.exports = Item
