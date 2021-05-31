const { nanoid } = require("nanoid")
const React = require("react")
const ReactDOM = require("react-dom")
const InputField = require("./Components/InputField.jsx")
const ItemList = require("./Components/ItemList.jsx")
const ItemChangeDialog = require("./Components/ItemChangeDialog.jsx")
const LanguageContext = require("./Contexts/LanguageContext.jsx")
require("./style.css")

class App extends React.Component {
  constructor(props) {
    super(props)
    this.state = {
      inputValue: "",
      items: [],
      idOfItemToChange: null,
    }

    this.onInputEnterPress = this.onInputEnterPress.bind(this)
    this.onDialogConfirm = this.onDialogConfirm.bind(this)
    this.onItemChange = this.onItemChange.bind(this)
    this.onItemDelete = this.onItemDelete.bind(this)
  }

  onInputEnterPress() {
    const { inputValue, items } = this.state
    this.setState({
      items: [{ id: nanoid(), value: inputValue }, ...items],
      inputValue: "",
    })
  }

  onDialogConfirm(result) {
    const { idOfItemToChange, items } = this.state
    const updatedItems = [...items]
    const itemIndex = items.findIndex((item) => item.id === idOfItemToChange)
    updatedItems[itemIndex].value = result
    this.setState({ items: updatedItems, idOfItemToChange: null })
  }

  onItemChange(id) {
    this.setState({ idOfItemToChange: id })
  }

  onItemDelete(id) {
    const { items } = this.state
    this.setState({ items: items.filter((item) => item.id !== id) })
  }

  render() {
    const { items, inputValue, idOfItemToChange } = this.state
    return (
      <LanguageContext.Provider value="en">
        <div className="container">
          <h1 className="title">Hello!</h1>
          <InputField
            value={inputValue}
            onChange={(e) => this.setState({ inputValue: e.target.value })}
            onEnterPress={this.onInputEnterPress}
          />
          <ItemList
            items={items}
            onItemChange={this.onItemChange}
            onItemDelete={this.onItemDelete}
          />

          <ItemChangeDialog
            show={idOfItemToChange !== null}
            onClose={() => this.setState({ idOfItemToChange: null })}
            onConfirm={this.onDialogConfirm}
          />
        </div>
      </LanguageContext.Provider>
    )
  }
}

ReactDOM.render(
  <App />,
  document.getElementById("root"),
)
