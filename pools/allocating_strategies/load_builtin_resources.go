package pools

import (
	"context"
	"github.com/net-auto/resourceManager/ent"
	"github.com/net-auto/resourceManager/ent/allocationstrategy"
	"github.com/net-auto/resourceManager/ent/propertytype"
	"github.com/net-auto/resourceManager/ent/resourcetype"
	"github.com/pkg/errors"
)

func loadIpv4Prefix(ctx context.Context, client *ent.Tx) error {

	exists, err := client.ResourceType.Query().Where(resourcetype.Name("ipv4_prefix")).Exist(ctx)
	if err != nil {
		return err
	}
	if exists {
		// TODO update if exists
		// TODO prevent users from overriding these
		return nil
	}

	propAddr, err := client.PropertyType.Create().
		SetName("address").
		SetType(propertytype.TypeString).
		Save(ctx)
	if err != nil {
		return err
	}
	propPrefix, err := client.PropertyType.Create().
		SetName("prefix").
		SetType(propertytype.TypeInt).
		Save(ctx)
	if err != nil {
		return err
	}

	_, err = client.ResourceType.Create().
		SetName("ipv4_prefix").
		AddPropertyTypes(propAddr, propPrefix).
		Save(ctx)
	if err != nil {
		return err
	}

	_, err = client.AllocationStrategy.Create().
		SetName("ipv4_prefix").
		SetLang(allocationstrategy.LangJs).
		SetScript(IPV4_PREFIX).
		Save(ctx)
	if err != nil {
		return err
	}

	return nil
}

func loadIpv6Prefix(ctx context.Context, client *ent.Tx) error {

	exists, err := client.ResourceType.Query().Where(resourcetype.Name("ipv6_prefix")).Exist(ctx)
	if err != nil {
		return err
	}
	if exists {
		// TODO update if exists
		// TODO prevent users from overriding these
		return nil
	}

	propAddr, err := client.PropertyType.Create().
		SetName("address").
		SetType(propertytype.TypeString).
		Save(ctx)
	if err != nil {
		return err
	}
	propPrefix, err := client.PropertyType.Create().
		SetName("prefix").
		SetType(propertytype.TypeInt).
		Save(ctx)
	if err != nil {
		return err
	}

	_, err = client.ResourceType.Create().
		SetName("ipv6_prefix").
		AddPropertyTypes(propAddr, propPrefix).
		Save(ctx)
	if err != nil {
		return err
	}

	_, err = client.AllocationStrategy.Create().
		SetName("ipv6_prefix").
		SetLang(allocationstrategy.LangJs).
		SetScript(IPV6_PREFIX).
		Save(ctx)
	if err != nil {
		return err
	}

	return nil
}

func loadIpv4(ctx context.Context, client *ent.Tx) error {
	exists, err := client.ResourceType.Query().Where(resourcetype.Name("ipv4")).Exist(ctx)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}

	propAddr, err := client.PropertyType.Create().
		SetName("address").
		SetType(propertytype.TypeString).
		Save(ctx)
	if err != nil {
		return err
	}

	_, err = client.ResourceType.Create().
		SetName("ipv4").
		AddPropertyTypes(propAddr).
		Save(ctx)
	if err != nil {
		return err
	}

	_, err = client.AllocationStrategy.Create().
		SetName("ipv4").
		SetLang(allocationstrategy.LangJs).
		SetScript(IPV4).
		Save(ctx)
	if err != nil {
		return err
	}

	return nil
}

func loadIpv6(ctx context.Context, client *ent.Tx) error {
	exists, err := client.ResourceType.Query().Where(resourcetype.Name("ipv6")).Exist(ctx)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}

	propAddr, err := client.PropertyType.Create().
		SetName("address").
		SetType(propertytype.TypeString).
		Save(ctx)
	if err != nil {
		return err
	}

	_, err = client.ResourceType.Create().
		SetName("ipv6").
		AddPropertyTypes(propAddr).
		Save(ctx)
	if err != nil {
		return err
	}

	_, err = client.AllocationStrategy.Create().
		SetName("ipv6").
		SetLang(allocationstrategy.LangJs).
		SetScript(IPV6).
		Save(ctx)
	if err != nil {
		return err
	}

	return nil
}

