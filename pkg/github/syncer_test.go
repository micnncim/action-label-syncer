package github

import (
	"context"
	"testing"

	"github.com/micnncim/action-label-syncer/pkg/github/fake"
)

func TestLabelSyncer_SyncLabels(t *testing.T) {
	type fields struct {
		client Client
	}
	type args struct {
		ctx    context.Context
		owner  string
		repo   string
		labels []Label
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "create labels",
			fields: fields{
				client: &fake.Client{
					FakeCreateLabel: func(_ context.Context, owner, repo string, label Label) error {
						return nil
					},
					FakeGetLabels:   nil,
					FakeUpdateLabel: nil,
					FakeDeleteLabel: nil,
				},
			},
			args: args{
				ctx:    nil,
				owner:  "",
				repo:   "",
				labels: nil,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &LabelSyncer{
				client: tt.fields.client,
			}
			if err := s.SyncLabels(tt.args.ctx, tt.args.owner, tt.args.repo, tt.args.labels); (err != nil) != tt.wantErr {
				t.Errorf("SyncLabels() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
