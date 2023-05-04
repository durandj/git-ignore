package internal_test

type listCall struct {
}

type listReturnValue struct {
	options []string
	err     error
}

type generateCall struct {
	options []string
}

type generateReturnValue struct {
	content string
	err     error
}

type updateCall struct {
}

type updateReturnValue struct {
	err error
}

type fakeAdapter struct {
	listCalls            []listCall
	listReturnValues     []listReturnValue
	generateCalls        []generateCall
	generateReturnValues []generateReturnValue
	updateCalls          []updateCall
	updateReturnValues   []updateReturnValue
}

func (adapter *fakeAdapter) addListReturn(options []string, err error) {
	adapter.listReturnValues = append(adapter.listReturnValues, listReturnValue{
		options: options,
		err:     err,
	})
}

func (adapter *fakeAdapter) List() ([]string, error) {
	adapter.listCalls = append(adapter.listCalls, listCall{})

	returnValue := adapter.listReturnValues[0]
	adapter.listReturnValues = adapter.listReturnValues[1:]

	return returnValue.options, returnValue.err
}

func (adapter *fakeAdapter) addGenerateReturn(content string, err error) {
	adapter.generateReturnValues = append(adapter.generateReturnValues, generateReturnValue{
		content: content,
		err:     err,
	})
}

func (adapter *fakeAdapter) Generate(options []string) (string, error) {
	adapter.generateCalls = append(adapter.generateCalls, generateCall{
		options: options,
	})

	returnValue := adapter.generateReturnValues[0]
	adapter.generateReturnValues = adapter.generateReturnValues[1:]

	return returnValue.content, returnValue.err
}

func (adapter *fakeAdapter) getUpdateCalls() []updateCall {
	return adapter.updateCalls
}

func (adapter *fakeAdapter) addUpdateReturn(err error) {
	adapter.updateReturnValues = append(adapter.updateReturnValues, updateReturnValue{
		err: err,
	})
}

func (adapter *fakeAdapter) Update() error {
	adapter.updateCalls = append(adapter.updateCalls, updateCall{})

	returnValue := adapter.updateReturnValues[0]
	adapter.updateReturnValues = adapter.updateReturnValues[1:]

	return returnValue.err
}

func newFakeAdapter() fakeAdapter {
	return fakeAdapter{
		listCalls:            []listCall{},
		listReturnValues:     []listReturnValue{},
		generateCalls:        []generateCall{},
		generateReturnValues: []generateReturnValue{},
		updateCalls:          []updateCall{},
		updateReturnValues:   []updateReturnValue{},
	}
}
