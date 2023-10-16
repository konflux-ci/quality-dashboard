package client

import (
	"context"

	v1alphaPlugins "github.com/redhat-appstudio/quality-studio/api/apis/plugins/v1alpha1"
	"github.com/redhat-appstudio/quality-studio/pkg/storage/ent/db"
	"github.com/redhat-appstudio/quality-studio/pkg/storage/ent/db/plugins"
)

// ListPlugins extracts an array of plugins from the database.
func (d *Database) ListPlugins() ([]*db.Plugins, error) {
	plugins, err := d.client.Plugins.Query().All(context.TODO())
	if err != nil {
		return nil, convertDBError("list plugins: %w", err)
	}

	return plugins, nil
}

// GetPluginByName extract the plugin name fromt database.
func (d *Database) GetPluginByName(pluginName string) (*db.Plugins, error) {
	plugin, err := d.client.Plugins.Query().Where(plugins.Name(pluginName)).Only(context.Background())
	if err != nil {
		return nil, convertDBError("get plugin: %w", err)
	}

	return plugin, nil
}

// GetPluginByName extract the plugin name fromt database.
func (d *Database) GetPluginsByTeam(team *db.Teams) ([]*db.Plugins, error) {
	plugins, err := d.client.Teams.QueryPlugins(team).All(context.Background())

	if err != nil {
		return nil, convertDBError("get plugin: %w", err)
	}

	return plugins, nil
}

func (d *Database) CreatePlugin(plugin *v1alphaPlugins.Plugin) (*db.Plugins, error) {
	team, err := d.client.Plugins.Create().
		SetName(plugin.Name).
		SetCategory(plugin.Category).
		SetLogo(plugin.Logo).
		SetDescription(plugin.Description).
		SetStatus(plugin.Status).
		Save(context.Background())
	if err != nil {
		return nil, convertDBError("create plugins status: %w", err)
	}

	return team, nil
}

func (d *Database) InstallPlugin(team *db.Teams, plugin *db.Plugins) (db *db.Teams, err error) {
	return d.client.Teams.UpdateOne(team).AddPlugins(plugin).Save(context.Background())
}
