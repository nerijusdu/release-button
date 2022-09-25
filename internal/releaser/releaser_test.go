package releaser

import (
	"nerijusdu/release-button/internal/argoApi"
	"nerijusdu/release-button/internal/config"
	"nerijusdu/release-button/internal/mocks"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestIsInSyncIsTrue(t *testing.T) {
	ctrl := gomock.NewController(t)
	aApi := mocks.NewMockIArgoApi(ctrl)
	defer ctrl.Finish()
	apps := &argoApi.Applications{}
	c := &config.Config{
		Selectors: map[string]string{"foo": "bar"},
	}
	releaser := &Releaser{
		argoApi: aApi,
		configs: c,
	}

	aApi.EXPECT().GetApps(c.Selectors, true).Return(apps, nil)

	r, _ := releaser.IsInSync()
	assert.Equal(t, true, r)
}

func TestIsInSyncIsFalse(t *testing.T) {
	ctrl := gomock.NewController(t)
	aApi := mocks.NewMockIArgoApi(ctrl)
	defer ctrl.Finish()
	apps := &argoApi.Applications{
		Items: []argoApi.Application{{
			Metadata: argoApi.AppMeta{
				Name: "test",
			},
			Status: argoApi.AppStatus{
				Sync: argoApi.AppStatusSync{Status: "OutOfSync"},
			},
		}},
	}
	c := &config.Config{
		Selectors: map[string]string{"foo": "bar"},
		Allowed:   []string{"test"},
	}
	releaser := &Releaser{
		argoApi: aApi,
		configs: c,
	}

	aApi.EXPECT().GetApps(c.Selectors, true).Return(apps, nil)

	r, _ := releaser.IsInSync()
	assert.Equal(t, false, r)
}

func TestIsInSyncIsFalseAndIgnored(t *testing.T) {
	ctrl := gomock.NewController(t)
	aApi := mocks.NewMockIArgoApi(ctrl)
	defer ctrl.Finish()
	apps := &argoApi.Applications{
		Items: []argoApi.Application{{
			Metadata: argoApi.AppMeta{
				Name: "test",
			},
			Status: argoApi.AppStatus{
				Sync: argoApi.AppStatusSync{Status: "OutOfSync"},
			},
		}},
	}
	c := &config.Config{
		Selectors: map[string]string{"foo": "bar"},
		Ignore:    []string{"test"},
	}
	releaser := &Releaser{
		argoApi: aApi,
		configs: c,
	}

	aApi.EXPECT().GetApps(c.Selectors, true).Return(apps, nil)

	r, _ := releaser.IsInSync()
	assert.Equal(t, true, r)
}

func TestIsInSyncIsFalseAndNotAllowed(t *testing.T) {
	ctrl := gomock.NewController(t)
	aApi := mocks.NewMockIArgoApi(ctrl)
	defer ctrl.Finish()
	apps := &argoApi.Applications{
		Items: []argoApi.Application{{
			Metadata: argoApi.AppMeta{
				Name: "test",
			},
			Status: argoApi.AppStatus{
				Sync: argoApi.AppStatusSync{Status: "OutOfSync"},
			},
		}},
	}
	c := &config.Config{
		Selectors: map[string]string{"foo": "bar"},
	}
	releaser := &Releaser{
		argoApi: aApi,
		configs: c,
	}

	aApi.EXPECT().GetApps(c.Selectors, true).Return(apps, nil)

	r, _ := releaser.IsInSync()
	assert.Equal(t, true, r)
}
