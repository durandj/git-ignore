package internal_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/durandj/git-ignore/internal"
)

func TestClientListWithErrorInAdapterShouldFallbackToNextAdapter(t *testing.T) {
	t.Parallel()

	primaryAdapter := newFakeAdapter()
	secondaryAdapter := newFakeAdapter()

	client := internal.Client{
		Adapters: []internal.Adapter{
			&primaryAdapter,
			&secondaryAdapter,
		},
	}

	expectedOptions := []string{"c"}

	primaryAdapter.addListReturn(nil, errors.New("Primary error"))
	secondaryAdapter.addListReturn(expectedOptions, nil)

	options, err := client.List()

	require.NoError(t, err)
	require.Equal(t, expectedOptions, options)
}

func TestClientListWithErrorInAllAdaptersShouldReturnAnError(t *testing.T) {
	t.Parallel()

	primaryAdapter := newFakeAdapter()
	secondaryAdapter := newFakeAdapter()

	client := internal.Client{
		Adapters: []internal.Adapter{
			&primaryAdapter,
			&secondaryAdapter,
		},
	}

	expectedErr := errors.New("Test error")
	primaryAdapter.addListReturn(nil, expectedErr)
	secondaryAdapter.addListReturn(nil, expectedErr)

	_, err := client.List()

	require.Error(t, err)
}

func TestClientListShouldRetrieveAListOfOptions(t *testing.T) {
	t.Parallel()

	primaryAdapter := newFakeAdapter()
	secondaryAdapter := newFakeAdapter()

	client := internal.Client{
		Adapters: []internal.Adapter{
			&primaryAdapter,
			&secondaryAdapter,
		},
	}

	expectedOptions := []string{"c", "c++"}
	primaryAdapter.addListReturn(expectedOptions, nil)

	options, err := client.List()

	require.NoError(t, err)
	require.Equal(t, expectedOptions, options)
}

func TestClientGenerateWithNoOptionsShouldReturnAnError(t *testing.T) {
	t.Parallel()

	primaryAdapter := newFakeAdapter()
	secondaryAdapter := newFakeAdapter()

	client := internal.Client{
		Adapters: []internal.Adapter{
			&primaryAdapter,
			&secondaryAdapter,
		},
	}

	_, err := client.Generate(nil)

	require.Error(t, err)
}

func TestClientGenerateWithAnInvalidOptionShouldReturnAnError(t *testing.T) {
	t.Parallel()

	primaryAdapter := newFakeAdapter()
	secondaryAdapter := newFakeAdapter()

	client := internal.Client{
		Adapters: []internal.Adapter{
			&primaryAdapter,
			&secondaryAdapter,
		},
	}

	primaryAdapter.addListReturn([]string{"c"}, nil)

	_, err := client.Generate([]string{"doesnotexist"})

	require.Error(t, err)
}

func TestClientGenerateWithASingleOptionShouldGenerateAGitignoreFile(t *testing.T) {
	t.Parallel()

	primaryAdapter := newFakeAdapter()
	secondaryAdapter := newFakeAdapter()

	client := internal.Client{
		Adapters: []internal.Adapter{
			&primaryAdapter,
			&secondaryAdapter,
		},
	}

	primaryAdapter.addListReturn([]string{"c", "c++"}, nil)
	primaryAdapter.addGenerateReturn("### C ###", nil)

	file, err := client.Generate([]string{"c"})

	require.NoError(t, err)
	require.Contains(t, file, "### C ###")
}

func TestClientGenerateWithMultipleOptionsShouldGenerateAGitignoreFile(t *testing.T) {
	t.Parallel()

	primaryAdapter := newFakeAdapter()
	secondaryAdapter := newFakeAdapter()

	client := internal.Client{
		Adapters: []internal.Adapter{
			&primaryAdapter,
			&secondaryAdapter,
		},
	}

	primaryAdapter.addListReturn([]string{"c", "c++"}, nil)
	primaryAdapter.addGenerateReturn("### C ###\n\n### C++ ###", nil)

	file, err := client.Generate([]string{"c", "c++"})

	require.NoError(t, err)
	require.Contains(t, file, "### C ###")
	require.Contains(t, file, "### C++ ###")
}

func TestClientGenerateWithAnErrorInAnAdapterShouldFallbackToNextAdapter(t *testing.T) {
	t.Parallel()

	primaryAdapter := newFakeAdapter()
	secondaryAdapter := newFakeAdapter()

	client := internal.Client{
		Adapters: []internal.Adapter{
			&primaryAdapter,
			&secondaryAdapter,
		},
	}

	primaryAdapter.addListReturn(nil, errors.New("Test error"))
	secondaryAdapter.addListReturn([]string{"c"}, nil)
	secondaryAdapter.addGenerateReturn("### C ###", nil)

	file, err := client.Generate([]string{"c"})

	require.NoError(t, err)
	require.Contains(t, file, "### C ###")
}

func TestClientGenerateWithAnErrorInAllAdaptersShouldReturnAnError(t *testing.T) {
	t.Parallel()

	primaryAdapter := newFakeAdapter()
	secondaryAdapter := newFakeAdapter()

	client := internal.Client{
		Adapters: []internal.Adapter{
			&primaryAdapter,
			&secondaryAdapter,
		},
	}

	primaryAdapter.addListReturn(nil, errors.New("Test error"))
	secondaryAdapter.addListReturn(nil, errors.New("Test error"))

	_, err := client.Generate([]string{"c"})

	require.Error(t, err)
}

func TestClientUpdateShouldUpdateAllAdapters(t *testing.T) {
	t.Parallel()

	primaryAdapter := newFakeAdapter()
	secondaryAdapter := newFakeAdapter()

	client := internal.Client{
		Adapters: []internal.Adapter{
			&primaryAdapter,
			&secondaryAdapter,
		},
	}

	primaryAdapter.addUpdateReturn(nil)
	secondaryAdapter.addUpdateReturn(nil)

	err := client.Update()

	require.NoError(t, err)

	primaryUpdateCalls := primaryAdapter.getUpdateCalls()
	secondaryUpdateCalls := secondaryAdapter.getUpdateCalls()

	require.Len(t, primaryUpdateCalls, 1)
	require.Len(t, secondaryUpdateCalls, 1)
}

func TestClientUpdateWithAnErrorInOneOrMoreAdaptersShouldReturnAnError(t *testing.T) {
	t.Parallel()

	primaryAdapter := newFakeAdapter()
	secondaryAdapter := newFakeAdapter()

	client := internal.Client{
		Adapters: []internal.Adapter{
			&primaryAdapter,
			&secondaryAdapter,
		},
	}

	primaryAdapter.addUpdateReturn(errors.New("Test error"))
	secondaryAdapter.addUpdateReturn(nil)

	err := client.Update()

	require.Error(t, err)
}
