export function ChangeElement(idParent, newElement) {
	let parent = document.getElementById(idParent)
    const realDOM = CreateRealElement(newElement)
	if (parent) {
		parent.appendChild(realDOM)
	} else {
		document.body.appendChild(realDOM)
	}
}

function CreateRealElement(virtualNode) {
    if (typeof virtualNode === "string") {
        return document.createTextNode(virtualNode)
    }

    const element = document.createElement(virtualNode.tag)

    if (virtualNode.props) {
        Object.keys(virtualNode.props).forEach(propName => {
        	element[propName] = virtualNode.props[propName]
        })
    }
    virtualNode.children.forEach(child => {
		if (child) {
			element.appendChild(CreateRealElement(child))
		}
    })

    return element
}
