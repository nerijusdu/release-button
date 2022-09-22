package releaser

import (
	"nerijusdu/release-button/internal/api"
	"nerijusdu/release-button/internal/config"
	"nerijusdu/release-button/internal/mocks"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestIsInSyncIsTrue(t *testing.T) {
	ctrl := gomock.NewController(t)
	argoApi := mocks.NewMockIArgoApi(ctrl)
	defer ctrl.Finish()
	apps := &api.Applications{}
	c := &config.Config{
		Selectors: map[string]string{"foo": "bar"},
	}
	releaser := &Releaser{
		argoApi: argoApi,
		configs: c,
	}

	argoApi.EXPECT().GetApps(c.Selectors, true).Return(apps, nil)

	r, _ := releaser.IsInSync()
	assert.Equal(t, true, r)
}

func TestIsInSyncIsFalse(t *testing.T) {
	ctrl := gomock.NewController(t)
	argoApi := mocks.NewMockIArgoApi(ctrl)
	defer ctrl.Finish()
	apps := &api.Applications{
		Items: []api.Application{{
			Status: api.AppStatus{
				Sync: api.AppStatusSync{Status: "OutOfSync"},
			},
		}},
	}
	c := &config.Config{
		Selectors: map[string]string{"foo": "bar"},
	}
	releaser := &Releaser{
		argoApi: argoApi,
		configs: c,
	}

	argoApi.EXPECT().GetApps(c.Selectors, true).Return(apps, nil)

	r, _ := releaser.IsInSync()
	assert.Equal(t, false, r)
}

func TestIsInSyncIsFalseAndIgnored(t *testing.T) {
	ctrl := gomock.NewController(t)
	argoApi := mocks.NewMockIArgoApi(ctrl)
	defer ctrl.Finish()
	apps := &api.Applications{
		Items: []api.Application{{
			Metadata: api.AppMeta{
				Name: "test",
			},
			Status: api.AppStatus{
				Sync: api.AppStatusSync{Status: "OutOfSync"},
			},
		}},
	}
	c := &config.Config{
		Selectors: map[string]string{"foo": "bar"},
		Ignore:    []string{"test"},
	}
	releaser := &Releaser{
		argoApi: argoApi,
		configs: c,
	}

	argoApi.EXPECT().GetApps(c.Selectors, true).Return(apps, nil)

	r, _ := releaser.IsInSync()
	assert.Equal(t, true, r)
}