func loadVlanRange(ctx context.Context, client *ent.Tx) error {

	exists, err := client.ResourceType.Query().Where(resourcetype.Name("vlan_range")).Exist(ctx)
	if err != nil {
		return err
	}
	if exists {
		// TODO update if exists
		// TODO prevent users from overriding these
		return nil
	}

	propFrom, err := client.PropertyType.Create().
		SetName("from").
		SetType(propertytype.TypeInt).
		Save(ctx)
	if err != nil {
		return err
	}

	propTo, err := client.PropertyType.Create().
		SetName("to").
		SetType(propertytype.TypeInt).
		Save(ctx)
	if err != nil {
		return err
	}

	_, err = client.ResourceType.Create().
		SetName("vlan_range").
		AddPropertyTypes(propFrom, propTo).
		Save(ctx)
	if err != nil {
		return err
	}

	_, err = client.AllocationStrategy.Create().
		SetName("vlan_range").
		SetLang(allocationstrategy.LangJs).
		SetScript(VLAN_RANGE).
		Save(ctx)
	if err != nil {
		return err
	}

	return nil
}

func loadVlan(ctx context.Context, client *ent.Tx) error {

	exists, err := client.ResourceType.Query().Where(resourcetype.Name("vlan")).Exist(ctx)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}

	propAddr, err := client.PropertyType.Create().
		SetName("vlan").
		SetType(propertytype.TypeInt).
		Save(ctx)
	if err != nil {
		return err
	}

	_, err = client.ResourceType.Create().
		SetName("vlan").
		AddPropertyTypes(propAddr).
		Save(ctx)
	if err != nil {
		return err
	}

	_, err = client.AllocationStrategy.Create().
		SetName("vlan").
		SetLang(allocationstrategy.LangJs).
		SetScript(VLAN).
		Save(ctx)
	if err != nil {
		return err
	}

	return nil
}

func loadInner(ctx context.Context, client *ent.Tx) error {
	err := loadIpv4Prefix(ctx, client)
	if err != nil {
		return errors.Wrapf(err, "Unable to load ipv4_prefix resource type")
	}
	err = loadIpv4(ctx, client)
	if err != nil {
		return errors.Wrapf(err, "Unable to load ipv4 resource type")
	}

	err = loadVlanRange(ctx, client)
	if err != nil {
		return errors.Wrapf(err, "Unable to load vlan_range resource type")
	}
	err = loadVlan(ctx, client)
	if err != nil {
		return errors.Wrapf(err, "Unable to load vlan resource type")
	}

	err = loadIpv6Prefix(ctx, client)
	if err != nil {
		return errors.Wrapf(err, "Unable to load ipv6_prefix resource type")
	}
	err = loadIpv6(ctx, client)
	if err != nil {
		return errors.Wrapf(err, "Unable to load ipv6 resource type")
	}

	return nil
}

// LoadBuiltinTypes loads IP, VLAN etc. resource types and allocation strategies into DB
//  does not overwrite existing resources and strategies
func LoadBuiltinTypes(ctx context.Context, client *ent.Client) error {
	tx, err := client.Tx(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if v := recover(); v != nil {
			tx.Rollback()
			panic(v)
		}
	}()
	if err := loadInner(ctx, tx); err != nil {
		if rerr := tx.Rollback(); rerr != nil {
			err = errors.Wrapf(err, "rolling back transaction: %v", rerr)
		}
		return err
	}
	if err := tx.Commit(); err != nil {
		return errors.Wrapf(err, "committing transaction: %v", err)
	}
	return nil
}

