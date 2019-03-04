package gitignore_test

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

type sourceCall struct {
}

type sourceReturnValue struct {
	mappings map[string]string
	err      error
}

type cacheCall struct {
	mappings map[string]string
}

type cacheReturnValue struct {
	err error
}

type fakeAdapter struct {
	listCalls            []listCall
	listReturnValues     []listReturnValue
	generateCalls        []generateCall
	generateReturnValues []generateReturnValue
	sourceCalls          []sourceCall
	sourceReturnValues   []sourceReturnValue
	cacheCalls           []cacheCall
	cacheReturnValues    []cacheReturnValue
}

func (adapter *fakeAdapter) getListCalls() []listCall {
	return adapter.listCalls
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

func (adapter *fakeAdapter) getGenerateCalls() []generateCall {
	return adapter.generateCalls
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

func (adapter *fakeAdapter) getSourceCalls() []sourceCall {
	return adapter.sourceCalls
}

func (adapter *fakeAdapter) addSourceReturn(mappings map[string]string, err error) {
	adapter.sourceReturnValues = append(adapter.sourceReturnValues, sourceReturnValue{
		mappings: mappings,
		err:      err,
	})
}

func (adapter *fakeAdapter) Source() (map[string]string, error) {
	adapter.sourceCalls = append(adapter.sourceCalls, sourceCall{})

	returnValue := adapter.sourceReturnValues[0]
	adapter.sourceReturnValues = adapter.sourceReturnValues[1:]

	return returnValue.mappings, returnValue.err
}

func (adapter *fakeAdapter) getCacheCalls() []cacheCall {
	return adapter.cacheCalls
}

func (adapter *fakeAdapter) addCacheReturn(err error) {
	adapter.cacheReturnValues = append(adapter.cacheReturnValues, cacheReturnValue{
		err: err,
	})
}

func (adapter *fakeAdapter) Cache(ignoreMapping map[string]string) error {
	adapter.cacheCalls = append(adapter.cacheCalls, cacheCall{
		mappings: ignoreMapping,
	})

	returnValue := adapter.cacheReturnValues[0]
	adapter.cacheReturnValues = adapter.cacheReturnValues[1:]

	return returnValue.err
}

func newFakeAdapter() fakeAdapter {
	return fakeAdapter{
		listCalls:            []listCall{},
		listReturnValues:     []listReturnValue{},
		generateCalls:        []generateCall{},
		generateReturnValues: []generateReturnValue{},
		sourceCalls:          []sourceCall{},
		sourceReturnValues:   []sourceReturnValue{},
		cacheCalls:           []cacheCall{},
		cacheReturnValues:    []cacheReturnValue{},
	}
}
