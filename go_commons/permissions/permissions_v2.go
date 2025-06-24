package permissions

const (
	// Order permissions
	OrderView           PermissionType = "view_order"
	OrderCreate         PermissionType = "create_order"
	OrderApprove        PermissionType = "approve_order"
	OrderCreateShipment PermissionType = "create_shipment_order"
	OrderCancel         PermissionType = "cancel_order"
	OrderViewLog        PermissionType = "view_log_order"
	OrderCreateTag      PermissionType = "create_tag_order"
	OrderEditTag        PermissionType = "edit_tag_order"
	OrderViewTag        PermissionType = "view_tag_order"
	OrderEditStatus     PermissionType = "edit_status_order"
	OrderEditDetails    PermissionType = "edit_details_order"
	OrderForceFulRoute  PermissionType = "forceful_route_order"
	OrderViewTimeline   PermissionType = "view_timeline_order"
	OrderReassignHub    PermissionType = "reassign_hub_order"

	// On-Hold Order Permissions
	OnHoldOrderView PermissionType = "view_on_hold_order"
	OrderPutOnHold  PermissionType = "put_on_hold_order"

	// Return Order Permissions
	ReturnCreateRequest  PermissionType = "create_request_return"
	ReturnViewRequest    PermissionType = "view_request_return"
	ReturnEditRequest    PermissionType = "edit_request_return"
	ReturnCreateOrder    PermissionType = "create_order_return"
	ReturnEditOrder      PermissionType = "edit_order_return"
	ReturnCancelOrder    PermissionType = "cancel_order_return"
	ReturnViewOrder      PermissionType = "view_order_return"
	ReturnQCOrder        PermissionType = "qc_order_return"
	ReturnRefundOrder    PermissionType = "refund_order_return"
	ReturnCreateShipment PermissionType = "create_shipment_return"

	// Split Order Permissions
	SplitOrderCreate PermissionType = "create_split_order"
	SplitOrderEdit   PermissionType = "edit_split_order"
	SplitOrderView   PermissionType = "view_split_order"

	// Customer Permissions
	CustomerCreate PermissionType = "create_customer"
	CustomerEdit   PermissionType = "edit_customer"
	CustomerView   PermissionType = "view_customer"
	CustomerDelete PermissionType = "delete_customer"

	// Sales Person Permissions
	SalesPersonCreate PermissionType = "create_sales_person"
	SalesPersonEdit   PermissionType = "edit_sales_person"
	SalesPersonView   PermissionType = "view_sales_person"

	// Address Permissions
	CustomerCreateAddress PermissionType = "create_address_customer"
	CustomerViewAddress   PermissionType = "view_address_customer"
	CustomerEditAddress   PermissionType = "edit_address_customer"

	// Shipment Permissions
	ShipmentEditStatus PermissionType = "edit_status_shipment"
	ShipmentView       PermissionType = "view_shipment"

	// Configurations permissions
	ConfigurationCreate            PermissionType = "create_configuration"
	ConfigurationView              PermissionType = "view_configuration"
	ConfigurationOrderNotification PermissionType = "order_notification_configuration"

	// Manifest permissions
	ManifestView PermissionType = "view_manifest"
	ManifestEdit PermissionType = "edit_manifest"

	// STO permissions
	STOCreate  PermissionType = "create_sto"
	STOView    PermissionType = "view_sto"
	STOApprove PermissionType = "approve_sto"
	STOEdit    PermissionType = "edit_sto"
	STOReceive PermissionType = "receive_sto"

	// Hub permissions
	HubCreate            PermissionType = "create_hub"
	HubEdit              PermissionType = "edit_hub"
	HubView              PermissionType = "view_hub"
	HubEditConfiguration PermissionType = "edit_configuration_hub"
	HubViewConfiguration PermissionType = "view_configuration_hub"

	// Hub Location permissions
	HubLocationCreate         PermissionType = "create_hub_location"
	HubLocationEdit           PermissionType = "edit_hub_location"
	HubLocationView           PermissionType = "view_hub_location"
	HubLocationCreateBin      PermissionType = "create_bin_hub_location"
	HubLocationViewBin        PermissionType = "view_bin_hub_location"
	HubLocationViewSkuMapping PermissionType = "view_sku_mapping_hub_location"
	HubLocationEditSkuMapping PermissionType = "edit_sku_mapping_hub_location"

	//Sku Permissions
	SkuConfigurationSet  PermissionType = "set_sku_configuration"
	SkuConfigurationView PermissionType = "view_sku_configuration"

	// Product permissions
	ProductCreate PermissionType = "create_product"
	ProductView   PermissionType = "view_product"

	// Weighted Average Cost Permissions
	WeightedAverageCostView PermissionType = "view_weighted_average_cost"

	// Catalog Comparison permissions
	CatalogComparison        PermissionType = "compare_catalog"
	CatalogViewConfiguration PermissionType = "view_configuration_catalog"
	CatalogEditConfiguration PermissionType = "edit_configuration_catalog"

	// Purchase Order permissions
	PurchaseOrderCreate PermissionType = "create_purchase_order"
	PurchaseOrderView   PermissionType = "view_purchase_order"
	PurchaseOrderDelete PermissionType = "delete_purchase_order"

	// Billing Permissions
	BillingCreateProfile     PermissionType = "create_profile_billing"
	BillingViewProfile       PermissionType = "view_profile_billing"
	BillingEditProfile       PermissionType = "edit_profile_billing"
	BillingCreateContract    PermissionType = "create_contract_billing"
	BillingViewContract      PermissionType = "view_contract_billing"
	BillingTerminateContract PermissionType = "terminate_contract_billing"
	BillingCreateBill        PermissionType = "create_bill_billing"
	BillingViewBill          PermissionType = "view_bill_billing"
	BillingEditBill          PermissionType = "edit_bill_billing"
	BillingFinalizeBill      PermissionType = "finalize_bill_billing"
	BillingGenerateBill      PermissionType = "generate_bill_billing"

	//Adhoc Permissions
	AdhocCreate      PermissionType = "create_adhoc"
	AdhocView        PermissionType = "view_adhoc"
	AdhocCreateUsage PermissionType = "create_usage_adhoc"
	AdhocViewUsage   PermissionType = "view_usage_adhoc"
	AdhocExportUsage PermissionType = "export_usage_adhoc"

	// Supplier permissions
	SupplierCreate PermissionType = "create_supplier"
	SupplierEdit   PermissionType = "edit_supplier"
	SupplierView   PermissionType = "view_supplier"

	// Seller permissions
	SellerCreate PermissionType = "create_seller"
	SellerView   PermissionType = "view_seller"
	SellerEdit   PermissionType = "edit_seller"

	// Gate Entry permissions
	GateEntryCreate PermissionType = "create_gate_entry"
	GateEntryView   PermissionType = "view_gate_entry"

	// GRN permissions
	GrnCreate PermissionType = "create_grn"
	GrnView   PermissionType = "view_grn"

	// Inventory permissions
	InventoryView         PermissionType = "view_inventory"
	InventoryMovement     PermissionType = "movement_inventory"
	InventoryAdjust       PermissionType = "adjust_inventory"
	InventoryEdit         PermissionType = "edit_inventory"
	UploadInventory       PermissionType = "upload_inventory"
	ConfigureReOrderPoint PermissionType = "configure_reorder_point_inventory"
	ConfigureSafetyStock  PermissionType = "configure_safety_stock_inventory"
	ConfigureMaxShelfLife PermissionType = "configure_max_shelf_life_inventory"
	BreakCasePack         PermissionType = "break_case_pack_inventory"

	// User permissions
	UserCreate             PermissionType = "create_user"
	UserView               PermissionType = "view_user"
	UserEdit               PermissionType = "edit_user"
	UserResendPasswordLink PermissionType = "resend_password_link_user"

	//Omniful User permissions
	OmnifulUserCreate PermissionType = "create_omniful_user"
	OmnifulUserView   PermissionType = "view_omniful_user"

	// Role permissions
	RoleCreate PermissionType = "create_role"
	RoleView   PermissionType = "view_role"
	RoleEdit   PermissionType = "edit_role"

	// Sales Channel App permissions
	SalesChannelAppAdd  PermissionType = "add_sales_channel_app"
	SalesChannelAppView PermissionType = "view_sales_channel_app"
	SalesChannelAppEdit PermissionType = "edit_sales_channel_app"

	// Shipping App permissions
	ShippingAppAdd  PermissionType = "add_shipping_app"
	ShippingAppView PermissionType = "view_shipping_app"

	// Tenant City Mapping permissions
	TenantCityMappingAdd     PermissionType = "add_tenant_city_mapping"
	TenantCityMappingView    PermissionType = "view_tenant_city_mapping"
	TenantCityMappingRequest PermissionType = "request_tenant_city_mapping"

	// Shipping Rule permissions
	ShippingRuleAdd  PermissionType = "add_shipping_rule"
	ShippingRuleView PermissionType = "view_shipping_rule"
	ShippingRuleEdit PermissionType = "edit_shipping_rule"

	// Picking Wave permissions
	PickingCreateWave PermissionType = "create_wave_picking"
	PickingViewWave   PermissionType = "view_wave_picking"

	// Stock Ownership Transfer permissions
	StockOwnershipTransferCreate PermissionType = "create_stock_ownership_transfer"
	StockOwnershipTransferView   PermissionType = "view_stock_ownership_transfer"

	// BulkShip permissions
	BulkShipCreate           PermissionType = "create_bulk_ship"
	BulkShipView             PermissionType = "view_bulk_ship"
	BulkShipSaveSuggestion   PermissionType = "save_suggestion_bulk_ship"
	BulkShipDeleteSuggestion PermissionType = "delete_suggestion_bulk_ship"

	// PutAway permissions
	PutAwayView    PermissionType = "view_put_away"
	PutAwayPerform PermissionType = "perform_put_away"

	// Driver permissions
	//DriverView               PermissionType = "view_driver"
	//DriverCreate             PermissionType = "create_driver"
	//TripView                 PermissionType = "view_trip"
	//TripCreate               PermissionType = "create_trip"
	FleetConfigurationView   PermissionType = "view_fleet_configuration"
	FleetConfigurationCreate PermissionType = "create_fleet_configuration"

	// Reports permissions
	InwardingReport               PermissionType = "inwarding_report"
	FulfilmentReport              PermissionType = "fulfilment_report"
	ReturnReport                  PermissionType = "return_report"
	InventoryReport               PermissionType = "inventory_report"
	CommercialReport              PermissionType = "commercial_report"
	DeliveryReport                PermissionType = "delivery_report"
	PerformanceProductivityReport PermissionType = "performance_productivity_report"
	InvoicingReport               PermissionType = "invoicing_report"
	ViewReport                    PermissionType = "view_report"
	DownloadReport                PermissionType = "download_report"
	ValuationReport               PermissionType = "valuation_report"
	WorkforceProductivityReport   PermissionType = "workforce_productivity_report"
	ViewReportArchives            PermissionType = "view_archives_report"
	OmnishipReport                PermissionType = "omniship_report"
	POSReport                     PermissionType = "pos_report"

	GenerateForecast PermissionType = "generate_forecast"

	// Cycle Count Permissions
	CycleCountCreate PermissionType = "create_cycle_count"
	CycleCountView   PermissionType = "view_cycle_count"
	CycleCountCount  PermissionType = "count_cycle_count"

	//Assembly Permissions
	AssemblyView     PermissionType = "view_assembly"
	AssemblyAssemble PermissionType = "assemble_assembly"

	//Reasons Permissions
	ReasonCreate PermissionType = "create_reason"
	ReasonView   PermissionType = "view_reason"
	ReasonDelete PermissionType = "delete_reason"
	ReasonEdit   PermissionType = "edit_reason"

	// Dashboard Permissions
	DashboardViewHome            PermissionType = "view_home_dashboard"
	DashboardDownloadHome        PermissionType = "download_home_dashboard"
	DashboardViewSeller          PermissionType = "view_seller_dashboard"
	DashboardDownloadSeller      PermissionType = "download_seller_dashboard"
	DashboardViewHub             PermissionType = "view_hub_dashboard"
	DashboardDownloadHub         PermissionType = "download_hub_dashboard"
	DashboardViewShipping        PermissionType = "view_shipping_dashboard"
	DashboardDownloadShipping    PermissionType = "download_shipping_dashboard"
	DashboardViewFulfillment     PermissionType = "view_fulfillment_dashboard"
	DashboardDownloadFulfillment PermissionType = "download_fulfillment_dashboard"
	DashboardViewLastMile        PermissionType = "view_lastmile_dashboard"
	DashboardDownloadLastMile    PermissionType = "download_lastmile_dashboard"
	DashboardViewPOS             PermissionType = "view_pos_dashboard"
	DashboardDownloadPOS         PermissionType = "download_pos_dashboard"

	//Tenant Permissions
	TenantView PermissionType = "view_tenant"
	TenantEdit PermissionType = "edit_tenant"

	//Batch Permissions
	BatchEdit PermissionType = "edit_batch"

	//ShipmentOrder Permissions
	ShipmentOrderAdd          PermissionType = "add_shipment_order"
	ShipmentOrderView         PermissionType = "view_shipment_order"
	ShipmentOrderCancel       PermissionType = "cancel_shipment_order"
	ShipmentOrderViewLog      PermissionType = "view_log_shipment_order"
	ShipmentOrderEditStatus   PermissionType = "edit_status_shipment_order"
	ShipmentOrderEditDetails  PermissionType = "edit_details_shipment_order"
	ShipmentOrderBulkShipment PermissionType = "create_bulk_shipment_shipment_order"

	//ShipmentOrderV2 Permissions
	ShipmentOrderV2Add         PermissionType = "add_shipment_order_v2"
	ShipmentOrderV2View        PermissionType = "view_shipment_order_v2"
	ShipmentOrderV2Cancel      PermissionType = "cancel_shipment_order_v2"
	ShipmentOrderV2ViewLog     PermissionType = "view_log_shipment_order_v2"
	ShipmentOrderV2EditStatus  PermissionType = "edit_status_shipment_order_v2"
	ShipmentOrderV2EditDetails PermissionType = "edit_details_shipment_order_v2"

	// ShippingV2 Permissions
	ShippingRuleV2Add  PermissionType = "add_shipping_rule_v2"
	ShippingRuleV2View PermissionType = "view_shipping_rule_v2"
	ShippingRuleV2Edit PermissionType = "edit_shipping_rule_v2"

	// ShipmentV2 Permissions
	ShipmentV2Add     PermissionType = "add_shipment_v2"
	ShipmentV2View    PermissionType = "view_shipment_v2"
	ShipmentV2Edit    PermissionType = "edit_shipment_v2"
	ShipmentV2AWBView PermissionType = "view_awb_shipment_v2"
	AllShipmentV2View PermissionType = "view_all_shipment_v2"

	// Tax Authority Integrations
	TaxAuthorityCreate PermissionType = "create_tax_authority"
	TaxAuthorityView   PermissionType = "view_tax_authority"

	// Invoice Permissions
	InvoiceView      PermissionType = "view_invoice"
	InvoiceReport    PermissionType = "report_invoice"
	InvoiceConfigure PermissionType = "configure_invoice"

	// Product Price View Permissions
	ProductPriceView PermissionType = "view_price_product"

	//External WMS Permissions
	ExternalWMSView PermissionType = "view_external_wms"
	ExternalWMSEdit PermissionType = "edit_external_wms"

	//Reports APIs Permissions
	ReportsAPIsEdit PermissionType = "edit_reports_api"

	//Downloads Permissions
	DownloadOrders    PermissionType = "download_orders_dashboard"
	DownloadShipments PermissionType = "download_shipments_dashboard"

	// Priority Rule Permissions
	PriorityRuleAdd    PermissionType = "add_priority_rule"
	PriorityRuleView   PermissionType = "view_priority_rule"
	PriorityRuleEdit   PermissionType = "edit_priority_rule"
	PriorityRuleDelete PermissionType = "delete_priority_rule"

	// Delivery Zone Permissions
	DeliveryZoneView             PermissionType = "view_delivery_zone"
	DeliveryZoneCreate           PermissionType = "create_delivery_zone"
	DeliveryZoneUpdate           PermissionType = "update_delivery_zone"
	DeliveryZoneUpdatePrecedence PermissionType = "update_precedence_delivery_zone"

	/*
		TMS Permissions
	*/

	// Reports and Matrices

	TmsReportView PermissionType = "view_tms_report"
	TmsMatrixView PermissionType = "view_tms_matrix"

	// ShippingClient Permissions

	ShippingClientCreate PermissionType = "create_shipping_client"
	ShippingClientView   PermissionType = "view_shipping_client"
	ShippingClientDelete PermissionType = "delete_shipping_client"
	ShippingClientEdit   PermissionType = "edit_shipping_client"

	// ClientIdentifier Permissions

	ClientIdentifierCreate PermissionType = "create_client_identifier"
	ClientIdentifierView   PermissionType = "view_client_identifier"
	ClientIdentifierDelete PermissionType = "delete_client_identifier"
	ClientIdentifierEdit   PermissionType = "edit_client_identifier"

	// Webhook Permissions

	WebhookCreate PermissionType = "create_webhook"
	WebhookView   PermissionType = "view_webhook"
	WebhookDelete PermissionType = "delete_webhook"
	WebhookEdit   PermissionType = "edit_webhook"

	// WebhookLog Permissions

	WebhookLogView PermissionType = "view_webhook_log"

	// SortingHub Permissions

	SortingHubDelete PermissionType = "delete_sorting_hub"
	SortingHubCreate PermissionType = "create_sorting_hub"
	SortingHubView   PermissionType = "view_sorting_hub"
	SortingHubEdit   PermissionType = "edit_sorting_hub"

	// HubDriver Permissions

	HubDriverView PermissionType = "view_hub_driver"

	// Shipment Permissions

	ClientShipmentCreate PermissionType = "create_client_shipment"
	ClientShipmentView   PermissionType = "view_client_shipment"
	ClientShipmentDelete PermissionType = "delete_client_shipment"
	ClientShipmentEdit   PermissionType = "edit_client_shipment"

	// ClientCustomer Permissions

	ClientCustomerCreate PermissionType = "create_client_customer"
	ClientCustomerView   PermissionType = "view_client_customer"
	ClientCustomerEdit   PermissionType = "edit_client_customer"

	// Address Permissions

	AddressCreate PermissionType = "create_address"
	AddressView   PermissionType = "view_address"
	AddressDelete PermissionType = "delete_address"
	AddressEdit   PermissionType = "edit_address"

	// Package Permissions

	PackageCreate PermissionType = "create_package"
	PackageView   PermissionType = "view_package"
	PackageDelete PermissionType = "delete_package"
	PackageEdit   PermissionType = "edit_package"
	PackageReturn PermissionType = "return_package"

	// Fleet Permissions

	FleetCreate PermissionType = "create_fleet"
	FleetView   PermissionType = "view_fleet"
	FleetDelete PermissionType = "delete_fleet"
	FleetEdit   PermissionType = "edit_fleet"

	// FleetHub Permissions

	FleetHubView   PermissionType = "view_fleet_hub"
	FleetHubCreate PermissionType = "create_fleet_hub"

	// GeoZone Permissions

	GeoZoneCreate PermissionType = "create_geo_zone"
	GeoZoneView   PermissionType = "view_geo_zone"
	GeoZoneDelete PermissionType = "delete_geo_zone"
	GeoZoneEdit   PermissionType = "edit_geo_zone"

	// Driver Permissions

	DriverCreate  PermissionType = "create_driver"
	DriverView    PermissionType = "view_driver"
	DriverDelete  PermissionType = "delete_driver"
	DriverEdit    PermissionType = "edit_driver"
	DriverSuggest PermissionType = "suggest_driver"

	// DriverFleet Permissions

	DriverFleetCreate PermissionType = "create_driver_fleet"
	DriverFleetView   PermissionType = "view_driver_fleet"

	// DriverVehicle

	DriverVehicleEdit PermissionType = "edit_driver_vehicle"

	// Vehicle Permissions

	VehicleCreate PermissionType = "create_vehicle"
	VehicleView   PermissionType = "view_vehicle"
	VehicleDelete PermissionType = "delete_vehicle"
	VehicleEdit   PermissionType = "edit_vehicle"

	// Trip Permissions

	TripCreate         PermissionType = "create_trip"
	TripView           PermissionType = "view_trip"
	TripDelete         PermissionType = "delete_trip"
	TripEdit           PermissionType = "edit_trip"
	TripCancel         PermissionType = "cancel_trip"
	TripAssignPackages PermissionType = "assign_packages_trip"
	TripBroadcast      PermissionType = "broadcast_trip"

	// TripPackage Permissions

	TripPackageView PermissionType = "view_trip_package"

	// TripPlanning Permissions

	TripPlanningCreate PermissionType = "create_trip_planning"
	TripPlanningView   PermissionType = "view_trip_planning"
	TripPlanningDelete PermissionType = "delete_trip_planning"
	TripPlanningEdit   PermissionType = "edit_trip_planning"

	// TripAssignment Permissions

	TripAssignmentCreate        PermissionType = "create_trip_assignment"
	TripAssignmentView          PermissionType = "view_trip_assignment"
	TripAssignmentDelete        PermissionType = "delete_trip_assignment"
	TripAssignmentEdit          PermissionType = "edit_trip_assignment"
	TripAssignmentAssignDrivers PermissionType = "assign_drivers_trip_assignment"

	// ClientHub Permissions

	ClientHubCreate PermissionType = "create_client_hub"
	ClientHubView   PermissionType = "view_client_hub"
	ClientHubDelete PermissionType = "delete_client_hub"
	ClientHubEdit   PermissionType = "edit_client_hub"

	// ActionLog Permissions

	PackageActionLogView PermissionType = "view_package_action_log"
	TripActionLogView    PermissionType = "view_trip_action_log"

	// Account Permissions

	DriverAccountView PermissionType = "view_driver_account"
	UserAccountView   PermissionType = "view_user_account"

	// Transaction Permissions

	DriverTransactionView PermissionType = "view_driver_transaction"
	UserTransactionView   PermissionType = "view_user_transaction"
	UserTransactonEdit    PermissionType = "edit_user_transaction"

	// Trigger Permissions

	TriggerCreate   PermissionType = "create_trigger"
	TriggerView     PermissionType = "view_trigger"
	TriggerDelete   PermissionType = "delete_trigger"
	TriggerEdit     PermissionType = "edit_trigger"
	TriggerSchedule PermissionType = "schedule_trigger"

	// PackageTagRule Permissions

	PackageTagRuleCreate PermissionType = "create_package_tag_rule"
	PackageTagRuleView   PermissionType = "view_package_tag_rule"
	PackageTagRuleDelete PermissionType = "delete_package_tag_rule"
	PackageTagRuleEdit   PermissionType = "edit_package_tag_rule"

	// TripTagRule Permissions

	TripTagRuleCreate PermissionType = "create_trip_tag_rule"
	TripTagRuleView   PermissionType = "view_trip_tag_rule"
	TripTagRuleDelete PermissionType = "delete_trip_tag_rule"
	TripTagRuleEdit   PermissionType = "edit_trip_tag_rule"

	// NextHubRule Permissions

	NextHubRuleCreate PermissionType = "create_next_hub_rule"
	NextHubRuleView   PermissionType = "view_next_hub_rule"
	NextHubRuleDelete PermissionType = "delete_next_hub_rule"
	NextHubRuleEdit   PermissionType = "edit_next_hub_rule"

	// ShippingChannel Permissions

	ShippingChannelCreate PermissionType = "create_shipping_channel"
	ShippingChannelView   PermissionType = "view_shipping_channel"
	ShippingChannelEdit   PermissionType = "edit_shipping_channel"

	// Role Permissions

	TmsRoleCreate PermissionType = "create_tms_role"
	TmsRoleView   PermissionType = "view_tms_role"
	TmsRoleDelete PermissionType = "delete_tms_role"
	TmsRoleEdit   PermissionType = "edit_tms_role"

	// User Permissions

	TmsUserCreate PermissionType = "create_tms_user"
	TmsUserView   PermissionType = "view_tms_user"
	TmsUserDelete PermissionType = "delete_tms_user"
	TmsUserEdit   PermissionType = "edit_tms_user"

	// TmsTenant Permissions

	TmsTenantView PermissionType = "view_tms_tenant"
	TmsTenantEdit PermissionType = "edit_tms_tenant"

	// Partner Configuration Permissions

	PartnerConfigurationEdit PermissionType = "edit_partner_configuration"
	PartnerConfigurationView PermissionType = "view_partner_configuration"

	// Client Configuration Permissions

	ClientConfigurationCreate PermissionType = "create_client_configuration"
	ClientConfigurationView   PermissionType = "view_client_configuration"
	ClientConfigurationDelete PermissionType = "delete_client_configuration"
	ClientConfigurationEdit   PermissionType = "edit_client_configuration"

	// POS Permission

	POSView       PermissionType = "view_pos"
	POSCreate     PermissionType = "create_pos"
	POSSettingSet PermissionType = "set_pos_setting"

	RegisterCreate PermissionType = "create_register"
	RegisterView   PermissionType = "view_register"
	RegisterClose  PermissionType = "close_register"
	RegisterOpen   PermissionType = "open_register"

	CashManagementView PermissionType = "view_cash_management"

	//Automation Permission
	AddPackagingRule  PermissionType = "add_packaging_rule_automation_rules"
	ViewPackagingRule PermissionType = "view_packaging_rule_automation_rules"
	EditPackagingRule PermissionType = "edit_packaging_rule_automation_rules"

	// Pickup Locations Permissions
	PickupLocationCreate   PermissionType = "create_pickup_location"
	PickupLocationView     PermissionType = "view_pickup_location"
	PickupLocationEdit     PermissionType = "edit_pickup_location"
	PickupLocationV2View   PermissionType = "view_pickup_location_v2"
	PickupLocationV2Create PermissionType = "create_pickup_location_v2"
	PickupLocationV2Edit   PermissionType = "edit_pickup_location_v2"

	//Hub Routiung Permission
	AddHubRoutingRule  PermissionType = "add_hub_routing_rule_automation_rules"
	ViewHubRoutingRule PermissionType = "view_hub_routing_rule_automation_rules"
	EditHubRoutingRule PermissionType = "edit_hub_routing_rule_automation_rules"

	// Add Item Rule Permission
	CreateAddItemsRule PermissionType = "create_add_items_rule_automation_rules"
	ViewAddItemsRule   PermissionType = "view_add_items_rule_automation_rules"
	EditAddItemsRule   PermissionType = "edit_add_items_rule_automation_rules"

	// Membership Type permissions
	MembershipTypeCreate PermissionType = "create_membership_type"
	MembershipTypeEdit   PermissionType = "edit_membership_type"
	MembershipTypeView   PermissionType = "view_membership_type"
	MembershipTypeDelete PermissionType = "delete_membership_type"

	// Member permissions
	MemberView   PermissionType = "view_member"
	MemberDelete PermissionType = "delete_member"

	//DrugAuthorityPermissions
	ViewDrugAuthority   PermissionType = "view_drug_authority"
	EditDrugAuthority   PermissionType = "edit_drug_authority"
	CreateDrugAuthority PermissionType = "create_drug_authority"

	// Omniship Configuration Permissions
	OmnishipConfigurationView   PermissionType = "view_omniship_configuration"
	OmnishipConfigurationCreate PermissionType = "create_omniship_configuration"

	STOVirtualHubCreate PermissionType = "create_sto_virtual_hub"
	STOVirtualHubView   PermissionType = "view_sto_virtual_hub"
	STOVirtualHubEdit   PermissionType = "edit_sto_virtual_hub"

	// Tenant Custom Integration permissions
	TenantCustomIntegrationView PermissionType = "view_tenant_custom_integration"
	TenantCustomIntegrationEdit PermissionType = "edit_tenant_custom_integration"

	// Return Request Approval Automation Rule permissions
	ReturnRequestApprovalAutomationRulesAdd  PermissionType = "add_return_request_approval_automation_rules"
	ReturnRequestApprovalAutomationRulesEdit PermissionType = "edit_return_request_approval_automation_rules"
	ReturnRequestApprovalAutomationRulesView PermissionType = "view_return_request_approval_automation_rules"

	// Hub Reassignment Permission
	AddHubReassignmentRule  PermissionType = "add_hub_reassignment_automation_rules"
	ViewHubReassignmentRule PermissionType = "view_hub_reassignment_automation_rules"
	EditHubReassignmentRule PermissionType = "edit_hub_reassignment_automation_rules"
)
