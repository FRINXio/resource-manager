// Code generated by entc, DO NOT EDIT.

// +build tools

// Package internal holds a loadable version of the latest schema.
package internal

const Schema = `{"Schema":"github.com/net-auto/resourceManager/ent/schema","Package":"github.com/net-auto/resourceManager/ent","Schemas":[{"name":"AllocationStrategy","config":{"Table":""},"edges":[{"name":"pools","type":"ResourcePool","ref_name":"allocation_strategy","inverse":true}],"fields":[{"name":"name","type":{"Type":7,"Ident":"","PkgPath":"","Nillable":false,"RType":null},"unique":true,"validators":1,"position":{"Index":0,"MixedIn":false,"MixinIndex":0}},{"name":"description","type":{"Type":7,"Ident":"","PkgPath":"","Nillable":false,"RType":null},"size":2147483647,"nillable":true,"optional":true,"position":{"Index":1,"MixedIn":false,"MixinIndex":0}},{"name":"lang","type":{"Type":6,"Ident":"allocationstrategy.Lang","PkgPath":"","Nillable":false,"RType":null},"enums":[{"N":"py","V":"py"},{"N":"js","V":"js"},{"N":"go","V":"go"}],"default":true,"default_value":"js","position":{"Index":2,"MixedIn":false,"MixinIndex":0}},{"name":"script","type":{"Type":7,"Ident":"","PkgPath":"","Nillable":false,"RType":null},"size":2147483647,"validators":1,"position":{"Index":3,"MixedIn":false,"MixinIndex":0}}],"policy":[{"Index":0,"MixedIn":false,"MixinIndex":0}]},{"name":"PoolProperties","config":{"Table":""},"edges":[{"name":"pool","type":"ResourcePool","ref_name":"poolProperties","unique":true,"inverse":true},{"name":"resourceType","type":"ResourceType","annotations":{"EntGQL":{"Bind":true,"Mapping":null,"OrderField":""}}},{"name":"properties","type":"Property","annotations":{"EntGQL":{"Bind":true,"Mapping":null,"OrderField":""}}}],"indexes":[{"edges":["pool"]}],"policy":[{"Index":0,"MixedIn":false,"MixinIndex":0}]},{"name":"Property","config":{"Table":""},"edges":[{"name":"type","type":"PropertyType","tag":"gqlgen:\"propertyType\"","unique":true,"required":true,"annotations":{"EntGQL":{"Bind":true,"Mapping":null,"OrderField":""}}},{"name":"resources","type":"Resource","ref_name":"properties","unique":true,"inverse":true}],"fields":[{"name":"int_val","type":{"Type":11,"Ident":"","PkgPath":"","Nillable":false,"RType":null},"tag":"json:\"intValue\" gqlgen:\"intValue\"","nillable":true,"optional":true,"position":{"Index":0,"MixedIn":false,"MixinIndex":0}},{"name":"bool_val","type":{"Type":1,"Ident":"","PkgPath":"","Nillable":false,"RType":null},"tag":"json:\"booleanValue\" gqlgen:\"booleanValue\"","nillable":true,"optional":true,"position":{"Index":1,"MixedIn":false,"MixinIndex":0}},{"name":"float_val","type":{"Type":19,"Ident":"","PkgPath":"","Nillable":false,"RType":null},"tag":"json:\"floatValue\" gqlgen:\"floatValue\"","nillable":true,"optional":true,"position":{"Index":2,"MixedIn":false,"MixinIndex":0}},{"name":"latitude_val","type":{"Type":19,"Ident":"","PkgPath":"","Nillable":false,"RType":null},"tag":"json:\"latitudeValue\" gqlgen:\"latitudeValue\"","nillable":true,"optional":true,"position":{"Index":3,"MixedIn":false,"MixinIndex":0}},{"name":"longitude_val","type":{"Type":19,"Ident":"","PkgPath":"","Nillable":false,"RType":null},"tag":"json:\"longitudeValue\" gqlgen:\"longitudeValue\"","nillable":true,"optional":true,"position":{"Index":4,"MixedIn":false,"MixinIndex":0}},{"name":"range_from_val","type":{"Type":19,"Ident":"","PkgPath":"","Nillable":false,"RType":null},"tag":"json:\"rangeFromValue\" gqlgen:\"rangeFromValue\"","nillable":true,"optional":true,"position":{"Index":5,"MixedIn":false,"MixinIndex":0}},{"name":"range_to_val","type":{"Type":19,"Ident":"","PkgPath":"","Nillable":false,"RType":null},"tag":"json:\"rangeToValue\" gqlgen:\"rangeToValue\"","nillable":true,"optional":true,"position":{"Index":6,"MixedIn":false,"MixinIndex":0}},{"name":"string_val","type":{"Type":7,"Ident":"","PkgPath":"","Nillable":false,"RType":null},"tag":"json:\"stringValue\" gqlgen:\"stringValue\"","nillable":true,"optional":true,"position":{"Index":7,"MixedIn":false,"MixinIndex":0}}],"indexes":[{"edges":["resources"]},{"edges":["type"]}],"policy":[{"Index":0,"MixedIn":false,"MixinIndex":0}]},{"name":"PropertyType","config":{"Table":""},"edges":[{"name":"properties","type":"Property","ref_name":"type","inverse":true},{"name":"resource_type","type":"ResourceType","ref_name":"property_types","unique":true,"inverse":true}],"fields":[{"name":"type","type":{"Type":6,"Ident":"propertytype.Type","PkgPath":"","Nillable":false,"RType":null},"enums":[{"N":"string","V":"string"},{"N":"int","V":"int"},{"N":"bool","V":"bool"},{"N":"float","V":"float"},{"N":"date","V":"date"},{"N":"enum","V":"enum"},{"N":"range","V":"range"},{"N":"email","V":"email"},{"N":"gps_location","V":"gps_location"},{"N":"datetime_local","V":"datetime_local"},{"N":"node","V":"node"}],"position":{"Index":0,"MixedIn":false,"MixinIndex":0}},{"name":"name","type":{"Type":7,"Ident":"","PkgPath":"","Nillable":false,"RType":null},"position":{"Index":1,"MixedIn":false,"MixinIndex":0}},{"name":"external_id","type":{"Type":7,"Ident":"","PkgPath":"","Nillable":false,"RType":null},"unique":true,"optional":true,"position":{"Index":2,"MixedIn":false,"MixinIndex":0}},{"name":"index","type":{"Type":11,"Ident":"","PkgPath":"","Nillable":false,"RType":null},"optional":true,"position":{"Index":3,"MixedIn":false,"MixinIndex":0}},{"name":"category","type":{"Type":7,"Ident":"","PkgPath":"","Nillable":false,"RType":null},"optional":true,"position":{"Index":4,"MixedIn":false,"MixinIndex":0}},{"name":"int_val","type":{"Type":11,"Ident":"","PkgPath":"","Nillable":false,"RType":null},"tag":"json:\"intValue\" gqlgen:\"intValue\"","nillable":true,"optional":true,"position":{"Index":5,"MixedIn":false,"MixinIndex":0}},{"name":"bool_val","type":{"Type":1,"Ident":"","PkgPath":"","Nillable":false,"RType":null},"tag":"json:\"booleanValue\" gqlgen:\"booleanValue\"","nillable":true,"optional":true,"position":{"Index":6,"MixedIn":false,"MixinIndex":0}},{"name":"float_val","type":{"Type":19,"Ident":"","PkgPath":"","Nillable":false,"RType":null},"tag":"json:\"floatValue\" gqlgen:\"floatValue\"","nillable":true,"optional":true,"position":{"Index":7,"MixedIn":false,"MixinIndex":0}},{"name":"latitude_val","type":{"Type":19,"Ident":"","PkgPath":"","Nillable":false,"RType":null},"tag":"json:\"latitudeValue\" gqlgen:\"latitudeValue\"","nillable":true,"optional":true,"position":{"Index":8,"MixedIn":false,"MixinIndex":0}},{"name":"longitude_val","type":{"Type":19,"Ident":"","PkgPath":"","Nillable":false,"RType":null},"tag":"json:\"longitudeValue\" gqlgen:\"longitudeValue\"","nillable":true,"optional":true,"position":{"Index":9,"MixedIn":false,"MixinIndex":0}},{"name":"string_val","type":{"Type":7,"Ident":"","PkgPath":"","Nillable":false,"RType":null},"tag":"json:\"stringValue\" gqlgen:\"stringValue\"","size":2147483647,"nillable":true,"optional":true,"position":{"Index":10,"MixedIn":false,"MixinIndex":0}},{"name":"range_from_val","type":{"Type":19,"Ident":"","PkgPath":"","Nillable":false,"RType":null},"tag":"json:\"rangeFromValue\" gqlgen:\"rangeFromValue\"","nillable":true,"optional":true,"position":{"Index":11,"MixedIn":false,"MixinIndex":0}},{"name":"range_to_val","type":{"Type":19,"Ident":"","PkgPath":"","Nillable":false,"RType":null},"tag":"json:\"rangeToValue\" gqlgen:\"rangeToValue\"","nillable":true,"optional":true,"position":{"Index":12,"MixedIn":false,"MixinIndex":0}},{"name":"is_instance_property","type":{"Type":1,"Ident":"","PkgPath":"","Nillable":false,"RType":null},"tag":"gqlgen:\"isInstanceProperty\"","default":true,"default_value":true,"position":{"Index":13,"MixedIn":false,"MixinIndex":0}},{"name":"editable","type":{"Type":1,"Ident":"","PkgPath":"","Nillable":false,"RType":null},"tag":"gqlgen:\"isEditable\"","default":true,"default_value":true,"position":{"Index":14,"MixedIn":false,"MixinIndex":0}},{"name":"mandatory","type":{"Type":1,"Ident":"","PkgPath":"","Nillable":false,"RType":null},"tag":"gqlgen:\"isMandatory\"","default":true,"default_value":false,"position":{"Index":15,"MixedIn":false,"MixinIndex":0}},{"name":"deleted","type":{"Type":1,"Ident":"","PkgPath":"","Nillable":false,"RType":null},"tag":"gqlgen:\"isDeleted\"","default":true,"default_value":false,"position":{"Index":16,"MixedIn":false,"MixinIndex":0}},{"name":"nodeType","type":{"Type":7,"Ident":"","PkgPath":"","Nillable":false,"RType":null},"optional":true,"position":{"Index":17,"MixedIn":false,"MixinIndex":0}}],"policy":[{"Index":0,"MixedIn":false,"MixinIndex":0}]},{"name":"Resource","config":{"Table":""},"edges":[{"name":"pool","type":"ResourcePool","ref_name":"claims","unique":true,"inverse":true},{"name":"properties","type":"Property","annotations":{"EntGQL":{"Bind":true,"Mapping":null,"OrderField":""}}},{"name":"nested_pool","type":"ResourcePool","unique":true,"annotations":{"EntGQL":{"Bind":true,"Mapping":null,"OrderField":""}}}],"fields":[{"name":"status","type":{"Type":6,"Ident":"resource.Status","PkgPath":"","Nillable":false,"RType":null},"enums":[{"N":"free","V":"free"},{"N":"claimed","V":"claimed"},{"N":"retired","V":"retired"},{"N":"bench","V":"bench"}],"position":{"Index":0,"MixedIn":false,"MixinIndex":0}},{"name":"description","type":{"Type":7,"Ident":"","PkgPath":"","Nillable":false,"RType":null},"size":2147483647,"nillable":true,"optional":true,"position":{"Index":1,"MixedIn":false,"MixinIndex":0}},{"name":"alternate_id","type":{"Type":3,"Ident":"map[string]interface {}","PkgPath":"","Nillable":true,"RType":null},"optional":true,"position":{"Index":2,"MixedIn":false,"MixinIndex":0}},{"name":"updated_at","type":{"Type":2,"Ident":"","PkgPath":"time","Nillable":false,"RType":null},"default":true,"update_default":true,"position":{"Index":3,"MixedIn":false,"MixinIndex":0}}],"indexes":[{"edges":["pool"]}],"policy":[{"Index":0,"MixedIn":false,"MixinIndex":0}]},{"name":"ResourcePool","config":{"Table":""},"edges":[{"name":"resource_type","type":"ResourceType","ref_name":"pools","unique":true,"inverse":true},{"name":"tags","type":"Tag","ref_name":"pools","inverse":true},{"name":"claims","type":"Resource","annotations":{"EntGQL":{"Bind":true,"Mapping":null,"OrderField":""}}},{"name":"poolProperties","type":"PoolProperties","unique":true,"annotations":{"EntGQL":{"Bind":true,"Mapping":null,"OrderField":""}}},{"name":"allocation_strategy","type":"AllocationStrategy","unique":true,"annotations":{"EntGQL":{"Bind":true,"Mapping":null,"OrderField":""}}},{"name":"parent_resource","type":"Resource","ref_name":"nested_pool","unique":true,"inverse":true}],"fields":[{"name":"name","type":{"Type":7,"Ident":"","PkgPath":"","Nillable":false,"RType":null},"unique":true,"validators":1,"position":{"Index":0,"MixedIn":false,"MixinIndex":0}},{"name":"description","type":{"Type":7,"Ident":"","PkgPath":"","Nillable":false,"RType":null},"size":2147483647,"nillable":true,"optional":true,"position":{"Index":1,"MixedIn":false,"MixinIndex":0}},{"name":"pool_type","type":{"Type":6,"Ident":"resourcepool.PoolType","PkgPath":"","Nillable":false,"RType":null},"enums":[{"N":"singleton","V":"singleton"},{"N":"set","V":"set"},{"N":"allocating","V":"allocating"}],"position":{"Index":2,"MixedIn":false,"MixinIndex":0}},{"name":"dealocation_safety_period","type":{"Type":11,"Ident":"","PkgPath":"","Nillable":false,"RType":null},"default":true,"default_value":0,"position":{"Index":3,"MixedIn":false,"MixinIndex":0}}],"indexes":[{"edges":["allocation_strategy"]}],"policy":[{"Index":0,"MixedIn":false,"MixinIndex":0}]},{"name":"ResourceType","config":{"Table":""},"edges":[{"name":"property_types","type":"PropertyType","annotations":{"EntGQL":{"Bind":true,"Mapping":null,"OrderField":""}}},{"name":"pools","type":"ResourcePool","annotations":{"EntGQL":{"Bind":true,"Mapping":null,"OrderField":""}}},{"name":"pool_properties","type":"PoolProperties","ref_name":"resourceType","inverse":true}],"fields":[{"name":"name","type":{"Type":7,"Ident":"","PkgPath":"","Nillable":false,"RType":null},"validators":1,"position":{"Index":0,"MixedIn":false,"MixinIndex":0}}],"policy":[{"Index":0,"MixedIn":false,"MixinIndex":0}]},{"name":"Tag","config":{"Table":""},"edges":[{"name":"pools","type":"ResourcePool","annotations":{"EntGQL":{"Bind":true,"Mapping":null,"OrderField":""}}}],"fields":[{"name":"tag","type":{"Type":7,"Ident":"","PkgPath":"","Nillable":false,"RType":null},"unique":true,"validators":1,"position":{"Index":0,"MixedIn":false,"MixinIndex":0}}]}],"Features":["privacy","schema/snapshot"]}`
