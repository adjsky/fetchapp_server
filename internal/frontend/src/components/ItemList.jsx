const React = require("react")
const PropTypes = require("prop-types")
const Item = require("./Item.jsx")
const LanguageContext = require("../contexts/LanguageContext.jsx")

class ItemList extends React.Component {
  constructor(props) {
    super(props)

    this.generateItemList = this.generateItemList.bind(this)
  }

  generateItemList() {
    const { items, onItemChange, onItemDelete } = this.props
    return items.map((item) => (
      <Item key={item.id} item={item} onChange={onItemChange} onDelete={onItemDelete} />
    ))
  }

  render() {
    return (
      <div className="items">
        {this.generateItemList()}
      </div>
    )
  }
}

ItemList.contextType = LanguageContext

ItemList.propTypes = {
  items: PropTypes.arrayOf(PropTypes.object).isRequired,
  onItemChange: PropTypes.func.isRequired,
  onItemDelete: PropTypes.func.isRequired,
}

module.exports = ItemList
