package models

const (
	ReadOrganization     string = "Read Organization"
	UpdateOrganization   string = "Update Organization"
	DeleteOrganization   string = "Delete Organization"
	CreateRole           string = "Create Role"
	ReadRole             string = "Read Role"
	UpdateRole           string = "Update Role"
	DeleteRole           string = "Delete Role"
	CreateUser           string = "Create User"
	ReadUser             string = "Read User"
	UpdateUser           string = "Update User"
	DeleteUser           string = "Delete User"
	CreateCategoryOne    string = "Create Category One"
	ReadCategoryOne      string = "Read Category One"
	UpdateCategoryOne    string = "Update Category One"
	DeleteCategoryOne    string = "Delete Category One"
	CreateCategoryTwo    string = "Create Category Two"
	ReadCategoryTwo      string = "Read Category Two"
	UpdateCategoryTwo    string = "Update Category Two"
	DeleteCategoryTwo    string = "Delete Category Two"
	CreateSKU            string = "Create SKU"
	ReadSKU              string = "Read SKU"
	UpdateSKU            string = "Update SKU"
	DeleteSKU            string = "Delete SKU"
	CreateOrder          string = "Create Order"
	ReadOrder            string = "Read Order"
	UpdateOrder          string = "Update Order"
	DeleteOrder          string = "Delete Order"
	CreateContract       string = "Create Contract"
	ReadContract         string = "Read Contract"
	UpdateContract       string = "Update Contract"
	DeleteContract       string = "Delete Contract"
	CreateDistributor    string = "Create Distributor"
	ReadDistributor      string = "Read Distributor"
	UpdateDistributor    string = "Update Distributor"
	DeleteDistributor    string = "Delete Distributor"
	CreateContainer      string = "Create Container"
	ReadContainer        string = "Read Container"
	UpdateContainer      string = "Update Container"
	DeleteContainer      string = "Delete Container"
	CreatePallet         string = "Create Pallet"
	ReadPallet           string = "Read Pallet"
	UpdatePallet         string = "Update Pallet"
	DeletePallet         string = "Delete Pallet"
	CreateCarton         string = "Create Carton"
	ReadCarton           string = "Read Carton"
	UpdateCarton         string = "Update Carton"
	DeleteCarton         string = "Delete Carton"
	CreateProduct        string = "Create Product"
	ReadProduct          string = "Read Product"
	UpdateProduct        string = "Update Product"
	DeleteProduct        string = "Delete Product"
	CreateTask           string = "Create Task"
	ReadTask             string = "Read Task"
	UpdateTask           string = "Update Task"
	DeleteTask           string = "Delete Task"
	CreatePurchaseRecord string = "Create Purchase Record"
	ReadPurchaseRecord   string = "Read Purchase Record"
	UpdatePurchaseRecord string = "Update Purchase Record"
	DeletePurchaseRecord string = "Delete Purchase Record"
	CreateConsumerOrder  string = "Create Consumer Order"
	ReadConsumerOrder    string = "Read Consumer Order"
	UpdateConsumerOrder  string = "Update Consumer Order"
	DeleteConsumerOrder  string = "Delete Consumer Order"
	CreateTrackAction    string = "Create Track Action"
	ReadTrackAction      string = "Read Track Action"
	UpdateTrackAction    string = "Update Track Action"
	DeleteTrackAction    string = "Delete Track Action"
)

func ListPermissions() []string {
	return []string{
		ReadOrganization,
		UpdateOrganization,
		DeleteOrganization,
		CreateRole,
		ReadRole,
		UpdateRole,
		DeleteRole,
		CreateUser,
		ReadUser,
		UpdateUser,
		DeleteUser,
		CreateCategoryOne,
		ReadCategoryOne,
		UpdateCategoryOne,
		DeleteCategoryOne,
		CreateCategoryTwo,
		ReadCategoryTwo,
		UpdateCategoryTwo,
		DeleteCategoryTwo,
		CreateSKU,
		ReadSKU,
		UpdateSKU,
		DeleteSKU,
		CreateOrder,
		ReadOrder,
		UpdateOrder,
		DeleteOrder,
		CreateContract,
		ReadContract,
		UpdateContract,
		DeleteContract,
		CreateDistributor,
		ReadDistributor,
		UpdateDistributor,
		DeleteDistributor,
		CreateContainer,
		ReadContainer,
		UpdateContainer,
		DeleteContainer,
		CreatePallet,
		ReadPallet,
		UpdatePallet,
		DeletePallet,
		CreateCarton,
		ReadCarton,
		UpdateCarton,
		DeleteCarton,
		CreateProduct,
		ReadProduct,
		UpdateProduct,
		DeleteProduct,
		CreateTask,
		ReadTask,
		UpdateTask,
		DeleteTask,
		CreatePurchaseRecord,
		ReadPurchaseRecord,
		UpdatePurchaseRecord,
		DeletePurchaseRecord,
		CreateConsumerOrder,
		ReadConsumerOrder,
		UpdateConsumerOrder,
		DeleteConsumerOrder,
		CreateTrackAction,
		ReadTrackAction,
		UpdateTrackAction,
		DeleteTrackAction,
	}
}
